package s3

import "errors"

var (
	// Configuration errors
	ErrS3EmptyRegion    = errors.New("s3: empty region")
	ErrS3EmptyBucket    = errors.New("s3: empty bucket")
	ErrS3EmptyAccessKey = errors.New("s3: empty access key")
	ErrS3EmptySecretKey = errors.New("s3: empty secret key")
	ErrS3ConfigFailed   = errors.New("s3: config failed")

	// Operation errors
	ErrS3EmptyKey      = errors.New("s3: empty key")
	ErrS3EmptyData     = errors.New("s3: empty data")
	ErrS3EmptyLocation = errors.New("s3: empty location")

	// Upload/Download errors
	ErrS3UploadFailed   = errors.New("s3: upload failed")
	ErrS3DownloadFailed = errors.New("s3: download failed")
	ErrS3DeleteFailed   = errors.New("s3: delete failed")
	ErrS3ListFailed     = errors.New("s3: list failed")

	// Input validation errors
	ErrS3InvalidInput  = errors.New("s3: invalid input")
	ErrS3InternalError = errors.New("s3: internal error")
)
