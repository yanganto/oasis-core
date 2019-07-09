package storage

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/oasislabs/ekiden/go/common"
	"github.com/oasislabs/ekiden/go/common/cbor"
	"github.com/oasislabs/ekiden/go/common/crypto/hash"
	"github.com/oasislabs/ekiden/go/storage/api"

	pb "github.com/oasislabs/ekiden/go/grpc/storage"
)

const (
	// WriteLogIteratorChunkSize defines the chunk size of write log entries
	// for GetCheckpoint and GetDiff methods.
	WriteLogIteratorChunkSize int = 10
)

var _ pb.StorageServer = (*grpcServer)(nil)

type grpcServer struct {
	backend api.Backend
}

func (s *grpcServer) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	var ns common.Namespace
	if err := ns.UnmarshalBinary(req.GetNamespace()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal namespace")
	}

	var srcRoot, dstRoot hash.Hash
	if err := srcRoot.UnmarshalBinary(req.GetSrcRoot()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal src root")
	}
	if err := dstRoot.UnmarshalBinary(req.GetDstRoot()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal dst root")
	}

	var log api.WriteLog
	for _, item := range req.GetLog() {
		log = append(log, api.LogEntry{
			Key:   item.GetKey(),
			Value: item.GetValue(),
		})
	}

	<-s.backend.Initialized()
	receipts, err := s.backend.Apply(ctx, ns, req.GetSrcRound(), srcRoot, req.GetDstRound(), dstRoot, log)

	if err != nil {
		return nil, err
	}

	return &pb.ApplyResponse{Receipts: cbor.Marshal(receipts)}, nil
}

func (s *grpcServer) ApplyBatch(ctx context.Context, req *pb.ApplyBatchRequest) (*pb.ApplyBatchResponse, error) {
	var ns common.Namespace
	if err := ns.UnmarshalBinary(req.GetNamespace()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal namespace")
	}

	var ops []api.ApplyOp
	for _, op := range req.GetOps() {
		var srcRoot, dstRoot hash.Hash
		if err := srcRoot.UnmarshalBinary(op.GetSrcRoot()); err != nil {
			return nil, errors.Wrap(err, "storage: failed to unmarshal src root")
		}
		if err := dstRoot.UnmarshalBinary(op.GetDstRoot()); err != nil {
			return nil, errors.Wrap(err, "storage: failed to unmarshal dst root")
		}

		var log api.WriteLog
		for _, item := range op.GetLog() {
			log = append(log, api.LogEntry{
				Key:   item.GetKey(),
				Value: item.GetValue(),
			})
		}

		ops = append(ops, api.ApplyOp{
			SrcRound: op.GetSrcRound(),
			SrcRoot:  srcRoot,
			DstRoot:  dstRoot,
			WriteLog: log,
		})
	}

	<-s.backend.Initialized()
	receipts, err := s.backend.ApplyBatch(ctx, ns, req.GetDstRound(), ops)

	if err != nil {
		return nil, err
	}

	return &pb.ApplyBatchResponse{Receipts: cbor.Marshal(receipts)}, nil
}

func (s *grpcServer) GetSubtree(ctx context.Context, req *pb.GetSubtreeRequest) (*pb.GetSubtreeResponse, error) {
	var root api.Root
	if err := root.UnmarshalCBOR(req.GetRoot()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal root")
	}

	maxDepth := uint8(req.GetMaxDepth())

	nid := req.GetId()
	var path hash.Hash
	if err := path.UnmarshalBinary(nid.GetPath()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal id")
	}

	nodeID := api.NodeID{
		Path:  path,
		Depth: uint8(nid.GetDepth()),
	}

	<-s.backend.Initialized()
	subtree, err := s.backend.GetSubtree(ctx, root, nodeID, maxDepth)
	if err != nil {
		return nil, err
	}

	serializedSubtree, err := subtree.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &pb.GetSubtreeResponse{Subtree: serializedSubtree}, nil
}

func (s *grpcServer) GetPath(ctx context.Context, req *pb.GetPathRequest) (*pb.GetPathResponse, error) {
	var root api.Root
	if err := root.UnmarshalCBOR(req.GetRoot()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal root")
	}

	var key hash.Hash
	if err := key.UnmarshalBinary(req.GetKey()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal key")
	}

	startDepth := uint8(req.GetStartDepth())

	<-s.backend.Initialized()
	subtree, err := s.backend.GetPath(ctx, root, key, startDepth)
	if err != nil {
		return nil, err
	}

	serializedSubtree, err := subtree.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &pb.GetPathResponse{Subtree: serializedSubtree}, nil
}

