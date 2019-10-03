package utils

import (
	"log"

	minio "github.com/minio/minio-go"
)

// MinioConnection Set up minio and initialize client
func MinioConnection() *minio.Client {
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

// MinioConnectionDEV Set up minio and initialize client
func MinioConnectionDEV() *minio.Client {
	// endpoint := "192.168.2.131:9000"
	endpoint := "localhost:9000"
	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient
}
