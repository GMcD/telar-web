package handlers

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
)

func ConnectAws() (*session.Session, error) {
	accessKeyID := os.Getenv("ASSET_ACCESS_ID")
	secretAccessKey := os.Getenv("ASSET_ACCESS_SECRET")
	myRegion := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(myRegion),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})

	if err != nil {
		log.Error("Error getting s3 session : %s\n", err)
	}

	return sess, err
}

// Upload File to Bucket with key ObjectName
func UploadImage(c *fiber.Ctx, fileHeader *multipart.FileHeader, objectName string) (string, error) {
	// var buf bytes.Buffer

	sess := c.Locals("aws").(*session.Session)

	uploader := s3manager.NewUploader(sess)

	myBucket := os.Getenv("ASSET_BUCKET")
	myDomain := os.Getenv("ASSET_HOST")

	//get bytes into buffer for reader
	file, err := fileHeader.Open()
	// _, _ = io.Copy(&buf, )

	//upload to the s3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(objectName),
		Body:   file, // bytes.NewReader(buf.Bytes()),
	})

	if err != nil {
		return "", err
	}
	filepath := fmt.Sprintf("https://%s.%s/%s", myBucket, myDomain, file)
	return filepath, nil
}
