package gateway

import (
	"context"

	// The high-level gateway module depends on the abstraction, not the detail.
	"nexusfs/internal/metadata"
)

// Service is the high-level module for the gateway.
// It embeds the dependencies it needs as interfaces.
type Service struct {
	metaDB metadata.MetadataReadWriter // Dependency is the interface.
	// TODO: Add other dependencies like an erasure coding client
	// and a gRPC client for storage nodes.
}

// NewService uses dependency injection to receive its dependencies.
// This decouples the gateway from the concrete implementation of the metadata store.
func NewService(metaDB metadata.MetadataReadWriter) *Service {
	return &Service{
		metaDB: metaDB,
	}
}

// GetObject is an example method that uses the metadata interface.
func (s *Service) GetObject(ctx context.Context, bucket, object string) (*metadata.Object, error) {
	// The service interacts with the abstraction, completely unaware of whether
	// the underlying database is PostgreSQL, TiKV, or something else.
	return s.metaDB.GetObjectMetadata(ctx, bucket, object)
}
