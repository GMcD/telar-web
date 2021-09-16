package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	keyFile := "/var/openfaas/secrets/media-access-key-id"
	secretFile := "/var/openfaas/secrets/media-secret-access-key"
	ak, akError := ioutil.ReadFile(keyFile)
	if akError != nil {
		return nil, akError
	}
	accessKeyID := string(ak)
	as, asError := ioutil.ReadFile(secretFile)
	if asError != nil {
		return nil, asError
	}
	secretAccessKey := string(as)

	log.Info("AWS Access Key %s, read from %s.\n", accessKeyID, keyFile)
	if accessKeyID == "" || secretAccessKey == "" {
		return nil, errors.New("Invalid AWS Access Credentials.")
	}

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

	//TODO: Move into handler setup?
	session, err := ConnectAws()
	if session == nil {
		return "", errors.New("S3 Session connection failed.")
	}

	uploader := s3manager.NewUploader(session)

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
	filepath := fmt.Sprintf("https://%s.%s/%s", myBucket, myDomain, objectName)
	return filepath, nil
}
