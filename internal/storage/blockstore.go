package storage

import "context"

// TODO: Implement the BlockStore for local disk.
// This will handle reading, writing, and deleting shard data from the filesystem.

type BlockStore struct {
	dataDir string
}

func NewBlockStore(dataDir string) *BlockStore {
	return &BlockStore{dataDir: dataDir}
}

func (bs *BlockStore) WriteShard(ctx context.Context, shardID string, data []byte) error {
	// TODO: Implement file writing logic.
	return nil
}

func (bs *BlockStore) ReadShard(ctx context.Context, shardID string) ([]byte, error) {
	// TODO: Implement file reading logic.
	return nil, nil
}
