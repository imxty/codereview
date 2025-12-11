package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取商户实体信息版本请求
func (s *UserAPIHandler) GetEntityRevision(ctx context.Context, req *pb.GetEntityRevisionRequest, rsp *pb.GetEntityRevisionResponse) error {
	err := validateGetEntityRevisionRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	revision, err := s.userStore.GetTenantEntityRevision(ctx, req.GetRevisionId())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	if revision == nil {
		return errors.Errorf(codes.InvalidRequest, "revision[%s] not found", req.GetRevisionId())
	}

	rsp.Entity = toProtoTenantEntityFromRevision(revision)

	return nil
}

func validateGetEntityRevisionRequest(req *pb.GetEntityRevisionRequest) error {
	if req.GetRevisionId() == "" {
		return gerr.New("revision id should not be empty")
	}
	return nil
}
