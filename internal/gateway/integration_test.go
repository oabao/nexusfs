package gateway

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"nexusfs/internal/erasure"
	"nexusfs/internal/metadata/mock"
	"nexusfs/internal/storage"
)

// TestIntegration_PutAndGetObject performs an end-to-end test of the Put and Get object lifecycle.
func TestIntegration_PutAndGetObject(t *testing.T) {
	// 1. Setup all "real" components
	metaStore := mock.NewStore()
	encoder, err := erasure.NewEncoder()
	if err != nil {
		t.Fatalf("Failed to create erasure encoder: %v", err)
	}

	dataDir, err := os.MkdirTemp("", "nexusfs-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp data dir: %v", err)
	}
	defer os.RemoveAll(dataDir)

	blockStore, err := storage.NewBlockStore(dataDir)
	if err != nil {
		t.Fatalf("Failed to create block store: %v", err)
	}

	service := NewService(metaStore, encoder, blockStore)
	api := NewS3Api(service)

	// 2. Setup the router and test server
	router := mux.NewRouter()
	router.HandleFunc("/{bucket}", api.CreateBucketHandler).Methods("PUT")
	router.HandleFunc("/{bucket}/{object:.+}", api.PutObjectHandler).Methods("PUT")
	router.HandleFunc("/{bucket}/{object:.+}", api.GetObjectHandler).Methods("GET")

	ts := httptest.NewServer(router)
	defer ts.Close()

	// 3. Create a bucket first
	bucketName := "test-bucket"
	createBucketReq, _ := http.NewRequest("PUT", ts.URL+"/"+bucketName, nil)
	createBucketResp, err := http.DefaultClient.Do(createBucketReq)
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}
	if createBucketResp.StatusCode != http.StatusOK {
		t.Fatalf("Create bucket failed with status: %s", createBucketResp.Status)
	}

	// 4. Perform the PUT request
	objectName := "my-test-object"
	objectContent := "hello world, this is a test of the nexusfs object storage system"
	putReq, _ := http.NewRequest("PUT", ts.URL+"/"+bucketName+"/"+objectName, strings.NewReader(objectContent))

	putResp, err := http.DefaultClient.Do(putReq)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if putResp.StatusCode != http.StatusOK {
		t.Fatalf("PUT handler returned wrong status code: got %v want %v", putResp.StatusCode, http.StatusOK)
	}

	// 5. Perform the GET request
	getReq, _ := http.NewRequest("GET", ts.URL+"/"+bucketName+"/"+objectName, nil)
	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET handler returned wrong status code: got %v want %v", getResp.StatusCode, http.StatusOK)
	}

	// 6. Verify the content
	retrievedContent, err := ioutil.ReadAll(getResp.Body)
	getResp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read GET response body: %v", err)
	}

	if !bytes.Equal([]byte(objectContent), retrievedContent) {
		t.Errorf("retrieved content does not match original content.\nOriginal:  %s\nRetrieved: %s", objectContent, string(retrievedContent))
	}

	t.Log("Integration test passed: Put and Get successful.")
}
