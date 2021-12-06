package enums

//go:generate go-enum -type=CompressionType
// CompressionType is an enumeration of GRPC traffic compression type values
type CompressionType int

const (
	None CompressionType = iota // default compression type (if not use compression)
	Gzip                        // gzip compression type
	Zstd                        // zstd compression type
)

//go:generate go-enum -type=SSLMode -transform=lower
// SSLMode is type of sslmode postgresql connection
type SSLMode int

const (
	Enable  SSLMode = iota // SSLMode postgres connection string sslmode Enable
	Disable                // SSLMode postgres connection string sslmode Disable
)
