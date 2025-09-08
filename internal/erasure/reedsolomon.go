package erasure

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

const (
	// DefaultDataShards is the number of data shards for EC.
	DefaultDataShards = 4
	// DefaultParityShards is the number of parity shards for EC.
	DefaultParityShards = 2
)

// Encoder wraps the Reed-Solomon library for encoding operations.
type Encoder struct {
	enc reedsolomon.Encoder
}

// NewEncoder creates a new erasure encoder with default settings.
func NewEncoder() (*Encoder, error) {
	enc, err := reedsolomon.New(DefaultDataShards, DefaultParityShards)
	if err != nil {
		return nil, err
	}
	return &Encoder{enc: enc}, nil
}

// Encode reads from the reader, splits the data into shards, and returns them.
func (e *Encoder) Encode(r io.Reader) ([][]byte, error) {
	// Read all data from the reader into a buffer.
	// Note: For very large objects, this is memory-inefficient. A real-world
	// implementation would use streaming and process the object in chunks.
	// For this phase, we'll buffer it for simplicity.
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Split the buffered data into shards.
	shards, err := e.enc.Split(buf)
	if err != nil {
		return nil, err
	}
	// Encode the parity shards.
	err = e.enc.Encode(shards)
	if err != nil {
		return nil, err
	}
	return shards, nil
}

// Reconstruct verifies and reconstructs the data shards.
// It can tolerate up to `DefaultParityShards` missing (nil) shards.
func (e *Encoder) Reconstruct(shards [][]byte) error {
	// This reconstructs the data in-place in the shards slice.
	return e.enc.Reconstruct(shards)
}

// Join writes the reconstructed data shards to a writer.
func (e *Encoder) Join(w io.Writer, shards [][]byte) error {
	return e.enc.Join(w, shards, len(shards[0])*DefaultDataShards)
}
