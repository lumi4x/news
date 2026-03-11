package filestorage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lumi4x/news/src/common/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var Default *FileStorage

type FileStorage struct {
	client     *s3.S3
	session    *session.Session
	bucketName string
	endpoint   string
}

const (
	MinIONotFound = "NotFound"
)

func init() {
	filestorage, err := new(config.FileStorageBucketName, config.FileStorageAccessKey, config.FileStorageSecretKey, config.FileStorageEndpoint)
	if err != nil {
		logrus.
			WithError(err).
			WithField("FileStorageBucketName", config.FileStorageBucketName).
			WithField("FileStorageAccessKey", config.FileStorageAccessKey).
			WithField("FileStorageSecretKey", config.FileStorageSecretKey).
			WithField("FileStorageEndpoint", config.FileStorageEndpoint).
			Panicf("failed initializing file storage")
	}

	Default = filestorage
}

func new(bucketName, accessKey, secretKey, endpoint string) (*FileStorage, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(true), // TODO: Remove
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed creating s3 session")
	}

	svc := s3.New(sess)

	// Check if the bucket already exists
	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// If the bucket doesn't exist, create it
		if aerr, ok := err.(awserr.Error); ok && (aerr.Code() == s3.ErrCodeNoSuchBucket || aerr.Code() == MinIONotFound) {
			_, err = svc.CreateBucket(&s3.CreateBucketInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed creating bucket")
			}
		} else {
			return nil, errors.Wrapf(err, "failed getting bucket")
		}
	}

	return &FileStorage{client: svc, session: sess, bucketName: bucketName, endpoint: endpoint}, nil
}

func (s *FileStorage) exists(objectKey string) (bool, error) {
	// Check if the file already exists in the bucket
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		// If the error is a NoSuchKey error, the object doesn't exist
		if awsErr, ok := err.(awserr.Error); ok {
			if s.isErrorNotFound(err) {
				return false, nil
			}

			logrus.Info(awsErr.Code())
		}
		// If the error is not a NoSuchKey error, return it
		return false, errors.Wrapf(err, "failed to check if object exists, code")
	}

	return true, nil
}

func (s *FileStorage) isErrorNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == s3.ErrCodeNoSuchKey || awsErr.Code() == MinIONotFound {
			return true
		}
	}

	return false
}

func (s *FileStorage) presign(objectKey string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})

	urlStr, err := req.Presign(expiration)
	if err != nil {
		return "", errors.Wrapf(err, "failed presigning filestorage object")
	}

	return urlStr, nil
}

func (s *FileStorage) Upload(directory string, r io.ReadSeeker) (string, error) {
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return "", errors.Wrapf(err, "failed seeking file to start")
	}

	hash := sha256.New()
	_, err = io.Copy(hash, r)
	if err != nil {
		return "", fmt.Errorf("failed to calculate hash: %v", err)
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek reader")
	}

	hashSum := hex.EncodeToString(hash.Sum(nil))
	objectKey := fmt.Sprintf("%s/%s", directory, hashSum)

	exists, err := s.exists(objectKey)
	if err != nil {
		return "", err
	}

	if exists {
		logrus.WithField("ObjectKey", objectKey).Info("Not uploading existing file")
		return objectKey, nil
	}

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
		Body:   r,
	})
	if err != nil {
		return "", err
	}

	logrus.WithField("ObjectKey", objectKey).Info("Uploaded file")
	return objectKey, nil
}

func (s *FileStorage) Download(objectKey string, w io.WriterAt, ctx context.Context) (bool, error) {
	downloader := s3manager.NewDownloader(s.session)

	// Required to support HTTP sequential writing
	downloader.Concurrency = 1

	_, err := downloader.DownloadWithContext(ctx, w, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		if s.isErrorNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
