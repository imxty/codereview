package product

import (
	"context"
	gerr "errors"
	"slices"
	"sort"
	"strings"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// 推荐商品数量上限
	productsAmountLimit = 6
	// 理疗湿气血瘀轻度标准
	physicalTherapyMildStandard = 3
)

func (s *ProductAPIHandler) ListRecommendedProducts(ctx context.Context, req *pb.ListRecommendedProductsRequest, rsp *pb.ListRecommendedProductsResponse) error {
	// 验证request
	err := validateListRecommendedProductsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 获取租户商品推荐状态
	getTenantRsp, err := s.userAPI.GetTenantEntity(ctx, &userpb.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return err
	}

	// 获取组织ID
	organizationID := getTenantRsp.GetEntity().GetOrganizationId()

	// 报告推荐方案版本id,如果存在返回以前的即可
	tenantTreatmentRevIDFromReport := req.GetTenantTreatmentRevId()
	// 判断应急态情况
	// 脏腑辩证和风险疾病为空 表明应急态
	if req.GetDirtyDialecticReportResult() == nil &&
		req.GetHighRiskyDiseaseReportResult() == nil &&
		req.GetMediumRiskyDiseaseReportResult() == nil {
		rsp.Products = nil
		rsp.TenantTreatmentRevId = ""
		rsp.AreLatestTreatments = true
		return nil
	}

	// 获取最新版方案版本id
	listTreatmentsRsp, err := s.userAPI.ListTreatmentsByTenantIDs(ctx, &userpb.ListTreatmentsByTenantIDsRequest{
		TenantIds: []string{req.GetTenantId()},
	})
	if err != nil {
		return err
	}
	// 获取绑定商户的方案
	treatmentIds := listTreatmentsRsp.GetTreatments()

	// 报告不是初次生成
	if !req.GetIsFirstReported() {
		// 方案版本id都为空 返回空
		if tenantTreatmentRevIDFromReport == "" {
			rsp.TenantTreatmentRevId = ""
			rsp.Products = nil
			rsp.AreLatestTreatments = true
			return nil
		}

		areLatestTreatments := false
		// 如果不存在绑定的方案
		tid, ok := treatmentIds[req.GetTenantId()]
		// 如果有采取匹配
		if ok {
			// 获取方案
			ts, err := s.productStore.GetTreatment(ctx, tid)
			if err != nil {
				return errors.Error(codes.DataAccessFailed, err.Error())
			}
			tenantLatestRevId := ts.GetLatestRev()
			// 判断最新的方案ID和报告记录的方案id是否一致，如果一致则是最新的
			if tenantLatestRevId == req.GetTenantTreatmentRevId() {
				areLatestTreatments = true
			}
		}

		// 得到商品列表
		pProducts, err := s.listRecommendedProductsByTreatmentRevId(ctx, organizationID, tenantTreatmentRevIDFromReport, req)
		if err != nil {
			return err
		}

		rsp.Products = pProducts
		rsp.AreLatestTreatments = areLatestTreatments
		rsp.TenantTreatmentRevId = tenantTreatmentRevIDFromReport
		return nil
	}

	// 如果是初次生成 获取 租户的最新方案tenantLatestRevId
	// 如果不存在绑定的方案，而且直接返回即可
	tid, ok := treatmentIds[req.GetTenantId()]
	if !ok {
		rsp.AreLatestTreatments = true
		return nil
	}
	// 获取方案
	ts, err := s.productStore.GetTreatment(ctx, tid)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	tenantLatestRevId := ts.GetLatestRev()
	// 得到商品列表
	pProducts, err := s.listRecommendedProductsByTreatmentRevId(ctx, organizationID, tenantLatestRevId, req)
	if err != nil {
		return err
	}

	rsp.TenantTreatmentRevId = tenantLatestRevId
	// 首次获取一定是最新
	rsp.AreLatestTreatments = true
	rsp.Products = pProducts

	return nil
}

// 验证request
func validateListRecommendedProductsRequest(req *pb.ListRecommendedProductsRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("invalid tenant_id")
	}
	if req.GetReportId() == "" {
		return gerr.New("invalid report_id")
	}
	return nil
}

