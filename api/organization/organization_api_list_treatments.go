package organization

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/organization/v1"
	productv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *OrganizationAPIHandler) ListTreatments(ctx context.Context, req *pb.ListTreatmentsRequest, rsp *pb.ListTreatmentsResponse) error {
	err := validateListTreatmentsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	listRsp, err := s.productAPI.ListTreatments(ctx, &productv1.ListTreatmentsRequest{
		// 组织id
		OrganizationId: req.GetOrganizationId(),
		IsPublish:      req.GetIsPublish(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	ts := listRsp.GetTreatments()
	pts := make([]*pb.Treatment, len(ts))
	for k, v := range ts {
		pts[k] = toAppTreatment(v, s.s3Domain)
	}
	rsp.Treatments = pts
	return nil
}

// 验证request
func validateListTreatmentsRequest(req *pb.ListTreatmentsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	return nil
}
