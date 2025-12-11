package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 上传图片
func (u *UserAPIHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest, rsp *pb.UploadImageResponse) error {
	// 1.验证request
	err := validateUploadImageRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	imageUrl, err := u.uploadImage(req.GetImage())
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.ImageUrl = imageUrl

	return nil
}

// 验证request
func validateUploadImageRequest(req *pb.UploadImageRequest) error {
	if req.GetImage() == nil {
		return gerr.New("image should not be nil")
	}
	return nil
}
