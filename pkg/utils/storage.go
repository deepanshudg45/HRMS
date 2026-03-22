package utils

import (
	"context"
	"fmt"
	"bytes"


	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"path/filepath"
)

// 🔹 Interface
type StorageService interface {
	Upload(ctx context.Context, key string, data []byte) error
	Delete(ctx context.Context, key string) error
	URL(key string) string
}





type LocalAdapter struct {
	BasePath string
	BaseURL  string
}

func (l *LocalAdapter) Upload(ctx context.Context, key string, data []byte) error {
	path := filepath.Join(l.BasePath, key)
	return os.WriteFile(path, data, 0644)
}

func (l *LocalAdapter) Delete(ctx context.Context, key string) error {
	path := filepath.Join(l.BasePath, key)
	return os.Remove(path)
}

func (l *LocalAdapter) URL(key string) string {
	return fmt.Sprintf("%s/%s", l.BaseURL, key)
}




type S3Adapter struct {
	Client *s3.Client
	Bucket string
	Region string
}

func (s *S3Adapter) Upload(ctx context.Context, key string, data []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3Adapter) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	return err
}

func (s *S3Adapter) URL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.Bucket, s.Region, key)
}