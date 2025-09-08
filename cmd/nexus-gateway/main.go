package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"nexusfs/internal/erasure"
	"nexusfs/internal/gateway"
	"nexusfs/internal/metadata/mock" // Use the mock store
	"nexusfs/internal/storage"
)

func main() {
	// Use the functional, in-memory mock store for metadata.
	// This makes the server stateful for its runtime.
	metaStore := mock.NewStore()

	// Initialize the erasure encoder
	encoder, err := erasure.NewEncoder()
	if err != nil {
		log.Fatalf("Failed to create erasure encoder: %v", err)
	}

	// Initialize the block store in a temporary directory
	dataDir, err := os.MkdirTemp("", "nexusfs-data")
	if err != nil {
		log.Fatalf("Failed to create temp data dir: %v", err)
	}
	defer os.RemoveAll(dataDir) // Clean up on exit
	log.Printf("Using data directory: %s", dataDir)

	blockStore, err := storage.NewBlockStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to create block store: %v", err)
	}

	// The service now depends on the mock store and the block store.
	gwService := gateway.NewService(metaStore, encoder, blockStore)
	s3Api := gateway.NewS3Api(gwService)

	r := mux.NewRouter()

	// Register S3 API handlers
	r.HandleFunc("/", s3Api.ListBucketsHandler).Methods("GET")
	r.HandleFunc("/{bucket}", s3Api.CreateBucketHandler).Methods("PUT")
	r.HandleFunc("/{bucket}", s3Api.ListObjectsV2Handler).Methods("GET").Queries("list-type", "2")
	r.HandleFunc("/{bucket}/{object:.+}", s3Api.PutObjectHandler).Methods("PUT")
	r.HandleFunc("/{bucket}/{object:.+}", s3Api.GetObjectHandler).Methods("GET")
	r.HandleFunc("/{bucket}/{object:.+}", s3Api.DeleteObjectHandler).Methods("DELETE")

	log.Println("Starting NexusFS Gateway on :9000")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
