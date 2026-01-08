package chunkstore

type Config struct {
	// BucketName is the name of the bucket that stores blobs (S3 or OCI).
	BucketName string `docs:"required"`
	// Backend is the backend to use for object storage (s3 or oci).
	Backend string
}
