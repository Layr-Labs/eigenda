package blobstore

type ObjectStorageBackend string

const (
	S3Backend  ObjectStorageBackend = "s3"
	OCIBackend ObjectStorageBackend = "oci"
)

type Config struct {
	BucketName string
	TableName  string
	Backend    ObjectStorageBackend
	// OCI-specific configuration
	OCINamespace     string
	OCIRegion        string
	OCICompartmentID string
}
