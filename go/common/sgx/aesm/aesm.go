// Package aesm provides a client for AESMD.
package aesm

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/oasislabs/ekiden/go/common/sgx/ias"
)

//go:generate protoc --go_out=. aesm_proto.proto

var errMalformedResponse = errors.New("aesm: malformed response")

var (
	// localAESMTimeout is the timeout for local requests.
	localAESMTimeout = 1 * time.Second

	// remoteAESMTimeout is the timeout for remote requests.
	remoteAESMTimeout = 30 * time.Second
)

// QuoteInfo is the quote information.
type QuoteInfo struct {
	// TargetInfo is the target enclave info.
	TargetInfo []byte
	// GID is an EPID group ID.
	GID []byte
}

// Client is an AESM client.
type Client struct {
	path string
}

// NewClient creates a new AESM client.
func NewClient(path string) *Client {
	return &Client{
		path: path,
	}
}

func (c *Client) transact(ctx context.Context, request *Request) (*Response, error) {
	// The AESM socket only accepts one request per connection, so we
	// need to establish a new connection for each request.
	conn, err := net.Dial("unix", c.path)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rsp := make(chan interface{})
	go func() {
		defer close(rsp)

		// Marshal request and send it as a length-prefixed blob.
		body, err := proto.Marshal(request)
		if err != nil {
			rsp <- err
			return
		}
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(body)))
		if _, err := conn.Write(buf); err != nil {
			rsp <- err
			return
		}
		if _, err := conn.Write(body); err != nil {
			rsp <- err
			return
		}

		// Receive size and encoded response.
		if _, err := io.ReadFull(conn, buf); err != nil {
			rsp <- err
			return
		}
		len := binary.LittleEndian.Uint32(buf)
		buf = make([]byte, len)
		if _, err := io.ReadFull(conn, buf); err != nil {
			rsp <- err
			return
		}

		response := &Response{}
		if err := proto.Unmarshal(buf, response); err != nil {
			rsp <- err
			return
		}

		rsp <- response
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-rsp:
		switch r := r.(type) {
		case error:
			return nil, r
		case *Response:
			return r, nil
		default:
			panic("invalid response type")
		}
	}
}

// InitQuote retrieves the quote info required for generating a report
// that can be exchanged for a quote.
func (c *Client) InitQuote(ctx context.Context) (*QuoteInfo, error) {
	timeout := uint32(localAESMTimeout.Nanoseconds() / 1000)
	resp, err := c.transact(ctx, &Request{
		InitQuoteReq: &Request_InitQuoteRequest{
			Timeout: &timeout,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.InitQuoteRes == nil {
		return nil, errMalformedResponse
	}
	if errCode := resp.InitQuoteRes.GetErrorCode(); errCode != 0 {
		return nil, fmt.Errorf("aesm: error %d", errCode)
	}

	qi := &QuoteInfo{
		TargetInfo: resp.InitQuoteRes.TargetInfo,
		GID:        resp.InitQuoteRes.Gid,
	}
	return qi, nil
}

// GetQuote retrieves the quote based on the provided report.
func (c *Client) GetQuote(
	ctx context.Context,
	report []byte,
	quoteType ias.SignatureType,
	spid ias.SPID,
	nonce []byte,
	sigRL []byte,
) ([]byte, error) {
	// Compute the buffer size for the quote. For some reason Intel decided
	// that you need to specify the buffer size in a protobuf-based API.
	const (
		// Offset of signature_len in sgx_quote_t (see common/inc/sgx_quote.h).
		sgxOffsetQuoteSigLen = 432
		// Size of sgx_quote_t without signature (see common/inc/sgx_quote.h).
		sgxQuoteStructLen = 436
		// Size of se_wrap_key_t (see common/inc/internal/se_quote_internal.h).
		sgxWrapKeyLen = 288
		// Size of quote IV (see common/inc/internal/se_quote_internal.h).
		sgxQuoteIVLen = 12
		// Size of payload_size field (see common/inc/internal/se_quote_internal.h).
		sgxQuotePayloadSize = 4
		// Size of sgx_mac_t (see common/inc/sgx_report.h).
		sgxMACLen = 16
		// Size of quote without signature (see SE_QUOTE_LENGTH_WITHOUT_SIG in common/inc/internal/se_quote_internal.h).
		sgxQuoteLen = sgxQuoteStructLen + sgxWrapKeyLen + sgxQuoteIVLen + sgxQuotePayloadSize + sgxMACLen

		// Size of EPID basic signature (see external/epid-sdk/epid/common/types.h).
		sgxBasicSigLen = 352
		// Size of RLver_t (see external/epid-sdk/epid/common/types.h).
		sgxRLVerLen = 4
		// Size of RLCount (see external/epid-sdk/epid/common/types.h).
		sgxRLCountLen = 4
		// Size of static signature part (without SigRL).
		sgxSigLen = sgxBasicSigLen + sgxRLVerLen + sgxRLCountLen
	)

	// This is the truly correct way to compute sigLen:
	//
	//   const sgxNrProofLen = 160
	//   sigLen := sgxSigLen + sigRLEntries * sgxNrProofLen
	//
	// Instead we do something that should be conservative, and doesn't
	// require interpreting the sigRL structure to determine the entry
	// count. An NrProof is 5 field elements, a sigRL entry is four.
	// Add some slop for sigRL headers.
	sigLen := sgxSigLen + (len(sigRL) * 5 / 4) + 128
	bufSize := uint32(sgxQuoteLen + sigLen)

	qeReport := true
	qt := uint32(int(quoteType))

	if len(sigRL) == 0 {
		sigRL = nil
	}

	timeout := uint32(remoteAESMTimeout.Nanoseconds() / 1000)
	resp, err := c.transact(ctx, &Request{
		GetQuoteReq: &Request_GetQuoteRequest{
			Report:    report,
			QuoteType: &qt,
			Spid:      spid[:],
			Nonce:     nonce,
			SigRl:     sigRL,
			BufSize:   &bufSize,
			QeReport:  &qeReport,
			Timeout:   &timeout,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.GetQuoteRes == nil {
		return nil, errMalformedResponse
	}
	if errCode := resp.GetQuoteRes.GetErrorCode(); errCode != 0 {
		return nil, fmt.Errorf("aesm: error %d", errCode)
	}

	quote := resp.GetQuoteRes.Quote
	// AESM allocates a buffer of the size we supplied and returns the whole
	// thing to us, regardless of how much space QE needed. Trim the excess.
	// The signature length is a little endian word at offset 432 in the quote
	// structure. See "QUOTE Structure" in the IAS API Spec.
	sigLen = int(binary.LittleEndian.Uint32(quote[sgxOffsetQuoteSigLen : sgxOffsetQuoteSigLen+4]))
	newLen := sgxQuoteStructLen + sigLen
	if len(quote) < newLen {
		// Quote is already too short, should not happen.
		// Probably we are interpreting the quote structure incorrectly.
		return nil, errMalformedResponse
	}
	quote = quote[:newLen]

	return quote, nil
}
