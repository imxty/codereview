package product

import (
	"context"
	gerr "errors"
	"time"

	tt "github.com/jinmukeji/huimaibao-service/pkg/time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

func (s *ProductAPIHandler) GetProductStatistics(ctx context.Context, req *pb.GetProductStatisticsRequest, rsp *pb.GetProductStatisticsResponse) error {
	// 验证request
	err := validateGetProductStatisticsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取商品
	product, err := s.productStore.GetProduct(ctx, req.GetProductId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if product == nil {
		return errors.Errorf(ErrProductNotFound, "product[%s] not found", req.GetProductId())
	}
	if product.GetOrganizationID() != req.GetOrganizationId() {
		return errors.Errorf(codes.InvalidRequest, "product[%s] not belong to organization[%s]", req.GetProductId(), req.GetOrganizationId())
	}
	// 只获取本月的商品统计
	now := time.Now().In(tt.LocBeijing)
	// 获取本月开始时间和结束时间
	// [本月的开始，和下个月的第一天的0点)
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, tt.LocBeijing)
	endTime := tt.GetLastDateOfMonth(now).AddDate(0, 0, 1)

	// 获取商品所有的曝光数据
	rs, err := s.productStore.GetProductStatisticsWithTimeRange(ctx, req.GetProductId(), startTime, endTime)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	totalCount, symptomExpose, productExpose := countProductExpose(rs)

	rsp.Product = toProtoProductFromRecommendedProduct(product)
	rsp.ProductExposure = productExpose
	rsp.SymptomExposure = symptomExpose
	rsp.ProductInstantTotalCount = totalCount

	return nil
}

// 验证request
func validateGetProductStatisticsRequest(req *pb.GetProductStatisticsRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization id should not be empty")
	}
	if req.GetProductId() == "" {
		return gerr.New("product id should not be empty")
	}
	return nil
}

// countProductExpose 根据时间从早到晚
func countProductExpose(rp []domain.RecommendedProductStatisticsIntf) (int32, []*pb.SymptomExposure, []*pb.ProductExposure) {
	// 没有曝光数据直接返回即可
	if len(rp) == 0 {
		return 0, nil, nil
	}
	expose := []*pb.ProductExposure{}
	// 获取曝光总数
	var totalCount int32 = int32(len(rp))
	// 获取不同病症曝光次数
	se := make(map[string]int32)
	for _, v := range rp {
		se[v.GetSymptom()]++
	}
	symptomExpose := make([]*pb.SymptomExposure, len(se))
	index := 0
	for k, v := range se {
		symptomExpose[index] = &pb.SymptomExposure{
			SymptomName:   k,
			ExposureCount: v,
		}
		index++
	}
	var exposeCount int32 = 0
	now := time.Now().In(tt.LocBeijing)
	// 获取第一笔数据的日期
	day := rp[0].GetExposedAt().In(tt.LocBeijing).Day()
	for _, v := range rp {
		rpDay := v.GetExposedAt().In(tt.LocBeijing).Day()
		if rpDay != day {
			expose = append(expose, &pb.ProductExposure{
				// 统计日期
				ExposedDate: &pb.Date{
					Year:  int32(now.Year()),
					Month: int32(now.Month()),
					Day:   int32(day),
				},
				// 某一天曝光次数
				ExposureCount: exposeCount,
			})
			// 日期设置为新的一天
			day = rpDay
			// 统计重新开始
			exposeCount = 0
		}
		exposeCount++
	}
	// 因为最后一天的数据无法append手动append
	expose = append(expose, &pb.ProductExposure{
		// 统计日期
		ExposedDate: &pb.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(day),
		},
		// 某一天曝光次数
		ExposureCount: exposeCount,
	})
	return totalCount, symptomExpose, expose
}
