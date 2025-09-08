package gateway

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"

	"nexusfs/internal/erasure"
	"nexusfs/internal/metadata"
	"nexusfs/internal/storage"

	"github.com/google/uuid"
)

// Service is the high-level module for the gateway.
// It embeds the dependencies it needs as interfaces.
type Service struct {
	metaDB     metadata.MetadataReadWriter
	encoder    *erasure.Encoder
	blockStore *storage.BlockStore
}

// NewService uses dependency injection to receive its dependencies.
func NewService(metaDB metadata.MetadataReadWriter, encoder *erasure.Encoder, blockStore *storage.BlockStore) *Service {
	return &Service{
		metaDB:     metaDB,
		encoder:    encoder,
		blockStore: blockStore,
	}
}

// PutObject handles the business logic of a PutObject call.
func (s *Service) PutObject(ctx context.Context, bucketName, objectName string, data io.Reader) (string, error) {
	// 1. Tee the reader to calculate the ETag (md5sum) while streaming.
	hash := md5.New()
	teeReader := io.TeeReader(data, hash)

	// 2. Erasure code the data.
	shards, err := s.encoder.Encode(teeReader)
	if err != nil {
		return "", fmt.Errorf("failed to erasure code object: %w", err)
	}

	// 3. Write each shard to the block store.
	shardLocations := make(map[string]string)
	for i, shardData := range shards {
		shardID := uuid.New().String()
		if err := s.blockStore.WriteShard(ctx, shardID, shardData); err != nil {
			// TODO: In a real system, this would need a cleanup/rollback mechanism.
			return "", fmt.Errorf("failed to write shard %d: %w", i, err)
		}
		shardLocations[fmt.Sprintf("shard_%d", i)] = shardID
	}

	// 4. Calculate the ETag.
	etag := hex.EncodeToString(hash.Sum(nil))

	// 5. Write the object's metadata to the database.
	objectMeta := &metadata.Object{
		Name:           objectName,
		Bucket:         bucketName,
		ETag:           etag,
		ShardLocations: shardLocations,
	}

	if err := s.metaDB.WriteObjectMetadata(ctx, objectMeta); err != nil {
		// TODO: Implement cleanup of orphaned shards if metadata write fails.
		return "", fmt.Errorf("failed to write object metadata: %w", err)
	}

	return etag, nil
}

// GetObject handles the business logic of a GetObject call.
func (s *Service) GetObject(ctx context.Context, bucketName, objectName string, w io.Writer) error {
	// 1. Get the object's metadata.
	meta, err := s.metaDB.GetObjectMetadata(ctx, bucketName, objectName)
	if err != nil {
		return fmt.Errorf("failed to get object metadata: %w", err)
	}
	if meta == nil {
		return fmt.Errorf("object %s/%s not found", bucketName, objectName)
	}

	// 2. Read all shards from the block store.
	shards := make([][]byte, erasure.DefaultDataShards+erasure.DefaultParityShards)
	for i := 0; i < len(shards); i++ {
		shardKey := fmt.Sprintf("shard_%d", i)
		shardID := meta.ShardLocations[shardKey]
		shards[i], err = s.blockStore.ReadShard(ctx, shardID)
		if err != nil {
			// A missing shard is okay (up to the parity limit). We pass nil to Reconstruct.
			log.Printf("Warning: could not read shard %s: %v", shardID, err)
			shards[i] = nil
		}
	}

	// 3. Reconstruct the data shards.
	if err := s.encoder.Reconstruct(shards); err != nil {
		return fmt.Errorf("failed to reconstruct object shards: %w", err)
	}

	// 4. Join the data shards and stream them to the writer.
	if err := s.encoder.Join(w, shards); err != nil {
		return fmt.Errorf("failed to join object shards: %w", err)
	}

	return nil
}

// CreateBucket creates a new bucket.
func (s *Service) CreateBucket(ctx context.Context, bucketName string) error {
	bucket := metadata.Bucket{
		Name:    bucketName,
		Created: time.Now().UTC(),
	}
	return s.metaDB.CreateBucket(ctx, bucket)
}

// ListBuckets lists all buckets.
func (s *Service) ListBuckets(ctx context.Context) ([]metadata.Bucket, error) {
	return s.metaDB.ListBuckets(ctx)
}

// ListObjects lists objects in a bucket.
func (s *Service) ListObjects(ctx context.Context, bucketName string) ([]metadata.Object, error) {
	return s.metaDB.ListObjects(ctx, bucketName)
}

// DeleteObject deletes an object's metadata.
func (s *Service) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	// In a real system, this would also trigger a garbage collection
	// job for the data shards.
	return s.metaDB.DeleteObjectMetadata(ctx, bucketName, objectName)
}
