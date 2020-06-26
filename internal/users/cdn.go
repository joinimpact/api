package users

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joinimpact/api/internal/config"
)

type cdnClient struct {
	config  *config.Config
	session *session.Session
}

func newCDNClient(config *config.Config) *cdnClient {
	key := config.CDNKey
	secret := config.CDNSecret

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(config.CDNEndpoint),
		Region:      aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("connected to s3")

	return &cdnClient{
		config,
		newSession,
	}
}

// uploadImage uploads an image to the Spaces CDN. On success, it returns the full URL of the profile picture.
func (c *cdnClient) uploadImage(imageName string, reader io.Reader) (string, error) {
	uploader := s3manager.NewUploader(c.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(c.config.CDNBucket),
		Key:         aws.String(imageName),
		Body:        reader,
		ContentType: aws.String("image/png"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://cdn.joinimpact.org/%s", imageName), nil
}
