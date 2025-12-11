package s3

import (
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinmukeji/huimaibao-service/pkg/filestore"
)

// s3Store 与 aws 通信
type s3Store struct {
	options *Options
	sess    *session.Session
	s3      *s3.S3
}

// NewS3Store 返回一个新的 aws 连接
func NewS3Store(opts ...Option) (filestore.FileStore, error) {
	s3Store := new(s3Store)
	s3Store.options = newOptions(opts...)
	creds := credentials.NewStaticCredentials(s3Store.options.AccessKeyID, s3Store.options.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(s3Store.options.Region),
	})
	if err != nil {
		return nil, err
	}
	s3Store.sess = sess
	s3Store.s3 = s3.New(sess)
	return s3Store, nil
}

// uploadFile 存储一个输入流到指定路径
func (store *s3Store) uploadFile(readerSeeker io.ReadSeeker, key string, acl string, mime string) (*s3.PutObjectOutput, error) {
	svc := store.s3
	return svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(store.options.BucketName),
		Key:         aws.String(key),
		Body:        readerSeeker,
		ACL:         aws.String(acl),
		ContentType: &mime,
	})
}

// Save 上传静态资源池资源
func (store *s3Store) Save(path, mime string, content io.ReadSeeker) (string, error) {
	s3Path := store.buildFullKey(path)
	_, err := store.uploadFile(content, s3Path, awss3.ObjectCannedACLPrivate, mime)
	if err != nil {
		return "", err
	}
	return s3Path, nil
}

// buildFullKey 构建 s3 存储时的完整的 key
func (store *s3Store) buildFullKey(key string) string {
	return path.Join(store.options.KeyPrefix, key)
}