func (s *grpcServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	var root api.Root
	if err := root.UnmarshalCBOR(req.GetRoot()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal root")
	}

	nid := req.GetId()
	var path hash.Hash
	if err := path.UnmarshalBinary(nid.GetPath()); err != nil {
		return nil, errors.Wrap(err, "storage: failed to unmarshal id")
	}

	nodeID := api.NodeID{
		Path:  path,
		Depth: uint8(nid.GetDepth()),
	}

	<-s.backend.Initialized()
	node, err := s.backend.GetNode(ctx, root, nodeID)
	if err != nil {
		return nil, err
	}

	serializedNode, err := node.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &pb.GetNodeResponse{Node: serializedNode}, nil
}

// writeLogService implements sending write log iterator for GetDiff and GetCheckpoint methods.
type writeLogService struct {
	opts     *pb.SyncOptions
	iterator api.WriteLogIterator
	send     func(*pb.WriteLogResponse) error
}

func (s *writeLogService) SendWriteLogIterator() error {
	var totalSent uint64
	skipping := true
	final := false
	done := false
	totalSent = 0

	if len(s.opts.GetOffsetKey()) == 0 {
		skipping = false
	}

	for {
		var entryArray []*pb.LogEntry
		for {
			more, err := s.iterator.Next()
			if err != nil {
				return err
			}
			if !more {
				final = true
				break
			}

			entry, err := s.iterator.Value()
			if err != nil {
				return err
			}

			if skipping {
				if bytes.Equal(entry.Key, s.opts.GetOffsetKey()) {
					skipping = false
				}
				continue
			}

			entryArray = append(entryArray, &pb.LogEntry{
				Key:   entry.Key,
				Value: entry.Value,
			})
			totalSent++
			if len(entryArray) >= WriteLogIteratorChunkSize {
				break
			}
			if s.opts.GetLimit() > 0 && totalSent >= s.opts.GetLimit() {
				done = true
				break
			}
		}
		resp := &pb.WriteLogResponse{
			Final: final,
			Log:   entryArray,
		}

		if err := s.send(resp); err != nil {
			return err
		}

		if done || final {
			break
		}
	}

	return nil
}

func (s *grpcServer) GetDiff(req *pb.GetDiffRequest, stream pb.Storage_GetDiffServer) error {
	var startRoot, endRoot api.Root
	if err := startRoot.UnmarshalCBOR(req.GetStartRoot()); err != nil {
		return errors.Wrap(err, "storage: failed to unmarshal start root")
	}
	if err := endRoot.UnmarshalCBOR(req.GetEndRoot()); err != nil {
		return errors.Wrap(err, "storage: failed to unmarshal end root")
	}

	<-s.backend.Initialized()

	it, err := s.backend.GetDiff(stream.Context(), startRoot, endRoot)
	if err != nil {
		return err
	}

	svc := &writeLogService{
		opts:     req.GetOpts(),
		iterator: it,
		send:     func(resp *pb.WriteLogResponse) error { return stream.Send(resp) },
	}

	return svc.SendWriteLogIterator()
}

func (s *grpcServer) GetCheckpoint(req *pb.GetCheckpointRequest, stream pb.Storage_GetCheckpointServer) error {
	var root api.Root
	if err := root.UnmarshalCBOR(req.GetRoot()); err != nil {
		return errors.Wrap(err, "storage: failed to unmarshal root")
	}

	<-s.backend.Initialized()

	it, err := s.backend.GetCheckpoint(stream.Context(), root)
	if err != nil {
		return err
	}

	svc := &writeLogService{
		opts:     req.GetOpts(),
		iterator: it,
		send:     func(resp *pb.WriteLogResponse) error { return stream.Send(resp) },
	}

	return svc.SendWriteLogIterator()
}

// NewGRPCServer intializes and registers a grpc storage server backed
// by the provided Backend.
func NewGRPCServer(srv *grpc.Server, b api.Backend) {
	s := &grpcServer{
		backend: b,
	}

	pb.RegisterStorageServer(srv, s)
}
