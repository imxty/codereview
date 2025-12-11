package product

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 通过treatmentId获取名称
func (s *ProductAPIHandler) ListTreatmentNameByIds(ctx context.Context, req *pb.ListTreatmentNameByIdsRequest, rsp *pb.ListTreatmentNameByIdsResponse) error {
	err := validateListTreatmentNameByIdsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取方案名称
	revs, err := s.productStore.ListTreatmentRevs(ctx, req.GetTreatmentIds())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 返回名称
	tName := make(map[string]string)
	for _, v := range revs {
		tName[v.GetTreatmentID()] = v.GetTreatmentName()
	}
	rsp.TreatmentNames = tName

	return nil
}

func validateListTreatmentNameByIdsRequest(req *pb.ListTreatmentNameByIdsRequest) error {
	if len(req.GetTreatmentIds()) == 0 {
		return gerr.New("treatment id should not be nil")
	}
	return nil
}
