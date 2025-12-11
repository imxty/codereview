package android

import (
	"encoding/json"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	// apk的描述信息文件名
	ApkComment = "apk_comment.json"
)

// ApkS3Client 获取安卓Apk信息的client
type ApkS3Client interface {
	GetApkInfo() (*ApkInfo, error)
}

// apk client
type ApkClient struct {
	options *ApkOptions
	sess    *session.Session
	s3      *s3.S3
}

// NewClient 返回一个新的 aws 连接
func NewApkClient(opts *ApkOptions) (*ApkClient, error) {
	client := new(ApkClient)
	client.options = opts
	creds := credentials.NewStaticCredentials(client.options.AccessKeyID, client.options.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(client.options.Region),
	})
	if err != nil {
		return nil, err
	}
	client.sess = sess
	client.s3 = s3.New(sess)
	return client, nil
}

// 获取apk的信息
func (client *ApkClient) GetApkInfo() (*ApkInfo, error) {
	// 获取apk comment
	comment, err := client.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(client.options.BucketName),
		Key:    aws.String(client.options.ApkPath),
	})
	if err != nil {
		return nil, err
	}
	defer comment.Body.Close() // nolint: errcheck
	// 读取comment
	comm, err := io.ReadAll(comment.Body)
	if err != nil {
		return nil, err
	}
	// 解析comment
	apkInfo := new(ApkInfo)
	err = json.Unmarshal(comm, apkInfo)
	if err != nil {
		return nil, err
	}

	return apkInfo, nil
}
