package s3

import (
	"net/url"
	"path"
)

type S3Repository struct {
	endpoint string
    bucket string
}

func NewS3Repository(endpoint, bucket string) *S3Repository {
	return &S3Repository{
		endpoint: endpoint,
        bucket: bucket,
	}
}

func (r *S3Repository) GetFileURL(key, bucket string) (string, error) {
	// Правильное объединение URL-путей
	u, _ := url.Parse(r.endpoint)
	u.Path = path.Join(bucket, key)
	return u.String()
}

// Для приватного бакета
// func (r *S3Repository) GetFileURL(ctx context.Context, key string) (string, error) {
// 	presignClient := s3.NewPresignClient(r.client)

// 	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
// 		Bucket: aws.String(r.bucket), // Используем сохраненный bucket
// 		Key:    aws.String(key),
// 	}, func(opts *s3.PresignOptions) {
// 		opts.Expires = 24 * time.Hour
// 	})

// 	if err != nil {
// 		return "", err
// 	}

// 	return req.URL, nil
// }