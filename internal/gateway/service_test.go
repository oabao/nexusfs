package gateway

import (
	"context"
	"os"
	"testing"

	"nexusfs/internal/erasure"
	"nexusfs/internal/metadata/mock"
	"nexusfs/internal/storage"
)

func TestService_CreateBucket(t *testing.T) {
	// Setup
	mockStore := mock.NewStore()
	encoder, err := erasure.NewEncoder()
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	// Create a temporary block store for the test
	dataDir, err := os.MkdirTemp("", "service-test")
	if err != nil {
		t.Fatalf("Failed to create temp data dir: %v", err)
	}
	defer os.RemoveAll(dataDir)
	blockStore, err := storage.NewBlockStore(dataDir)
	if err != nil {
		t.Fatalf("Failed to create block store: %v", err)
	}

	service := NewService(mockStore, encoder, blockStore)
	bucketName := "test-bucket"

	// Execute
	err = service.CreateBucket(context.Background(), bucketName)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if mockStore.CreatedBucket != bucketName {
		t.Errorf("expected bucket '%s' to be created, but got '%s'", bucketName, mockStore.CreatedBucket)
	}
}

// TODO: Add tests for PutObject, GetObject, etc.
