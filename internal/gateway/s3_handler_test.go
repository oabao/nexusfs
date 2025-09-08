package gateway

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	"nexusfs/internal/erasure"
	"nexusfs/internal/metadata/mock"
	"nexusfs/internal/storage"
)

func TestS3Api_CreateBucketHandler(t *testing.T) {
	// Setup
	mockStore := mock.NewStore()
	encoder, err := erasure.NewEncoder()
	if err != nil {
		t.Fatalf("failed to create encoder: %v", err)
	}

	dataDir, err := os.MkdirTemp("", "handler-test")
	if err != nil {
		t.Fatalf("Failed to create temp data dir: %v", err)
	}
	defer os.RemoveAll(dataDir)
	blockStore, err := storage.NewBlockStore(dataDir)
	if err != nil {
		t.Fatalf("Failed to create block store: %v", err)
	}

	service := NewService(mockStore, encoder, blockStore)
	api := NewS3Api(service)

	bucketName := "test-bucket"
	req, err := http.NewRequest("PUT", "/"+bucketName, nil)
	if err != nil {
		t.Fatal(err)
	}

	// We need a router because the handler uses mux.Vars
	router := mux.NewRouter()
	router.HandleFunc("/{bucket}", api.CreateBucketHandler)

	// Use httptest.ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if mockStore.CreatedBucket != bucketName {
		t.Errorf("expected bucket '%s' to be created, but got '%s'", bucketName, mockStore.CreatedBucket)
	}
}

// TODO: Add tests for ListBuckets, PutObject, etc.
