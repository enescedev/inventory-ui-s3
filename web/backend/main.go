package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"backend/handlers"
)

func initS3() (*minio.Client, string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, bucket, err
	}

	exists, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return nil, bucket, err
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, bucket, err
		}
	}
	return client, bucket, nil
}

func main() {
	client, bucket, err := initS3()
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}

	h := &handlers.TableHandler{Client: client, Bucket: bucket}
	r := mux.NewRouter()
	r.HandleFunc("/table", h.GetTable).Methods(http.MethodGet)
	r.HandleFunc("/table", h.PutTable).Methods(http.MethodPut)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
