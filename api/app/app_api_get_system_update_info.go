package app

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *AppAPIHandler) GetSystemUpdateInfo(ctx context.Context, req *pb.GetSystemUpdateInfoRequest, rsp *pb.GetSystemUpdateInfoResponse) error {
	info, err := s.apkClient.GetApkInfo()
	if err != nil {
		return errors.Error(codes.InvalidOperation, err.Error())
	}

	rsp.Version = info.Version
	rsp.ApkLink = info.ApkLink
	rsp.ApkSize = info.ApkSize
	rsp.UpdateInfo = info.UpdateInfo
	return nil
}
