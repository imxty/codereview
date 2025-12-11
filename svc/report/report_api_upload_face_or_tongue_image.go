package report

import (
	"bytes"
	"context"
	gerr "errors"
	"fmt"
	"io"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/image"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

// UploadFaceOrTongueImage 上传人脸或舌图
func (s *ReportAPIHandler) UploadFaceOrTongueImage(ctx context.Context, req *pb.UploadFaceOrTongueImageRequest, rsp *pb.UploadFaceOrTongueImageResponse) error {
	err := validateUploadFaceOrTongueImageRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	imageUrl, err := s.uploadImage(req.GetImage())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.ImageUrl = imageUrl
	return nil
}

// 验证request
func validateUploadFaceOrTongueImageRequest(req *pb.UploadFaceOrTongueImageRequest) error {
	if req.GetImage() == nil {
		return gerr.New("image should not be nil")
	}
	return nil
}

// 上传图片
func (u *ReportAPIHandler) uploadImage(logo *pb.UploadingImage) (string, error) {
	const (
		ImageSizeLimit = 4 * 1024 * 1024
	)
	logoUrl := ""
	if logo == nil {
		return "", nil
	}
	// 检查 mime 和大小
	if len(logo.GetImage()) > ImageSizeLimit {
		return "", gerr.New("image should not exceed 4M")
	}
	// 检测 mime
	suffix, err := image.GetImageSuffix(logo.GetMime())
	if err != nil {
		return "", err
	}
	// 如果传了 logo 就去上传 logo
	// 生成图片的 ID
	logoUrl = fmt.Sprintf("%s.%s", xid.New().String(), suffix)
	path, err := u.s3Store.Save(logoUrl, logo.GetMime(), byteReader(logo.GetImage()))
	if err != nil {
		return "", err
	}
	return path, nil
}

// byteReader 将字节格式化为可读文本流
func byteReader(img []byte) io.ReadSeeker {
	return bytes.NewReader(img)
}
