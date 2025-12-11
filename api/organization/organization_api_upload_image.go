package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest, rsp *pb.UploadImageResponse) error {
	err := validateUploadImageRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	uploadRsp, err := s.userAPI.UploadImage(ctx, &userv1.UploadImageRequest{
		Image: toSvcImage(req.GetImage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.ImageUrl = uploadRsp.GetImageUrl()
	return nil
}

// 验证request
func validateUploadImageRequest(req *pb.UploadImageRequest) error {
	if req.GetImage() == nil {
		return gerr.New("image should not be nil")
	}
	return nil
}
