package gateway

import "net/http"

// S3Api is responsible for handling the S3-compatible HTTP requests.
type S3Api struct {
	service *Service
}

func NewS3Api(service *Service) *S3Api {
	return &S3Api{service: service}
}

// GetObjectHandler is the HTTP handler for GET Object requests.
func (api *S3Api) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse bucket and object from request.
	// Call api.service.GetObject(...)
	// Translate service response to an S3-compatible HTTP response.
}

// PutObjectHandler is the HTTP handler for PUT Object requests.
func (api *S3Api) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement S3 PUT object logic.
}
