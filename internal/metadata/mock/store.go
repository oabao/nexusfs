package mock

import (
	"context"
	"fmt"

	"nexusfs/internal/metadata"
)

// Store is a mock implementation of the MetadataReadWriter interface for testing.
type Store struct {
	Buckets      map[string]metadata.Bucket
	Objects      map[string]metadata.Object
	Err          error
	CreatedBucket string
	WrittenObject *metadata.Object
}

// NewStore creates a new mock store.
func NewStore() *Store {
	return &Store{
		Buckets: make(map[string]metadata.Bucket),
		Objects: make(map[string]metadata.Object),
	}
}

func (s *Store) GetObjectMetadata(ctx context.Context, bucket, object string) (*metadata.Object, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	key := fmt.Sprintf("%s/%s", bucket, object)
	if obj, ok := s.Objects[key]; ok {
		return &obj, nil
	}
	return nil, nil // Not found
}

func (s *Store) WriteObjectMetadata(ctx context.Context, obj *metadata.Object) error {
	if s.Err != nil {
		return s.Err
	}
	s.WrittenObject = obj
	key := fmt.Sprintf("%s/%s", obj.Bucket, obj.Name)
	s.Objects[key] = *obj
	return nil
}

func (s *Store) ListObjects(ctx context.Context, bucketName string) ([]metadata.Object, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	var result []metadata.Object
	for _, obj := range s.Objects {
		if obj.Bucket == bucketName {
			result = append(result, obj)
		}
	}
	return result, nil
}

func (s *Store) DeleteObjectMetadata(ctx context.Context, bucketName, objectName string) error {
	if s.Err != nil {
		return s.Err
	}
	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	delete(s.Objects, key)
	return nil
}

func (s *Store) CreateBucket(ctx context.Context, bucket metadata.Bucket) error {
	if s.Err != nil {
		return s.Err
	}
	s.CreatedBucket = bucket.Name
	s.Buckets[bucket.Name] = bucket
	return nil
}

func (s *Store) GetBucket(ctx context.Context, name string) (*metadata.Bucket, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	if bucket, ok := s.Buckets[name]; ok {
		return &bucket, nil
	}
	return nil, nil // Not found
}

func (s *Store) ListBuckets(ctx context.Context) ([]metadata.Bucket, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	var result []metadata.Bucket
	for _, bucket := range s.Buckets {
		result = append(result, bucket)
	}
	return result, nil
}
