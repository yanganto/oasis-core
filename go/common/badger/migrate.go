package badger

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"math"

	badgerV2 "github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/pb"
	"github.com/golang/protobuf/proto" //nolint: staticcheck
)

// Adapted from Badger v2 which is Copyright 2017 Dgraph Labs, Inc. and Contributors, released
// under the Apache-2 license.

func backup(db *badgerV2.DB, w io.Writer, managed bool) (uint64, error) {
	var stream *badgerV2.Stream
	switch managed {
	case true:
		stream = db.NewStreamAt(math.MaxUint64)
	case false:
		stream = db.NewStream()
	}

	stream.LogPrefix = "migration"
	stream.KeyToList = func(key []byte, itr *badgerV2.Iterator) (*pb.KVList, error) {
		list := &pb.KVList{}
		for ; itr.Valid(); itr.Next() {
			item := itr.Item()
			if !bytes.Equal(item.Key(), key) {
				return list, nil
			}

			var valCopy []byte
			var meta byte
			switch item.IsDeletedOrExpired() {
			case true:
				// No need to copy value, if item is deleted or expired.
				// Set delete bit.
				meta = 1 << 0 // bitDelete
			case false:
				var err error
				valCopy, err = item.ValueCopy(nil)
				if err != nil {
					return nil, err
				}
			}

			kv := &pb.KV{
				Key:       item.KeyCopy(nil),
				Value:     valCopy,
				UserMeta:  []byte{item.UserMeta()},
				Version:   item.Version(),
				ExpiresAt: item.ExpiresAt(),
				Meta:      []byte{meta},
			}
			list.Kv = append(list.Kv, kv)
		}
		return list, nil
	}

	var maxVersion uint64
	stream.Send = func(list *pb.KVList) error {
		for _, kv := range list.Kv {
			if maxVersion < kv.Version {
				maxVersion = kv.Version
			}
		}
		if err := binary.Write(w, binary.LittleEndian, uint64(proto.Size(list))); err != nil {
			return err
		}
		buf, err := proto.Marshal(list)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}

	if err := stream.Orchestrate(context.Background()); err != nil {
		return 0, err
	}
	return maxVersion, nil
}