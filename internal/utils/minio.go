package utils

import (
	"log"

	minio "github.com/minio/minio-go"
)

// MinioConnection Set up minio and initialize client
func MinioConnection() *minio.Client {
	// endpoint := "localhost:9111"
	// accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	// secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	endpoint := "s3:9000"
	accessKeyID := "AKIAIOSFODNN7JASUINM"
	secretAccessKey := "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"

	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
