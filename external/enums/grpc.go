package enums

//go:generate go-enum -type=CompressionType
// CompressionType is an enumeration of GRPC traffic compression type values
type CompressionType int

const (
	None CompressionType = iota // default compression type (if not use compression)
	Gzip                        // gzip compression type
	Zstd                        // zstd compression type
)
