package metadata

import "context"

import "time"

// MetadataReadWriter is the abstraction that high-level modules
// will depend on for all metadata operations. It defines the contract
// for interacting with any underlying metadata storage system.
type MetadataReadWriter interface {
	// Object methods
	GetObjectMetadata(ctx context.Context, bucket, object string) (*Object, error)
	WriteObjectMetadata(ctx context.Context, obj *Object) error
	ListObjects(ctx context.Context, bucketName string) ([]Object, error)
	DeleteObjectMetadata(ctx context.Context, bucketName, objectName string) error

	// Bucket methods
	CreateBucket(ctx context.Context, bucket Bucket) error
	GetBucket(ctx context.Context, name string) (*Bucket, error)
	ListBuckets(ctx context.Context) ([]Bucket, error)
}

// Bucket represents the metadata for a bucket.
type Bucket struct {
	Name    string
	Created time.Time
}

// Object represents the metadata for a stored object. This is the
// data structure that is passed across the abstraction boundary.
type Object struct {
	Name           string
	Bucket         string
	Size           int64
	ETag           string
	ShardLocations map[string]string // e.g., { "shard_1": "storage_node_A:/path/..." }
}
