package main

import (
    "fmt"
    "os"
    "log"
    "encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
	}
}


func getS3FileData(bucket, key string) ([]byte, error) {
	region := os.Getenv("REGION")
	access_key := os.Getenv("MY_AWS_ACCESS_KEY_ID")
	secret := os.Getenv("MY_AWS_SECRET_ACCESS_KEY")

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(access_key, secret, ""),
	}))

	downloader := s3manager.NewDownloader(sess)
    write_buffer := aws.NewWriteAtBuffer([]byte{})

	result, err := downloader.Download(write_buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	fmt.Printf("File downloaded: %d bytes\n", result)
	fmt.Println(string(write_buffer.Bytes()))

	return write_buffer.Bytes(), nil
}

func isValidRequest(bucket, validation_key, access_token string) (bool, error) {
    file_data, err := getS3FileData(bucket, validation_key)
	if err != nil {
		return false, err
	}

	return string(file_data) == access_token, nil
}

func getMatchJSON(bucket, match_id string) (map[string]string, error) {
	var json_data map[string]string

    data, err := getS3FileData(bucket, match_id)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &json_data)
	fmt.Println(json_data["firstname"])

	return json_data, nil
}
