package dataplane

// S3Client is a convenience wrapper for uploading and downloading files from amazon S3. May
// break down files into smaller parts for upload (to improve latency), and if so the files are
// reassembled on download. This tool is not intended to be used for reading and writing files
// that are consumed by utilities that are not aware of the multipart upload/download process.
//
// Implementations of this interface are required to be thread-safe.
type S3Client interface {
	// Upload uploads a file to S3. The fragmentSize parameter specifies the maximum size of each
	// file uploaded to S3. If the file is larger than fragmentSize then it will be broken into
	// smaller parts and uploaded in parallel. The file will be reassembled on download.
	Upload(key string, data []byte, fragmentSize int) error
	// Download downloads a file from S3, as written by Upload. The fileSize (in bytes) and fragmentSize
	// must be the same as the values used in the Upload call.
	Download(key string, fileSize int, fragmentSize int) ([]byte, error)
	// Close closes the S3 client.
	Close() error
}
