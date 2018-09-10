package env

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func downloadFileFromS3(svc *s3.S3, bucket string, key string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func SetEnvFromS3(s3URL string) error {
	data, err := readEnv(s3URL)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		pieces := strings.SplitN(line, "=", 2)
		if len(pieces) > 1 {
			os.Setenv(pieces[0], pieces[1])
		}
	}
	return nil
}

func readEnv(s3URL string) ([]byte, error) {
	u, err := url.Parse(s3URL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "s3" {
		bucket := u.Host
		path := u.Path

		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{Region: aws.String("us-east-1")},
		})
		if err != nil {
			return nil, err
		}

		svc := s3.New(sess)

		data, err := downloadFileFromS3(svc, bucket, path)
		if err != nil {
			return nil, fmt.Errorf("Error downloading env file from S3: %s", err.Error())
		}
		return data, nil
	} else {
		return nil, fmt.Errorf("Unknown URL scheme %s for -env-file", s3URL)
	}
}
