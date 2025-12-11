package domain

import (
	"context"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type UserRepository interface {
	dbutils.Tx

	// CreateUser 创建用户
	CreateUser(ctx context.Context, user UserIntf) error
	// UpdateTenantUser 更新商户 user
	UpdateTenantUser(ctx context.Context, tid string, updateMap map[string]interface{}) error
	// CheckOrganizationIDExist 检测组织 ID 是否存在
	CheckOrganizationIDExist(ctx context.Context, id string) (bool, error)
	// GetOrganizationByPhone 通过手机号查询组织信息
	GetOrganizationByPhone(ctx context.Context, phone string) (UserIntf, error)
	// GetOrganizationByOrganizationId 通过组织 ID 查询组织信息
	GetOrganizationByOrganizationId(ctx context.Context, organizationID string) (UserIntf, error)
	// GetOrganizationByUsername 通过用户名查询组织信息
	GetOrganizationByUsername(ctx context.Context, username string) (UserIntf, error)
	// GetOrganizationByVagueUsername 通过模糊用户名查询组织信息
	GetOrganizationByVagueUsername(ctx context.Context, username string) ([]UserIntf, error)

	// GetOrganizationTenant 获取组织商户信息
	GetOrganizationTenant(ctx context.Context, organizationID, tenantID string) (OrganizationTenantIntf, error)
	// GetOrganizationTenantByTenantId 通过组织ID获取组织商户信息
	GetOrganizationTenantByTenantId(ctx context.Context, tenantID string) (OrganizationTenantIntf, error)
	// CreateOrganizationTenant 创建组织商户信息
	CreateOrganizationTenant(ctx context.Context, ot OrganizationTenantIntf) error
	// CreateTenantSubscriptionTimeline 创建商户订阅时间线
	CreateTenantSubscriptionTimeline(ctx context.Context, timeline SubscriptionTimelineIntf) error
	// BatchCreateTenantSubscriptionTimeline 批量创建商户订阅时间线
	BatchCreateTenantSubscriptionTimeline(ctx context.Context, timelines []SubscriptionTimelineIntf) error
	// ChangeOrganizationTenantReviewStatus 改变组织商户审核状态
	ChangeOrganizationTenantReviewStatus(ctx context.Context, tenantId string, status, rev int32) error
	// ListOrganizationTenants 获取组织商户列表信息
	ListOrganizationTenants(ctx context.Context, organizationID string, tenantID []string) ([]OrganizationTenantIntf, error)
	// ListTenants 获取商户列表信息
	ListTenants(ctx context.Context, tenantID []string) ([]OrganizationTenantIntf, error)
	// ListTenantsContainsDeleted 获取商户列表信息
	ListTenantsContainsDeleted(ctx context.Context, tenantID []string) ([]OrganizationTenantIntf, error)
	// ListTenantsStatus 获取商户列表状态
	ListTenantsStatus(ctx context.Context, tenantID []string) ([]OrganizationTenantIntf, error)
	// ListOrganizationTenantsWithoutTenantId 获取组织商户列表信息
	ListOrganizationTenantsWithoutTenantId(ctx context.Context, organizationID string, offset, size int32) ([]OrganizationTenantIntf, int64, error)
	// ListOrganizationExistTenants 获取所有未删除的商户ID
	ListOrganizationExistTenants(ctx context.Context, organizationID string) ([]OrganizationTenantIntf, error)
	// DeleteOrganizationTenant 删除商户组织关系
	DeleteOrganizationTenant(ctx context.Context, organizationID, tenantID string) error
	// ListOrganizationAuthTenants 获取所有认证组织商户
	ListOrganizationAuthTenants(ctx context.Context, organizationID string) ([]OrganizationTenantIntf, error)
	// ListAllOrganizationTenantIds 获取组织下所有商户ID
	ListAllOrganizationTenantIds(ctx context.Context, oid string) ([]string, error)

	// GetTenantEntity 获取租户实体信息
	GetTenantEntity(ctx context.Context, tid string) (TenantEntityIntf, error)
	// ListTenantEntity 获取商户信息列表
	ListTenantEntity(ctx context.Context, tids []string) ([]TenantEntityIntf, error)
	// GetTenantUser 获取租户User信息
	GetTenantUser(ctx context.Context, tid string) (UserIntf, error)
	// GetTenantUserByPhone 通过手机号获取租户User信息
	GetTenantUserByPhone(ctx context.Context, phone string) (UserIntf, error)
	// UpdateTenantUserPassword 更新租户密码
	UpdateTenantUserPassword(ctx context.Context, tenantId string, newHashedPassword string, rev int32) error
	// UpdateOrganizationUserPassword 更新组织密码
	UpdateOrganizationUserPassword(ctx context.Context, organizationId string, newHashedPassword string, rev int32) error
	// GetActivateStaffByPhone 通过手机号获取已经激活员工
	GetActivateStaffByPhone(ctx context.Context, phone string) (UserIntf, error)
	// ListStaffs 获取员工列表
	ListStaffs(ctx context.Context, tid string) ([]UserIntf, error)
	// ListActivatedStaffs 获取激活员工列表
	ListActivatedStaffs(ctx context.Context, tid string) ([]UserIntf, error)
	// GetTenantStaffByPhone 通过手机号获取商户员工
	GetTenantStaffByPhone(ctx context.Context, tenantId, phone string) (UserIntf, error)
	// CreateStaff 创建员工
	CreateStaff(ctx context.Context, staff UserIntf) error

	// CreateTenantEntityRevision 创建商户副本
	CreateTenantEntityRevision(ctx context.Context, revision TenantEntityRevisionIntf) error
	// CreateTenantEntity 创建商户
	CreateTenantEntity(ctx context.Context, te TenantEntityIntf) error
	// GetStaffByID 通过员工ID获取员工
	GetStaffByID(ctx context.Context, staffId string) (UserIntf, error)
	// BatchGetStaffs 批量获取员工
	BatchGetStaffs(ctx context.Context, staffIds []string) ([]UserIntf, error)
	// UpdateStaffNickname 更新员工的昵称
	UpdateStaffNickname(ctx context.Context, staffId, newName string, rev int32) error
	// DeleteStaff 删除员工
	DeleteStaff(ctx context.Context, staffId string) error
	// GetLatestTenantSubscriptionTimeline 获取最新的商户订阅时间线
	GetLatestTenantSubscriptionTimeline(ctx context.Context, organizationId, tenantId string) (SubscriptionTimelineIntf, error)
	// BatchGetLatestTenantSubscriptionTimelinesDetails 批量获取商户订阅时间线
	BatchGetLatestTenantSubscriptionTimelinesDetails(ctx context.Context, tenantIds []string) ([]SubscriptionTimelineIntf, int64, error)
	// GetLatestTenantEntityRevision 获取最新的revision
	GetLatestTenantEntityRevision(ctx context.Context, tenantId string) (TenantEntityRevisionIntf, error)
	// GetTenantEntityRevision 获取最新的revision
	GetTenantEntityRevision(ctx context.Context, revisionId string) (TenantEntityRevisionIntf, error)

	// GetStaffByPhone 通过手机号查询员工
	GetStaffByPhone(ctx context.Context, areaCode, phone string) (UserIntf, error)
	// UpdateStaff 更新员工信息
	UpdateStaff(ctx context.Context, staffId, phone, name, hashedPassword string, rev int32) error
	// LeaveTenant 员工离开租户
	LeaveTenant(ctx context.Context, staffId, tenantId string, rev int32) error
	// UpdateStaffPassword 更新员工密码
	UpdateStaffPassword(ctx context.Context, staffId, hashedPassword string, rev int32) error
	// ListAllStaffs 获取所有的员工包括删除的
	ListAllStaffs(ctx context.Context, tenantId string) ([]UserIntf, error)
	// UpdateTenantPassword 更新商户密码
	UpdateTenantPassword(ctx context.Context, tid, hashedPassword string, rev int32) error
	// UpdateOrganizationPassword 更新组织密码
	UpdateOrganizationPassword(ctx context.Context, oid, hashedPassword string, rev int32) error
	// UpdateTenantEntity 更新租户实体
	UpdateTenantEntity(ctx context.Context, revision TenantEntityRevisionIntf, rev int32) error
	// UpdateOrganizationTenantStatus
	UpdateOrganizationTenantStatus(ctx context.Context, tenantId string, reviewStatus, rev int32) error

	// ModifyReportSharingStatus 修改报告共享状态(status为目标状态)
	ModifyReportSharingStatus(ctx context.Context, tid string, status bool, rev int32) error

	// GetSystemUserByPhone 通过手机号获取系统用户
	GetSystemUserByPhone(ctx context.Context, phone string) (UserIntf, error)
	// GetSystemUserByID 通过用户ID获取系统用户
	GetSystemUserByID(ctx context.Context, systemUserID string) (UserIntf, error)
	// ListSystemUsers 获取系统用户列表
	ListSystemUsers(ctx context.Context, userIds []string) ([]UserIntf, error)
	// ListSystemOperators 获取系统操作员
	ListSystemOperators(ctx context.Context) ([]UserIntf, error)

	// GetSubscriptionTimeline 获取订阅时间线
	GetSubscriptionTimeline(ctx context.Context, tenantId string) (SubscriptionTimelineIntf, error)
	// ListSubscriptionTimeline 获取订阅时间线
	ListSubscriptionTimeline(ctx context.Context, tenantIds []string) ([]SubscriptionTimelineIntf, error)
	// CreateSubscriptionTimeline 创建订阅
	CreateSubscriptionTimeline(ctx context.Context, timelines []SubscriptionTimelineIntf) error
	// RenewSubscriptionTimeline 续费订阅(未过期)
	RenewSubscriptionTimeline(ctx context.Context, tenantIds []string, renewDays int32) error
	// ListOrganizations 获取组织
	ListOrganizations(ctx context.Context, oids []string) ([]UserIntf, error)
	// ListAllOrganizations 获取所有组织信息
	ListAllOrganizations(ctx context.Context) ([]UserIntf, error)

	// ListTenantsActivationTimeline 获取商户订阅时间线.
	ListTenantsActivationTimeline(ctx context.Context) ([]SubscriptionTimelineIntf, error)

	// CreatePrivilegeGroup 创建权限组
	CreatePrivilegeGroup(ctx context.Context, privilegeId, privilegeName, remark string) error
	// CreatePagePrivileges 创建页面权限
	CreatePagePrivileges(ctx context.Context, privilege_id string, preset_ids []string) error
	// CreateDataPrivileges 创建数据权限
	CreateDataPrivileges(ctx context.Context, privilegeId string, organizationIds []string) error
	// GetPrivilegeGroupById 通过 ID 获取权限组
	GetPrivilegeGroupById(ctx context.Context, privilege_id string) (PrivilegePolicyIntf, error)
	// GetPrivilegeGroupsByName 通过权限组名获取权限组
	GetPrivilegeGroupsByName(ctx context.Context, privilege_name string) ([]PrivilegePolicyIntf, error)

	// CreateSystemUser 创建系统用户
	CreateSystemUser(ctx context.Context, systemUser UserIntf) error
	// DeleteSystemUser 删除系统用户
	DeleteSystemUser(ctx context.Context, user_id string) error
	// DeleteTenantUser 删除商户
	DeleteTenantUser(ctx context.Context, tenant_id string) error

	// BatchGetTenantNamesByIds 通过 ID 批量获取商户名
	BatchGetTenantNamesByIds(ctx context.Context, tenant_ids []string) ([]TenantEntityIntf, error)
	// GetUsingTenantIds 获取使用中商户 ID
	GetUsingTenantIds(ctx context.Context, organization_id string) ([]string, error)
	// GetTodayAddedTenantIds 获取今日新增商户 ID
	GetTodayAddedTenantIds(ctx context.Context, organization_id string) ([]string, error)
	// GetExpiringTenantIds 获取本月即将到期商户数量
	GetExpiringTenantIds(ctx context.Context, organization_id string) ([]string, error)

	// ListPrivilegeGroups 获取权限组列表
	ListPrivilegeGroups(ctx context.Context) ([]PrivilegePolicyIntf, error)

	// SearchOrganizationByName 根据组织名模糊查询组织名
	SearchOrganizationByName(ctx context.Context, organization_name string) ([]UserIntf, error)

	// SearchSystemUser 查询系统用户
	SearchSystemUser(ctx context.Context, staff_name, contact_phone string, privilege_group_ids []string, is_activated *wrapperspb.BoolValue, offset, size int32) ([]UserIntf, int32, error)

	// SearchTenantIdsByName 通过名称查询商户 ID
	SearchTenantIdsByName(ctx context.Context, tenant_name string) ([]string, error)

	// SubmitTreatmentToTenant 分配方案
	SubmitTreatmentToTenant(ctx context.Context, organization_id string, tenant_ids []string, treatment_id string) error

	UpdateTreatmentToTenant(ctx context.Context, organization_id string, tenant_ids []string, treatment_id string) error

	// UpdateOrganizationPhone 修改组织手机号
	UpdateOrganizationPhone(ctx context.Context, organization_id string, new_phone string) error

	// UpdateTenantConstitutionSwitch 修改商户体质辩证开关
	UpdateTenantConstitutionSwitch(ctx context.Context, tenant_id string, constitution_switch_status bool) error

	// UpdateTenantReview 修改复查天数
	UpdateTenantReview(ctx context.Context, tenant_id string, days int32) error

	// BatchGetOrganizationNamesByIDs 通过组织 ID 批量获取组织名
	BatchGetOrganizationNamesByIDs(ctx context.Context, organization_ids []string) ([]UserIntf, error)

	// GetDataPrivilegeByOrganizationID 通过组织 ID 获取数据权限
	GetDataPrivilegeByOrganizationID(ctx context.Context, privilege_id, organization_id string) (PrivilegePolicyDataIntf, error)

	// GetDataPrivilegeByID 通过权限组 ID 获取数据权限
	GetDataPrivilegeByID(ctx context.Context, privilege_id string) ([]string, error)

	// GetPagePrivilegeByID 通过权限组 ID 获取页面权限
	GetPagePrivilegeByID(ctx context.Context, privilege_id string) ([]PrivilegePolicyPageIntf, error)

	// SearchStaffsByNameAndStatus 根据员工名和状态模糊查询
	SearchStaffsByNameAndStatus(ctx context.Context, tenant_id, staff_name string, is_activated *wrapperspb.BoolValue) ([]UserIntf, error)

	// GetTenantByName 根据商户名获取商户
	GetTenantByName(ctx context.Context, name string) (UserIntf, error)

	// ResetPassword 重置平台管理员密码
	ResetPassword(ctx context.Context, user_id, new_plain_password string) error

	// CancelRevisionIssue 取消 issue
	CancelRevisionIssue(ctx context.Context, revisionId string) error

	// ListTreatmentsByTenantIDs 根据商户 ID 批量获取已分配方案
	ListTreatmentsByTenantIDs(ctx context.Context, tids []string) ([]TenantTreatmentIntf, error)

	// SearchTenantByNameAndOrganizationID 通过组织 ID 和商户名获取商户 ID
	SearchTenantByNameAndOrganizationID(ctx context.Context, organization_id, tenant_name string) ([]string, error)

	// SearchTenantsByName 搜索商户
	SearchTenantsByName(ctx context.Context, tids []string, tenant_name, contact_name, contact_phone string, status, offset, size int32) ([]TenantEntityIntf, int32, error)

	// BatchGetTenantsSubscriptionPeriods 批量获取商户订阅信息
	BatchGetTenantsSubscriptionPeriods(ctx context.Context, tenant_ids []string, start_time, end_time time.Time, offset, size int32) ([]SubscriptionPeriodIntf, int64, error)

	// UpdateTenantRevisionFailReason 更新商户审核失败原因
	UpdateTenantRevisionFailReason(ctx context.Context, revision_id, fail_reason string) error

	// BatchCreateTenantSubscriptionPeriod 批量创建订阅周期
	BatchCreateTenantSubscriptionPeriod(ctx context.Context, periods []SubscriptionPeriodIntf) error

	// BatchAddTenantTimelines 批量添加订阅年限
	BatchAddTenantTimelines(ctx context.Context, tids []string, end_time time.Time, years int32) error

	// BatchGetOrganizationTenants 批量获取组织商户关系
	BatchGetOrganizationTenants(ctx context.Context, tids []string) ([]OrganizationTenantIntf, error)

	// GetMenus 获取菜单
	GetMenus(ctx context.Context, privilege_id string) ([]MenuIntf, error)

	// BatchUpdateTenantTimelines 批量更新订阅年限
	BatchUpdateTenantTimelines(ctx context.Context, tids []string, add_time int32, years int32) error

	// UpdateStaffInfo 更新员工信息
	UpdateStaffInfo(ctx context.Context, staff_id, staff_name, staff_phone, privilege_id, remark string) error

	// SubmitFeedback 提交反馈
	SubmitFeedback(ctx context.Context, feedback FeedbackIntf) error

	// CancelTreatmentToTenants 取消方案分配
	CancelTreatmentToTenants(ctx context.Context, tenant_ids []string) error

	// ResetAppStaffPassword 重置 APP 登录密码
	ResetAppStaffPassword(ctx context.Context, user_id string, new_password string) error

	// UpdateSafePhone 更新User中tenant安全手机号
	UpdateTenantUserSafePhone(ctx context.Context, tid string, safePhone string, rev int32) error

	// UpdateTenantSafePhone 更新tenantEntity安全手机号
	UpdateTenantSafePhone(ctx context.Context, tid string, safePhone string, rev int32) error

	// GetTenantTreatments 获取商户方案
	GetTenantTreatments(ctx context.Context, tids []string) ([]TenantTreatmentIntf, error)

	// GetDataPrivilegePagination 分页获取数据权限
	GetDataPrivilegePagination(ctx context.Context, privilege_id string, size, offset int32) ([]string, int32, error)

	// RecoverStaff 恢复已删除的用户
	RecoverStaff(ctx context.Context, staffId, tenantId string, rev int32) error

	// GetOrganizationSubscriptionsMonthly 获取组织月份商户续费信息
	GetOrganizationSubscriptionsMonthly(ctx context.Context, organization_id string, start_time, end_time time.Time, size, offset int32) ([]TenantSubscriptionMonthlyIntf, int32, error)

	// ListAllOrganizationTenantIdsIncludeDeleted 获取组织下所有商户 ID（包括已删除）
	ListAllOrganizationTenantIdsIncludeDeleted(ctx context.Context, oids []string) ([]string, error)

	// GetMenusByPaths 通过路径获取目录
	GetMenusByPaths(ctx context.Context, paths []string) ([]MenuIntf, error)

	// CreateMenus 创建目录
	CreateMenus(ctx context.Context, privilege_id string, menu_ids []string) error

	// GetTreatmentTenant 获取被分配方案的商户
	GetTreatmentTenant(ctx context.Context, treatment_id string) (TenantTreatmentIntf, error)

	// DeletePrivilegeGroup 通过 ID 删除权限组
	DeletePrivilegeGroup(ctx context.Context, privilege_id string) error

	// UpdatePrivilegeGroup 修改权限组信息
	UpdatePrivilegeGroup(ctx context.Context, privilege_id, privilege_name, remark string) error

	// DeletePrivilegePages 删除页面权限
	DeletePrivilegePages(ctx context.Context, privilege_id string) error

	// DeletePrivilegeMenus 删除目录
	DeletePrivilegeMenus(ctx context.Context, privilege_id string) error

	// SearchTenantsByNameAndTime 通过商户名和时间模糊搜索商户
	SearchTenantsByNameAndTime(ctx context.Context, tids []string, tenant_name string, contact_name string, contact_phone string, status, offset, size int32, start_time, end_time *timestamppb.Timestamp) ([]TenantEntityIntf, int32, error)

	// GetUserByPrivilegeID 通过权限组 ID 获取用户
	GetUserByPrivilegeID(ctx context.Context, privilege_id string) (UserIntf, error)

	// BatchGetDataPrivileges 获取数据权限信息
	BatchGetDataPrivileges(ctx context.Context, privilege_id string) ([]PrivilegePolicyDataIntf, error)

	// DeleteDataPrivileges 删除数据权限
	DeleteDataPrivileges(ctx context.Context, ids []string) error

	// BatchGetTenantsSubscriptionPeriodsDetails 查询商户订阅详情
	BatchGetTenantsSubscriptionPeriodsDetails(ctx context.Context, tenant_ids []string, contact_name, contact_phone string, start_time time.Time, end_time time.Time, offset int32, size int32) ([]SubscriptionPeriodIntf, int64, int32, error)

	// ListPrivilegeGroupsPagination 获取权限组列表
	ListPrivilegeGroupsPagination(ctx context.Context, size, offset int32) ([]PrivilegePolicyIntf, int64, error)

	// UpdateOrganizationContactPhone 更新组织联系人电话
	UpdateOrganizationContactPhone(ctx context.Context, organization_id, organization_contact_phone string) error

	// ListFaqs 获取所有faq
	ListFaqs(ctx context.Context) ([]FaqIntf, error)
}
