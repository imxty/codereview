package product

import (
	"context"
	"errors"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/product/domain"
	"gorm.io/gorm"
)

type ProductStore struct {
	// 数据库管理
	*dbutils.Connection
}

func NewProductStore(management *dbutils.Connection) *ProductStore {
	return &ProductStore{
		management,
	}
}

// 实现 Domain 行为
var _ domain.ProductRepository = (*ProductStore)(nil)

// GetTreatmentRev 通过Treatment获取方案版本
func (p *ProductStore) GetTreatmentRev(ctx context.Context, organizationId, treatmentId string) (domain.TreatmentRevIntf, error) {
	// 根据treatment_id 查询方案最新版本信息
	var tr TreatmentRev
	err := p.GetConnection(ctx).Raw(`SELECT 
        treatment_rev.treatment_rev_id,
        treatment_rev.treatment_id,
        treatment_rev.treatment_name,
        treatment_rev.remarks,
        treatment_rev.organization_id,
        treatment_rev.is_published,
        treatment_rev.treatment_status,
        treatment_rev.rev,
        treatment_rev.created_at,
        treatment_rev.updated_at,
        treatment_rev.deleted_at
    FROM
        treatment_rev
            LEFT JOIN
        treatment ON treatment.latest_rev = treatment_rev.treatment_rev_id
            AND (treatment_rev.deleted_at IS NULL
            AND treatment.deleted_at IS NULL)
    WHERE
        treatment.treatment_id = ?
    AND treatment_rev.organization_id = ?
    `, treatmentId, organizationId).First(&tr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 转换
	dTreatmentRev := tr.ToDomainTreatmentRev()
	return dTreatmentRev, nil
}

// ListLatestTreatmentRevs 批量获取最新方案版本信息和方案总数
func (p *ProductStore) ListLatestTreatmentRevs(ctx context.Context, organizationID string) ([]domain.TreatmentRevIntf, error) {
	// 查询所有方案版本信息
	var treatmentRevs []TreatmentRev
	err := p.GetConnection(ctx).Raw(`SELECT 
    treatment_rev.treatment_rev_id,
    treatment_rev.treatment_id,
    treatment_rev.treatment_name,
    treatment_rev.remarks,
    treatment_rev.organization_id,
    treatment_rev.is_published,
    treatment_rev.treatment_status,
    treatment_rev.rev,
    treatment_rev.created_at,
    treatment_rev.updated_at,
    treatment_rev.deleted_at
    FROM
    treatment_rev
        LEFT JOIN
    treatment ON treatment.latest_rev = treatment_rev.treatment_rev_id
        AND (treatment_rev.deleted_at IS NULL
        AND treatment.deleted_at IS NULL)
WHERE
    treatment.organization_id = ?
    ORDER BY treatment_rev.created_at DESC, treatment_rev.treatment_name
    `, organizationID).Find(&treatmentRevs).Error
	if err != nil {
		return nil, err
	}

	// 转换
	dTreatmentRevs := make([]domain.TreatmentRevIntf, len(treatmentRevs))
	for i, v := range treatmentRevs {
		dTreatmentRevs[i] = v.ToDomainTreatmentRev()
	}
	return dTreatmentRevs, nil
}

// ListLatestTreatmentRevsWithStatus 批量获取方案最新版本信息和方案总数
func (p *ProductStore) ListLatestTreatmentRevsWithStatus(ctx context.Context, organizationID string, isPublish bool) ([]domain.TreatmentRevIntf, error) {
	// 查询所有方案版本信息
	var treatmentRevs []TreatmentRev
	err := p.GetConnection(ctx).Raw(`SELECT 
    treatment_rev.treatment_rev_id,
    treatment_rev.treatment_id,
    treatment_rev.treatment_name,
    treatment_rev.remarks,
    treatment_rev.organization_id,
    treatment_rev.is_published,
    treatment_rev.treatment_status,
    treatment_rev.rev,
    treatment_rev.created_at,
    treatment_rev.updated_at,
    treatment_rev.deleted_at
    FROM
    treatment_rev
        LEFT JOIN
    treatment ON treatment.latest_rev = treatment_rev.treatment_rev_id
        AND (treatment_rev.deleted_at IS NULL AND treatment_rev.is_published = ?
        AND treatment.deleted_at IS NULL)
WHERE
    treatment.organization_id = ?
    ORDER BY treatment_rev.created_at DESC, treatment_rev.treatment_name
    `, isPublish, organizationID).Find(&treatmentRevs).Error
	if err != nil {
		return nil, err
	}

	// 转换
	dTreatmentRevs := make([]domain.TreatmentRevIntf, len(treatmentRevs))
	for i, v := range treatmentRevs {
		dTreatmentRevs[i] = v.ToDomainTreatmentRev()
	}
	return dTreatmentRevs, nil
}

// GetTreatmentRevByTreatmentRevId 通过方案版本获取方案
func (p *ProductStore) GetTreatmentRevByTreatmentRevId(ctx context.Context, treatmentRevId string) (domain.TreatmentRevIntf, error) {
	var rev TreatmentRev
	err := p.GetConnection(ctx).Model(&TreatmentRev{}).Where("treatment_rev_id = ?", treatmentRevId).First(&rev).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return rev.ToDomainTreatmentRev(), nil
}

// ChangeTreatmentRevStatus 改变方案版本状态
func (p *ProductStore) ChangeTreatmentRevStatus(ctx context.Context, treatmentRevId string, rev, status int32) error {
	return p.GetConnection(ctx).Model(&TreatmentRev{}).Where("treatment_rev_id = ? and rev = ?", treatmentRevId, rev).Updates(map[string]interface{}{
		"rev":              rev + 1,
		"treatment_status": status,
	}).Error
}

// BatchGetRecommendedProductsByTreatmentRevId 通过方案id批量获取方案的所有商品信息
func (p *ProductStore) BatchGetRecommendedProductsByTreatmentRevId(ctx context.Context, treatmentRevId string) ([]domain.RecommendedProductIntf, error) {
	// 关联查询 方案表 方案配置表 商品表
	var products []RecommendedProduct
	err := p.GetConnection(ctx).Raw(`
    SELECT 
    recommended_product.recommended_product_id,
    recommended_product.organization_id,
    recommended_product.product_image_url,
    recommended_product.product_name,
    recommended_product.product_type,
    recommended_product.product_status,
    recommended_product.product_introduction,
    recommended_product.remarks,
    recommended_product.drug_name,
    recommended_product.drug_classification,
    recommended_product.approved_number,
    recommended_product.drug_validity_period,
    recommended_product.symptoms,
    recommended_product.rev
FROM
    recommended_product
      where recommended_product_id in
    (select recommended_product_id from treatment_item WHERE treatment_rev_id = ?)
        AND recommended_product.deleted_at IS NULL
    `, treatmentRevId).Find(&products).Error
	if err != nil {
		return nil, err
	}
	// 转换
	dProducts := make([]domain.RecommendedProductIntf, len(products))
	for i, v := range products {
		dProducts[i] = v.ToDomainRecommendedProduct()
	}
	return dProducts, nil
}

// ChangeProductStatus 修改商品状态
func (p *ProductStore) ChangeProductStatus(ctx context.Context, productId string, rev, status int32) error {
	return p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("recommended_product_id = ? AND rev = ?", productId, rev).Updates(map[string]interface{}{
		"rev":            rev + 1,
		"product_status": status,
	}).Error
}

// BatchChangeProductStatus 批量修改商品状态
func (p *ProductStore) BatchChangeProductStatus(ctx context.Context, productId []string, status int32) error {
	return p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("recommended_product_id IN (?)", productId).Updates(map[string]interface{}{
		"product_status": status,
	}).Error
}

// AddProduct 添加商品
func (p *ProductStore) AddProduct(ctx context.Context, organizationId string, product domain.RecommendedProductIntf) error {
	// 转换product
	var rp RecommendedProduct
	rp.FromDomainRecommendedProduct(product)

	// 创建商品记录
	return p.GetConnection(ctx).Model(&RecommendedProduct{}).Create(&rp).Error
}

// GetProduct 获取商品
func (p *ProductStore) GetProduct(ctx context.Context, productId string) (domain.RecommendedProductIntf, error) {
	var rp RecommendedProduct
	err := p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("recommended_product_id = ?", productId).First(&rp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return rp.ToDomainRecommendedProduct(), nil
}

// DeleteProduct 删除商品
func (p *ProductStore) DeleteProduct(ctx context.Context, productId string, rev int32) error {
	// 根据商品id删除商品
	err := p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("recommended_product_id = ?", productId).
		Updates(map[string]interface{}{
			"deleted_at":     time.Now().UTC(),
			"rev":            rev + 1,
			"product_status": domain.ProductStatusDeleted,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

// CreateTreatment 创建方案表
func (p *ProductStore) CreateTreatment(ctx context.Context, treatment domain.TreatmentIntf) error {
	var t Treatment
	t.FromDomainTreatment(treatment)
	return p.GetConnection(ctx).Model(&Treatment{}).Create(&t).Error
}

// CreateTreatmentRev 创建方案版本表
func (p *ProductStore) CreateTreatmentRev(ctx context.Context, treatmentRev domain.TreatmentRevIntf) error {
	var tr TreatmentRev
	tr.FromDomainTreatmentRev(treatmentRev)
	return p.GetConnection(ctx).Model(&TreatmentRev{}).Create(&tr).Error
}

// CreateSubTreatmentItems 创建多个方案配置表
func (p *ProductStore) CreateTreatmentItems(ctx context.Context, treatmentItems []domain.TreatmentItemIntf) error {
	// 转换
	dbTreatmentItems := make([]TreatmentItem, len(treatmentItems))
	for i, v := range treatmentItems {
		dbTreatmentItems[i].FromDomainTreatmentItem(v)
	}

	// 创建TreatmentItem表
	err := p.GetConnection(ctx).Model(&TreatmentItem{}).Create(&dbTreatmentItems).Error
	if err != nil {
		return err
	}
	return nil
}

// ListTreatmentItems 获取方案配置
func (p *ProductStore) ListTreatmentItems(ctx context.Context, treatmentId string) ([]domain.TreatmentProductAssociationIntf, error) {
	// 根据treatment_id关联查询treatment_items
	var associations []TreatmentProductAssociation
	err := p.GetConnection(ctx).Raw(`
    SELECT 
    treatment_item.treatment_item_id,
    treatment_item.treatment_rev_id,
    treatment_item.treatment_item_type,
    treatment_item.symptom,
    recommended_product.recommended_product_id,
    recommended_product.organization_id,
    recommended_product.product_image_url,
    recommended_product.product_name,
    recommended_product.product_type,
    recommended_product.product_status,
    recommended_product.product_introduction,
    recommended_product.remarks,
    recommended_product.drug_name,
    recommended_product.drug_classification,
    recommended_product.approved_number,
    recommended_product.drug_validity_period
FROM
    recommended_product
        LEFT JOIN
    treatment_item ON treatment_item.recommended_product_id = recommended_product.recommended_product_id
        AND (treatment_item.deleted_at IS NULL
        AND recommended_product.deleted_at IS NULL)
        LEFT JOIN
    treatment ON treatment.latest_rev = treatment_item.treatment_rev_id
        AND treatment.deleted_at IS NULL
WHERE
    treatment.treatment_id = ?
	`, treatmentId).Find(&associations).Error
	if err != nil {
		return nil, err
	}

	// 转换
	dAssociations := make([]domain.TreatmentProductAssociationIntf, len(associations))
	for i, v := range associations {
		dAssociations[i] = v.ToDomainTreatmentProductAssociation()
	}
	return dAssociations, nil
}

// ListTreatmentRevItems 获取方案rev配置
func (p *ProductStore) ListTreatmentRevItems(ctx context.Context, treatmentRevId string) ([]domain.TreatmentProductAssociationIntf, error) {
	// 根据treatment_id关联查询treatment_items
	var associations []TreatmentProductAssociation
	err := p.GetConnection(ctx).Raw(`
    SELECT 
        treatment_item.treatment_item_id,
        treatment_item.treatment_rev_id,
        treatment_item.treatment_item_type,
        treatment_item.symptom,
        recommended_product.recommended_product_id,
        recommended_product.organization_id,
        recommended_product.product_image_url,
        recommended_product.product_name,
        recommended_product.product_type,
        recommended_product.product_status,
        recommended_product.product_introduction,
        recommended_product.remarks,
        recommended_product.drug_name,
        recommended_product.drug_classification,
        recommended_product.approved_number,
        recommended_product.drug_validity_period
    FROM
        treatment_item
            INNER JOIN
        recommended_product ON treatment_item.recommended_product_id = recommended_product.recommended_product_id
    WHERE
        treatment_item.treatment_rev_id = ?  and recommended_product.deleted_at is null
	`, treatmentRevId).Find(&associations).Error
	if err != nil {
		return nil, err
	}

	// 转换
	dAssociations := make([]domain.TreatmentProductAssociationIntf, len(associations))
	for i, v := range associations {
		dAssociations[i] = v.ToDomainTreatmentProductAssociation()
	}
	return dAssociations, nil
}

// GetTreatment 获取方案信息
func (p *ProductStore) GetTreatment(ctx context.Context, treatmentId string) (domain.TreatmentIntf, error) {
	var tr Treatment
	err := p.GetConnection(ctx).Model(&Treatment{}).Where("treatment_id = ?", treatmentId).First(&tr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tr.ToDomainTreatment(), nil
}

// RemoveTreatmentByTreatmentRevId 下架方案
func (p *ProductStore) RemoveTreatmentByTreatmentRevId(ctx context.Context, treatmentRevId string, rev int32) error {
	return p.GetConnection(ctx).Model(&TreatmentRev{}).Where("is_published = 1 and rev = ? and treatment_rev_id = ?", rev, treatmentRevId).Updates(map[string]interface{}{
		"rev":          rev + 1,
		"is_published": 0,
	}).Error
}

// DeleteTreatment 删除方案表
func (p *ProductStore) DeleteTreatment(ctx context.Context, treatmentId string, rev int32) error {
	return p.GetConnection(ctx).Model(&Treatment{}).Where("treatment_id = ? AND rev = ?", treatmentId, rev).
		Updates(map[string]interface{}{
			"deleted_at": time.Now().UTC(),
			"rev":        rev + 1,
		}).Error
}

// GetOrganizationLatestTreatment 获取组织最新的方案
func (p *ProductStore) GetOrganizationLatestTreatmentRev(ctx context.Context, organizationID string) (domain.TreatmentRevIntf, error) {
	// 查询方案表 获取租户所有最新方案 从中获取发布状态为1的方案版本
	var tenantTreatmentRev TreatmentRev
	err := p.GetConnection(ctx).Model(&TreatmentRev{}).Where("is_published = 1 and organization_id = ?", organizationID).First(&tenantTreatmentRev).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tenantTreatmentRev.ToDomainTreatmentRev(), nil
}

// GetPublishedTreatmentRev 获取现有的已发布方案
func (p *ProductStore) GetPublishedTreatmentRev(ctx context.Context, organizationID string) (domain.TreatmentRevIntf, error) {
	// 查询方案表 获取租户所有最新方案 从中获取发布状态为1的方案版本
	var tenantTreatmentRev TreatmentRev
	err := p.GetConnection(ctx).Model(&TreatmentRev{}).Where("is_published = 1 and organization_id = ?", organizationID).First(&tenantTreatmentRev).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tenantTreatmentRev.ToDomainTreatmentRev(), nil
}

// ChangeTreatmentFromPublishedToUnpublished 修改商品方案的最新方案版本的发布状态为未发布
func (p *ProductStore) ChangeTreatmentFromPublishedToUnpublished(ctx context.Context, treatmentRevId string, rev int32) error {
	// 关联查询treatment
	return p.GetConnection(ctx).Model(&TreatmentRev{}).Where("treatment_rev_id = ? AND rev = ?", treatmentRevId, rev).
		Updates(map[string]interface{}{
			"rev":          rev + 1,
			"is_published": domain.TreatmentRevIsNotPublished,
		}).Error
}

// PublishTreatment 发布方案
func (p *ProductStore) PublishTreatment(ctx context.Context, treatmentId string, rev int32) error {
	// 子查询在treatment表中查treatment_rev_id 然后更新
	return p.GetConnection(ctx).Exec(`
        UPDATE treatment_rev
    SET
        treatment_rev.is_published = ?,
        treatment_rev.rev = ?,
        treatment_rev.updated_at = utc_timestamp()
    WHERE
        treatment_rev.rev = ?
    AND treatment_rev_id = (SELECT
                treatment.latest_rev
            FROM
                treatment
            WHERE
                treatment_id = ?)
            `, domain.TreatmentRevIsPublished, rev+1, rev, treatmentId).Error
}

// GetProductStatisticsWithTimeRange 根据获取商品统计
func (p *ProductStore) GetProductStatisticsWithTimeRange(ctx context.Context, productId string, startTime, endTime time.Time) ([]domain.RecommendedProductStatisticsIntf, error) {
	var ps []RecommendedProductStatistics
	err := p.GetConnection(ctx).Model(&RecommendedProductStatistics{}).Where("recommended_product_id = ? and exposed_at between ? and ?", productId, startTime, endTime).
		Find(&ps).Error
	if err != nil {
		return nil, err
	}
	dp := make([]domain.RecommendedProductStatisticsIntf, len(ps))
	for k, v := range ps {
		dp[k] = v.ToDomainRecommendedProductStatistics()
	}
	return dp, nil
}

// CountMonthlyTreatmentSubmitTimes 查询当月商品方案提审次数
func (p *ProductStore) CountMonthlyTreatmentSubmitTimes(ctx context.Context, organizationID string, startTime, endTime time.Time) (int32, error) {
	var times int64
	err := p.GetConnection(ctx).Table("review_issue").Where("submitter_organization_id = ? and target_type = 1 and created_at between ? and ?", organizationID, startTime.UTC(), endTime.UTC()).
		Count(&times).Error
	if err != nil {
		return 0, err
	}
	return int32(times), nil
}

// ListProducts 根据分页获取租户所有商品
func (p *ProductStore) ListProducts(ctx context.Context, organizationID string, status int32, offset, size int) ([]domain.RecommendedProductIntf, error) {
	// 按分页信息查询一个租户下的所有商品
	var products []RecommendedProduct
	err := p.GetConnection(ctx).Raw(`SELECT 
    r.recommended_product_id,
    r.organization_id,
    r.product_image_url,
    r.product_name,
    r.product_type,
    r.product_status,
    r.product_introduction,
    r.remarks,
    r.drug_name,
    r.drug_classification,
    r.approved_number,
    r.drug_validity_period,
    r.symptoms,
    r.rev,
    r.created_at,
    r.updated_at,
    r.deleted_at
    FROM 
        recommended_product r
    WHERE 
        r.organization_id = ? and product_status = ?
    ORDER BY r.product_name, r.created_at
    LIMIT ?
    OFFSET ?
    `, organizationID, status, size, offset).Find(&products).Error
	if err != nil {
		return nil, err
	}
	// 转换
	dProducts := make([]domain.RecommendedProductIntf, len(products))
	for i, v := range products {
		dProducts[i] = v.ToDomainRecommendedProduct()
	}
	return dProducts, nil
}

// ListAllProducts 获取租户下所有商品
func (p *ProductStore) ListAllProducts(ctx context.Context, organizationID string, offset, size int) ([]domain.RecommendedProductIntf, error) {

	var products []RecommendedProduct
	err := p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("organization_id = ?", organizationID).Offset(offset).Limit(size).Find(&products).Error
	if err != nil {
		return nil, err
	}
	// 转换
	dProducts := make([]domain.RecommendedProductIntf, len(products))
	for i, v := range products {
		dProducts[i] = v.ToDomainRecommendedProduct()
	}
	return dProducts, nil
}

// ListOrganizationProducts 获取组织下所有商品
func (p *ProductStore) ListOrganizationProducts(ctx context.Context, organizationID string) ([]domain.RecommendedProductIntf, error) {
	var products []RecommendedProduct
	err := p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("organization_id = ?", organizationID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	// 转换
	dProducts := make([]domain.RecommendedProductIntf, len(products))
	for i, v := range products {
		dProducts[i] = v.ToDomainRecommendedProduct()
	}
	return dProducts, nil
}

// GetProductsTotalCount 获取租户下商品总数量
func (p *ProductStore) GetProductsTotalCount(ctx context.Context, organizationId string) (int64, error) {
	// 查询所有商品数量
	var count int64
	err := p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("organization_id = ?", organizationId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductsTotalCountWithCertainStatus 获取租户下指定状态商品总数量
func (p *ProductStore) GetProductsTotalCountWithCertainStatus(ctx context.Context, organizationID string, status int32) (int64, error) {
	// 查询指定状态商品数量
	var count int64
	err := p.GetConnection(ctx).Raw(`SELECT COUNT(*) FROM recommended_product WHERE organization_id = ? and product_status = ?`, organizationID, status).Find(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateRecommendedProductStandingBook 创建商品统计台账
func (p *ProductStore) CreateRecommendedProductStandingBook(ctx context.Context, rps []domain.RecommendedProductStatisticsIntf) error {
	// 转换
	dbRps := make([]RecommendedProductStatistics, len(rps))
	for i, v := range rps {
		dbRps[i].FromDomainRecommendedProductStatistics(v)
	}
	// 创建商品统计表
	return p.GetConnection(ctx).Model(&RecommendedProductStatistics{}).Create(&dbRps).Error
}

// BatchGetTreatmentsWithCertainStatus 获取指定状态的方案
func (p *ProductStore) BatchGetTreatmentsWithCertainStatus(ctx context.Context, organizationID string, status int32, offset, size int) ([]domain.TreatmentRevIntf, error) {
	// 查询所有方案版本信息
	var treatmentRevs []TreatmentRev
	err := p.GetConnection(ctx).Raw(`SELECT 
    treatment_rev.treatment_rev_id,
    treatment_rev.treatment_id,
    treatment_rev.treatment_name,
    treatment_rev.remarks,
    treatment_rev.organization_id,
    treatment_rev.is_published,
    treatment_rev.treatment_status,
    treatment_rev.rev,
    treatment_rev.created_at,
    treatment_rev.updated_at,
    treatment_rev.deleted_at
    FROM
    treatment_rev
        LEFT JOIN
    treatment ON treatment.latest_rev = treatment_rev.treatment_rev_id
        AND (treatment_rev.deleted_at IS NULL
        AND treatment.deleted_at IS NULL)
    WHERE
        treatment.organization_id = ?
    AND treatment_rev.treatment_status = ?
        ORDER BY treatment_rev.created_at DESC, treatment_rev.treatment_name
        LIMIT ? OFFSET ?
    `, organizationID, status, size, offset).Find(&treatmentRevs).Error
	if err != nil {
		return nil, err
	}

	// 转换
	dTreatmentRevs := make([]domain.TreatmentRevIntf, len(treatmentRevs))
	for i, v := range treatmentRevs {
		dTreatmentRevs[i] = v.ToDomainTreatmentRev()
	}
	return dTreatmentRevs, nil
}

// GetTreatmentsCountWithCertainStatus 获取指定状态的方案总数
func (p *ProductStore) GetTreatmentsCountWithCertainStatus(ctx context.Context, organizationID string, status int32) (int32, error) {
	// 查询租户所有方案总数
	var count int64
	// err = db.Model(&Treatment{}).Where("tenant_id = ? and treatment_status = ?", tenantId, status).Count(&count).Error
	err := p.GetConnection(ctx).Raw(`SELECT 
    COUNT(*) FROM treatment
    LEFT JOIN treatment_rev
    ON treatment.latest_rev = treatment_rev.treatment_rev_id
    WHERE  treatment.organization_id = ?
    AND treatment_rev.treatment_status = ?
    AND treatment.deleted_at IS NULL
    `, organizationID, status).Find(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// GetTreatmentsCount 获取方案总数
func (p *ProductStore) GetTreatmentsCount(ctx context.Context, organizationID string) (int32, error) {
	var count int64
	err := p.GetConnection(ctx).Model(&Treatment{}).Where("organization_id = ?", organizationID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

// UpdateProduct 更新商品
func (p *ProductStore) UpdateProduct(ctx context.Context, productId string, updates map[string]interface{}, rev int32) error {
	// 更新字段 版本号加1
	updates["rev"] = rev + 1
	return p.GetConnection(ctx).Model(&RecommendedProduct{}).Where("recommended_product_id = ? AND rev = ?", productId, rev).
		Updates(updates).Error
}

// UpdateLatestTreatmentRev 更新方案表里的最新方案版本id
func (p *ProductStore) UpdateLatestTreatmentRev(ctx context.Context, treatmentId, latestRev string, rev int32) error {
	// 更新treatment表的latest_rev字段
	return p.GetConnection(ctx).Model(&Treatment{}).Where("treatment_id = ? AND rev = ?", treatmentId, rev).
		Updates(map[string]interface{}{
			"latest_rev": latestRev,
			"rev":        rev + 1,
		}).Error
}

// CheckProductUsage 检测商品使用状态.
func (p *ProductStore) CheckProductUsage(ctx context.Context, oid string) ([]domain.TreatmentRevItemIntf, error) {
	var items []TreatmentRevItem
	err := p.GetConnection(ctx).Raw(`
        SELECT 
        a.treatment_name,
        a.treatment_status,
        b.*  
        FROM
        (
            SELECT * FROM
            treatment_rev
        WHERE treatment_rev_id IN 
            (SELECT latest_rev FROM treatment WHERE organization_id = ? AND deleted_at IS NULL)
        )
        a left join treatment_item b on a.treatment_rev_id = b.treatment_rev_id;`, oid).Find(&items).Error
	if err != nil {
		return nil, err
	}

	dItems := make([]domain.TreatmentRevItemIntf, len(items))
	for k, v := range items {
		dItems[k] = v.ToDomainTreatmentRevItem()
	}
	return dItems, nil
}

// ListTreatmentProductsNotUsingWithoutItself 获取还在使用中的商品ID除去本身，返回还在使用中的id
func (p *ProductStore) ListTreatmentProductsNotUsingWithoutItself(ctx context.Context, productsIds []string, treatmentId string) ([]string, error) {
	var productIds []string
	err := p.GetConnection(ctx).Raw(`
        SELECT 
            recommended_product_id 
        FROM treatment_item
        WHERE treatment_rev_id in
            (   SELECT 
                    treatment_rev_id 
                FROM treatment_rev
                WHERE is_published = 1 and treatment_rev_id != ?
            )
        GROUP BY
            recommended_product_id 
        HAVING recommended_product_id 
        IN (?)`, treatmentId, productsIds).Find(&productIds).Error
	if err != nil {
		return nil, err
	}
	return productIds, nil
}

// ListTreatmentRevs
func (p *ProductStore) ListTreatmentRevs(ctx context.Context, treatmentIds []string) ([]domain.TreatmentRevIntf, error) {
	var treatmentRevs []TreatmentRev
	err := p.GetConnection(ctx).Raw(`
    SELECT b.* FROM
        (
            SELECT
                *
            FROM treatment 
            WHERE treatment_id IN (?)
        ) a
    LEFT JOIN
        treatment_rev b on a.latest_rev = b.treatment_rev_id
    `, treatmentIds).Find(&treatmentRevs).Error
	if err != nil {
		return nil, err
	}

	// 转化
	tRevs := make([]domain.TreatmentRevIntf, len(treatmentRevs))
	for k, v := range treatmentRevs {
		tRevs[k] = v.ToDomainTreatmentRev()
	}
	return tRevs, nil
}
