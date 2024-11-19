package common

// BlobSizeBucket maps the blob size into a bucket that's defined according to
// the power of 2.
func BlobSizeBucket(blobSize int) string {
	switch {
	case blobSize <= 1*1024:
		return "1KiB"
	case blobSize <= 2*1024:
		return "2KiB"
	case blobSize <= 4*1024:
		return "4KiB"
	case blobSize <= 8*1024:
		return "8KiB"
	case blobSize <= 16*1024:
		return "16KiB"
	case blobSize <= 32*1024:
		return "32KiB"
	case blobSize <= 64*1024:
		return "64KiB"
	case blobSize <= 128*1024:
		return "128KiB"
	case blobSize <= 256*1024:
		return "256KiB"
	case blobSize <= 512*1024:
		return "512KiB"
	case blobSize <= 1024*1024:
		return "1MiB"
	case blobSize <= 2*1024*1024:
		return "2MiB"
	case blobSize <= 4*1024*1024:
		return "4MiB"
	case blobSize <= 8*1024*1024:
		return "8MiB"
	case blobSize <= 16*1024*1024:
		return "16MiB"
	case blobSize <= 32*1024*1024:
		return "32MiB"
	default:
		return "invalid"
	}
}
