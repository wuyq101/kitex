package compress

// zstd  zlib gzip sa
type Compressor interface {
	// Compress writes the data written to wc to w after compressing it.  If an
	// error occurs while initializing the compressor, that error is returned
	// instead.
	Compress(dst, src []byte) []byte
	// Decompress reads data from r, decompresses it, and provides the
	// uncompressed data via the returned io.Reader.  If an error occurs while
	// initializing the decompressor, that error is returned instead.
	Decompress(src []byte) (int, error)
	// Name is the name of the compression codec and is used to set the content
	// coding header.  The result must be static; the result cannot change
	// between calls.

	DecompressedSize(src []byte) (int, error)

	Name() string
	// EXPERIMENTAL: if a Compressor implements
	// DecompressedSize(compressedBytes []byte) int, gRPC will call it
	// to determine the size of the buffer allocated for the result of decompression.
	// Return -1 to indicate unknown size.
}
