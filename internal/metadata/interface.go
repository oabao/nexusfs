package metadata

import "context"

// MetadataReadWriter is the abstraction that high-level modules
// will depend on for all metadata operations. It defines the contract
// for interacting with any underlying metadata storage system.
type MetadataReadWriter interface {
	GetObjectMetadata(ctx context.Context, bucket, object string) (*Object, error)
	WriteObjectMetadata(ctx context.Context, obj *Object) error
	// TODO: Add other necessary methods like DeleteObjectMetadata, ListObjects, etc.
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