// 得到商品列表
func (p *ProductAPIHandler) listRecommendedProductsByTreatmentRevId(ctx context.Context, organizationID, tenantTreatmentRevId string, req *pb.ListRecommendedProductsRequest) ([]*pb.Product, error) {
	// 查询商品
	products, err := p.productStore.ListTreatmentRevItems(ctx, tenantTreatmentRevId)
	if err != nil {
		return nil, errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 查到的associations长度为0
	if len(products) == 0 {
		return []*pb.Product{}, nil
	}
	// 将association转换成lookups
	// 一个症候可能对应多个商品
	recommendProductMap := make(map[string]domain.TreatmentProductAssociationIntf)
	tenantLookups := make(map[string][]string)
	for _, v := range products {
		tenantLookups[v.GetSymptom()] = append(tenantLookups[v.GetSymptom()], v.GetRecommendedProductID())
		// 构建商品map
		recommendProductMap[v.GetRecommendedProductID()] = v
	}

	// 按规则过滤商品并排序 获取有序的商品id列表 和 商品id症候map
	productsId, productSymptomMap, err := filterRecommendedProducts(tenantLookups, req)
	if err != nil {
		return nil, err
	}
	// 对商品进行排序，组合商品优先
	sort.Slice(productsId, func(i, j int) bool {
		return len(productsId[i]) > len(productsId[j])
	})

	// 报告初次生成 需要将报告出现的商品和统计数据存入数据库
	if req.GetIsFirstReported() {
		// 将报告中的商品和症候记录到数据库
		err = p.recordIntoDbStore(ctx, productSymptomMap, req.GetReportId(), organizationID, req.GetTenantId())
		if err != nil {
			return nil, err
		}
	}

	// 获取商品列表
	pProducts := []*pb.Product{}
	for _, pids := range productsId {
		// 获取对应商品的symptom
		productsKey := strings.Join(pids, ",")
		symptomKey := productSymptomMap[productsKey]
		// 构建商品
		for _, v := range pids {
			// 非空判断
			if recommendProductMap[v] != nil {
				pProducts = append(pProducts, toProtoProductFromTreatmentProductAssociationWithSymptom(recommendProductMap[v], symptomKey))
			}
		}
	}
	return pProducts, nil
}

// 按规则过滤推荐商品
// 商品推荐最多推荐6个药品，推荐顺序按照：
//  1. 中医报告脏腑辨证结果
//  2. 西医报告高风险疾病
//  3. 西医报告中风险疾病
//  4. 理疗报告湿气、血瘀
//  5. 轻度湿气不推药
//  6. 体质
//     一种症候推荐一个商品
func filterRecommendedProducts(lookups map[string][]string, req *pb.ListRecommendedProductsRequest) ([][]string, map[string]string, error) {
	if lookups == nil {
		return nil, nil, gerr.New("lookups should not be nil")
	}

	var dd, hrd, mrd, ptDeleted, symptomList []string
	var dderr, hrderr, mrderr error
	// 获取四种报告结果的症候列表
	// 脏腑辨证症候key
	dd, dderr = getSymptomList(req.GetDirtyDialecticReportResult())
	// 高风险疾病症候key
	hrd, hrderr = getSymptomList(req.GetHighRiskyDiseaseReportResult())
	// 中风险疾病症候key
	mrd, mrderr = getSymptomList(req.GetMediumRiskyDiseaseReportResult())
	if dderr != nil || hrderr != nil || mrderr != nil {
		return nil, nil, gerr.New("report result should not be nil")
	}
	// 去掉轻度湿气
	for _, v := range req.GetPhysicalTherapyReportResult() {
		if v.GetScore() > physicalTherapyMildStandard {
			ptDeleted = append(ptDeleted, v.GetSymptomKey())
		}
	}
	// 体质推药
	pd, err := getSymptomList([]*pb.ReportResult{req.GetPhysiqueDialecticsReportResult()})
	if err != nil {
		return nil, nil, gerr.New("invalid physical dialectics")
	}

	// 按规则排列症候报告结果
	resultList := [][]string{dd, hrd, mrd, ptDeleted, pd}
	for _, v := range resultList {
		symptomList = append(symptomList, v...)
	}

	// 根据症候列表[]string 遍历lookups获取对应症候下的一种商品
	var filterProductList [][]string
	// map[商品id]症候key
	productSymptomMap := make(map[string]string)
	// 去重商品(多个症候对应一个商品的情况)
	book := make(map[string]bool)
	for _, v := range symptomList {
		if len(filterProductList) >= productsAmountLimit {
			break
		}
		if products, ok := lookups[v]; ok {
			// 对products排序
			slices.Sort(products)
			// 构建key
			productsKey := strings.Join(products, ",")
			// 判断是否重复
			if !book[productsKey] {
				// 记录商品对应的symptom
				productSymptomMap[productsKey] = v
				filterProductList = append(filterProductList, products)
				book[productsKey] = true
				// 商品ID去重，如果组合商品包含了多个商品，需要对单个商品去重，如果之前已经存在则需要从filterProductList剔除
				// 只有组合商品会影响单个商品
				if len(products) > 1 {
					for _, pid := range products {
						book[pid] = true
						// 过滤ProductList
						filterProductList = filterProducts(filterProductList, pid)
					}
				}
			}
		}
	}

	return filterProductList, productSymptomMap, nil
}

// filterProducts
func filterProducts(s [][]string, goal string) [][]string {
	for index, v := range s {
		// 如果只有一个且正好和目标一样，即单个商品又在组合商品中出现，需要剔除
		if len(v) == 1 && v[0] == goal {
			return append(s[:index], s[index+1:]...)
		}
	}
	return s
}

func getSymptomList(rr []*pb.ReportResult) ([]string, error) {
	if rr == nil {
		return []string{}, nil
	}
	symptoms := make([]string, len(rr))
	for i, v := range rr {
		symptoms[i] = v.GetSymptomKey()
	}
	return symptoms, nil
}

// 曝光商品和症候记录到数据库
func (p *ProductAPIHandler) recordIntoDbStore(ctx context.Context, productSymptomMap map[string]string, reportId string, organizationID, tenantId string) error {
	if len(productSymptomMap) == 0 {
		return nil
	}
	productStatistics := []domain.RecommendedProductStatisticsIntf{}
	for idString, symptomKey := range productSymptomMap {
		// id可能为多个
		ids := strings.Split(idString, ",")
		for _, id := range ids {
			productStatistics = append(productStatistics, toDomainRecommendedProductStatistics(id, reportId, symptomKey, organizationID, tenantId))
		}
	}
	// 记录到数据库
	err := p.productStore.CreateRecommendedProductStandingBook(ctx, productStatistics)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	return nil
}

// 转换 toDomainRecommendedProductStatistics
func toDomainRecommendedProductStatistics(productId, reportId, symptom, organizationID, tenantID string) *domain.RecommendedProductStatistics {

	return &domain.RecommendedProductStatistics{
		RecommendedProductStatisticsID: xid.New().String(),
		OrganizationID:                 organizationID,
		TenantID:                       tenantID,
		RecommendedProductID:           productId,
		ReportID:                       reportId,
		Symptom:                        symptom,
		ExposedAt:                      time.Now().UTC(),
	}
}
