package testutils

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/exp/rand"
)

func RandStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandBytes(n int) []byte {
	return []byte(RandStr(n))
}

// Panics if the bucket does not exist
func RemoveBlobInfoFromBucket(bucketName string, blobInfo []byte) error {
	// Initialize minio client object.
	endpoint := minioEndpoint
	accessKeyID := minioAdmin
	secretAccessKey := minioAdmin
	useSSL := false
	minioClient, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
	// Panic, the bucket should already exist
	if err != nil {
		panic(err)
	}
	key := crypto.Keccak256(blobInfo[1:])
	objectName := hex.EncodeToString(key)
	ctx := context.Background()
	err = minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("Successfully removed %s from %s\n", objectName, bucketName))

	return nil
}

// Panics if the bucket does not exist
func ExistsBlobInfotInBucket(bucketName string, blobInfo []byte) (bool, error) {
	// Initialize minio client object.
	endpoint := minioEndpoint
	accessKeyID := minioAdmin
	secretAccessKey := minioAdmin
	useSSL := false
	minioClient, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
	// Panic, the bucket should already exist
	if err != nil {
		panic(err)
	}
	key := crypto.Keccak256(blobInfo[1:])
	objectName := hex.EncodeToString(key)
	ctx := context.Background()
	_, err = minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
