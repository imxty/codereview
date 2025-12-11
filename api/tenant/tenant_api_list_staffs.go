package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *TenantAPIHandler) ListStaffs(ctx context.Context, req *pb.ListStaffsRequest, rsp *pb.ListStaffsResponse) error {
	err := validateListStaffsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 发送获取员工列表请求
	listRsp, err := s.userAPI.ListStaffs(ctx, &userpb.ListStaffsRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	rsp.ActivatedStaffCount = listRsp.GetActivatedStaffCount()
	rsp.MaxStaffCount = listRsp.GetMaxStaffCount()
	pStaffs := make([]*pb.Staff, len(listRsp.GetStaffs()))
	for k, v := range listRsp.GetStaffs() {
		pStaffs[k] = toAppStaff(v)
	}
	rsp.Staffs = pStaffs
	return nil
}

// 验证request
func validateListStaffsRequest(req *pb.ListStaffsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("staff id should not be empty")
	}
	return nil
}
