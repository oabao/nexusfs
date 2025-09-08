package gateway

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// S3Api is responsible for handling the S3-compatible HTTP requests.
type S3Api struct {
	service *Service
}

func NewS3Api(service *Service) *S3Api {
	return &S3Api{service: service}
}

// ListAllMyBucketsResult is the S3-compatible XML structure for ListBuckets response.
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   struct {
		ID          string `xml:"ID"`
		DisplayName string `xml:"DisplayName"`
	} `xml:"Owner"`
	Buckets struct {
		Bucket []struct {
			Name         string `xml:"Name"`
			CreationDate string `xml:"CreationDate"`
		} `xml:"Bucket"`
	} `xml:"Buckets"`
}

// ListBucketsHandler handles the `GET /` request.
func (api *S3Api) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := api.service.ListBuckets(r.Context())
	if err != nil {
		// TODO: Write S3-compatible error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ListAllMyBucketsResult{}
	// TODO: Fill Owner info properly
	response.Owner.ID = "test-id"
	response.Owner.DisplayName = "test-user"

	for _, b := range buckets {
		response.Buckets.Bucket = append(response.Buckets.Bucket, struct {
			Name         string `xml:"Name"`
			CreationDate string `xml:"CreationDate"`
		}{
			Name:         b.Name,
			CreationDate: b.Created.Format(time.RFC3339),
		})
	}

	xmlBytes, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}

// CreateBucketHandler handles the `PUT /{bucket}` request.
func (api *S3Api) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]

	err := api.service.CreateBucket(r.Context(), bucketName)
	if err != nil {
		// TODO: Write S3-compatible error response (e.g., BucketAlreadyExists)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetObjectHandler is the HTTP handler for GET Object requests.
func (api *S3Api) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]
	objectName := vars["object"]

	// TODO: Set response headers like Content-Type, Content-Length, ETag.

	// The http.ResponseWriter can be used as an io.Writer, so we can stream
	// the reconstructed object data directly to the client.
	err := api.service.GetObject(r.Context(), bucketName, objectName, w)
	if err != nil {
		// TODO: Write S3-compatible error response (e.g., NoSuchKey).
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

// ListBucketResult is the S3-compatible XML structure for ListObjectsV2 response.
type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	MaxKeys     int      `xml:"MaxKeys"`
	IsTruncated bool     `xml:"IsTruncated"`
	Contents    []struct {
		Key          string    `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string    `xml:"ETag"`
		Size         int64     `xml:"Size"`
	} `xml:"Contents"`
	// TODO: Add CommonPrefixes for directory-like listing.
}

// ListObjectsV2Handler handles the `GET /{bucket}` request with `list-type=2` query param.
func (api *S3Api) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]

	objects, err := api.service.ListObjects(r.Context(), bucketName)
	if err != nil {
		// TODO: Write S3-compatible error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Handle prefix, marker, max-keys from query params.
	response := ListBucketResult{
		Name:        bucketName,
		MaxKeys:     1000,
		IsTruncated: false,
	}

	for _, obj := range objects {
		response.Contents = append(response.Contents, struct {
			Key          string    `xml:"Key"`
			LastModified string `xml:"LastModified"`
			ETag         string    `xml:"ETag"`
			Size         int64     `xml:"Size"`
		}{
			Key:          obj.Name,
			LastModified: time.Now().UTC().Format(time.RFC3339), // Placeholder
			ETag:         `"` + obj.ETag + `"`,
			Size:         obj.Size,
		})
	}

	xmlBytes, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlBytes)
}


// DeleteObjectHandler handles the `DELETE /{bucket}/{object}` request.
func (api *S3Api) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]
	objectName := vars["object"]

	err := api.service.DeleteObject(r.Context(), bucketName, objectName)
	if err != nil {
		// TODO: Write S3-compatible error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// S3 spec says to return 204 No Content on successful deletion.
	w.WriteHeader(http.StatusNoContent)
}

// PutObjectHandler is the HTTP handler for PUT Object requests.
func (api *S3Api) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucketName := vars["bucket"]
	objectName := vars["object"]

	// The request body is the object's data.
	defer r.Body.Close()

	etag, err := api.service.PutObject(r.Context(), bucketName, objectName, r.Body)
	if err != nil {
		// TODO: Write S3-compatible error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ETag header on success.
	w.Header().Set("ETag", `"`+etag+`"`)
	w.WriteHeader(http.StatusOK)
}
