package domain

import (
	"context"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type ProductRepository interface {
	dbutils.Tx

	// GetTreatmentRev 通过Treatment获取方案版本
	GetTreatmentRev(ctx context.Context, organizationId, treatmentId string) (TreatmentRevIntf, error)
	// ListLatestTreatmentRevs 批量获取方案最新版本信息和方案总数
	ListLatestTreatmentRevs(ctx context.Context, organizationID string) ([]TreatmentRevIntf, error)
	// ListLatestTreatmentRevsWithStatus 批量获取方案最新版本信息和方案总数
	ListLatestTreatmentRevsWithStatus(ctx context.Context, organizationID string, isPublish bool) ([]TreatmentRevIntf, error)
	// GetTreatmentRevByTreatmentRevId 通过方案版本获取方案
	GetTreatmentRevByTreatmentRevId(ctx context.Context, treatmentRevId string) (TreatmentRevIntf, error)
	// ChangeTreatmentRevStatus 改变方案版本状态
	ChangeTreatmentRevStatus(ctx context.Context, treatmentRevId string, rev, status int32) error
	// BatchGetRecommendedProductsByTreatmentRevId 通过方案版本id批量获取方案的所有商品信息
	BatchGetRecommendedProductsByTreatmentRevId(ctx context.Context, treatmentRevId string) ([]RecommendedProductIntf, error)
	// ChangeProductStatus 修改商品状态
	ChangeProductStatus(ctx context.Context, productId string, rev, status int32) error
	// BatchChangeProductStatus 批量修改商品状态
	BatchChangeProductStatus(ctx context.Context, productId []string, status int32) error
	// AddProduct 添加商品
	AddProduct(ctx context.Context, organizationId string, product RecommendedProductIntf) error
	// GetProduct 获取商品
	GetProduct(ctx context.Context, productId string) (RecommendedProductIntf, error)
	// DeleteProduct 删除商品
	DeleteProduct(ctx context.Context, productId string, rev int32) error

	// CreateTreatment 创建方案表
	CreateTreatment(ctx context.Context, treatment TreatmentIntf) error
	// CreateTreatmentRev 创建方案版本表
	CreateTreatmentRev(ctx context.Context, treatmentRev TreatmentRevIntf) error
	// CreateSubTreatmentItems 创建多个方案配置表
	CreateTreatmentItems(ctx context.Context, treatmentItems []TreatmentItemIntf) error
	// ListTreatmentItems 获取方案配置
	ListTreatmentItems(ctx context.Context, treatmentId string) ([]TreatmentProductAssociationIntf, error)
	// ListTreatmentRevItems 获取方案rev配置
	ListTreatmentRevItems(ctx context.Context, treatmentRevId string) ([]TreatmentProductAssociationIntf, error)
	// GetTreatment 获取方案信息
	GetTreatment(ctx context.Context, treatmentId string) (TreatmentIntf, error)
	// RemoveTreatmentByTreatmentRevId 下架方案
	RemoveTreatmentByTreatmentRevId(ctx context.Context, treatmentRevId string, rev int32) error
	// DeleteTreatment 删除方案表
	DeleteTreatment(ctx context.Context, treatmentId string, rev int32) error
	// GetOrganizationLatestTreatment 获取组织最新的方案
	GetOrganizationLatestTreatmentRev(ctx context.Context, organizationID string) (TreatmentRevIntf, error)
	// GetPublishedTreatmentRev 获取租户现有的已发布方案
	GetPublishedTreatmentRev(ctx context.Context, organizationID string) (TreatmentRevIntf, error)
	// ChangeTreatmentFromPublishedToUnpublished 修改商品方案的最新方案版本的发布状态为未发布
	ChangeTreatmentFromPublishedToUnpublished(ctx context.Context, treatmentRevId string, rev int32) error
	// PublishTreatment 发布方案
	PublishTreatment(ctx context.Context, treatmentId string, rev int32) error

	// GetProductStatisticsWithTimeRange 根据获取商品统计
	GetProductStatisticsWithTimeRange(ctx context.Context, productId string, startTime, endTime time.Time) ([]RecommendedProductStatisticsIntf, error)
	// CountMonthlyTreatmentSubmitTimes 查询当月商品方案提审次数
	CountMonthlyTreatmentSubmitTimes(ctx context.Context, organizationID string, startTime, endTime time.Time) (int32, error)

	// ListProducts 根据分页获取组织所有商品
	ListProducts(ctx context.Context, organizationID string, status int32, offset, size int) ([]RecommendedProductIntf, error)
	// ListAllProducts 获取组织下所有商品
	ListAllProducts(ctx context.Context, organizationID string, offset, size int) ([]RecommendedProductIntf, error)
	// ListOrganizationProducts 获取组织下所有商品
	ListOrganizationProducts(ctx context.Context, organizationID string) ([]RecommendedProductIntf, error)
	// GetProductsTotalCount 获取租户下商品总数量
	GetProductsTotalCount(ctx context.Context, organizationId string) (int64, error)
	// GetProductsTotalCountWithCertainStatus 获取租户下指定状态商品总数量
	GetProductsTotalCountWithCertainStatus(ctx context.Context, organizationID string, status int32) (int64, error)
	// CreateRecommendedProductStandingBook 创建商品统计台账
	CreateRecommendedProductStandingBook(ctx context.Context, rps []RecommendedProductStatisticsIntf) error
	// BatchGetTreatmentsWithCertainStatus 获取指定状态的方案
	BatchGetTreatmentsWithCertainStatus(ctx context.Context, organizationID string, status int32, offset, size int) ([]TreatmentRevIntf, error)
	// UpdateProduct 更新商品
	UpdateProduct(ctx context.Context, productId string, updates map[string]interface{}, rev int32) error
	// UpdateLatestTreatmentRev 更新方案表里的最新方案版本id
	UpdateLatestTreatmentRev(ctx context.Context, treatmentId, latestRev string, rev int32) error
	// CheckProductUsage 检测商品使用状态.
	CheckProductUsage(ctx context.Context, oid string) ([]TreatmentRevItemIntf, error)

	// ListTreatmentProductsNotUsingWithoutItself 获取还在使用中的商品ID除去本身，返回还在使用中的id
	ListTreatmentProductsNotUsingWithoutItself(ctx context.Context, productsIds []string, treatmentId string) ([]string, error)
	// ListTreatmentRevs
	ListTreatmentRevs(ctx context.Context, treatmentIds []string) ([]TreatmentRevIntf, error)
}
