package storage

import "context"

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// BlockStore handles the low-level reading and writing of data blocks (shards)
// to the local filesystem.
type BlockStore struct {
	dataDir string
}

// NewBlockStore creates a new BlockStore and ensures the data directory exists.
func NewBlockStore(dataDir string) (*BlockStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &BlockStore{dataDir: dataDir}, nil
}

// WriteShard writes a byte slice representing a shard to a file on disk.
// The file is named after the shardID.
func (bs *BlockStore) WriteShard(ctx context.Context, shardID string, data []byte) error {
	filePath := filepath.Join(bs.dataDir, shardID)
	return ioutil.WriteFile(filePath, data, 0644)
}

// ReadShard reads a shard's data from a file on disk.
func (bs *BlockStore) ReadShard(ctx context.Context, shardID string) ([]byte, error) {
	filePath := filepath.Join(bs.dataDir, shardID)
	return ioutil.ReadFile(filePath)
}
