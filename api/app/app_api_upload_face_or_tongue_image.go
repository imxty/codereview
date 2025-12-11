package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) UploadFaceOrTongueImage(ctx context.Context, req *pb.UploadFaceOrTongueImageRequest, rsp *pb.UploadFaceOrTongueImageResponse) error {
	err := validateUploadFaceOrTongueImageRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送更新请求
	uploadRsp, err := s.reportAPI.UploadFaceOrTongueImage(ctx, &reportpb.UploadFaceOrTongueImageRequest{
		Image: toSvcImage(req.GetImage()),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.ImageUrl = uploadRsp.GetImageUrl()
	return nil
}

// 验证request
func validateUploadFaceOrTongueImageRequest(req *pb.UploadFaceOrTongueImageRequest) error {
	if req.GetImage() == nil {
		return gerr.New("image should not be empty")
	}
	return nil
}
