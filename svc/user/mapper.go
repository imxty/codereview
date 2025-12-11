package user

import (
	"strconv"

	customerv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toProtoEntityFromRevision
func toProtoEntityFromRevision(entity domain.TenantEntityRevisionIntf) *pb.Entity {
	if entity == nil {
		return nil
	}
	return &pb.Entity{
		// 租户 ID
		TenantId: entity.GetTenantID(),
		// 租户名称
		EntityName: entity.GetStoreName(),
		// logo 地址
		LogoUrl: entity.GetLogoUrl(),
		// 地址
		Address: toProtoAddress(entity.GetProvince(), entity.GetDistrict(), entity.GetDistrict(), entity.GetStreet()),
		// 手机号
		ContactPhone: entity.GetContactPhone(),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 营业执照地址
		BusinessLicenseUrl: entity.GetBusinessLicenseUrl(),
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		SafePhone:        entity.GetSafePhone(),
	}
}

// toProtoTenantEntityFromRevision
func toProtoTenantEntityFromRevision(entity domain.TenantEntityRevisionIntf) *pb.TenantEntity {
	if entity == nil {
		return nil
	}
	return &pb.TenantEntity{
		// 租户 ID
		TenantId: entity.GetTenantID(),
		// 租户名称
		Name: entity.GetStoreName(),
		// 地址
		Address: toProtoAddress(entity.GetProvince(), entity.GetCity(), entity.GetDistrict(), entity.GetStreet()),
		// 营业执照地址
		BusinessLicenseUrl: entity.GetBusinessLicenseUrl(),
		// logo 地址
		LogoUrl: entity.GetLogoUrl(),
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 联系人手机号
		ContactPhone: entity.GetContactPhone(),
		SafePhone:    entity.GetSafePhone(),
		CreateAt:     timestamppb.New(entity.GetCreatedAt()),
	}
}

// toProtoEntity
func toProtoEntity(entity domain.TenantEntityIntf) *pb.Entity {
	if entity == nil {
		return nil
	}
	return &pb.Entity{
		// 租户
		TenantId: entity.GetTenantID(),
		// 租户名称
		EntityName: entity.GetStoreName(),
		// logo 地址
		LogoUrl: entity.GetLogoUrl(),
		// 地址
		Address: toProtoAddress(entity.GetProvince(), entity.GetCity(), entity.GetDistrict(), entity.GetStreet()),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 联系人电话
		ContactPhone: entity.GetContactPhone(),
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		// 营业执照地址
		BusinessLicenseUrl: entity.GetBusinessLicenseUrl(),
		// 复查时间
		Overdue:   entity.GetOverdue(),
		SafePhone: entity.GetSafePhone(),
	}
}

func toAppSystemUser(u domain.UserIntf) *pb.SystemUser {
	if u == nil {
		return nil
	}
	return &pb.SystemUser{
		// 用户 ID
		UserId: u.GetUserID(),
		// 手机号
		Phone: u.GetPhone(),
		// 姓名
		Nickname: u.GetNickname(),
		Remark:   u.GetRemark(),
		RoleType: strconv.Itoa(int(u.GetRoleType())),
	}
}

// toProtoAddress 转为 proto 的 address
func toProtoAddress(province, city, district, street string) *pb.Address {
	return &pb.Address{
		// 省
		Province: province,
		// 市
		City: city,
		// 区
		District: district,
		// 街道
		Street: street,
	}
}

// toProtoTenant domain->proto
func toProtoTenant(s domain.TenantEntityIntf) *pb.Tenant {
	if s == nil {
		return nil
	}
	return &pb.Tenant{
		// 租户 ID
		TenantId: s.GetTenantID(),
		// 租户名称
		TenantName: s.GetStoreName(),
		// logo 地址
		LogoUrl: s.GetLogoUrl(),
		// 地址
		Address: toProtoAddress(s.GetProvince(), s.GetCity(), s.GetDistrict(), s.GetStreet()),
		// 联系电话
		ContactPhone: s.GetContactPhone(),
	}
}

// toProtoStaff domain->proto
func toProtoStaff(s domain.UserIntf) *pb.Staff {
	if s == nil {
		return nil
	}
	return &pb.Staff{
		StaffId:     s.GetUserID(),
		Phone:       s.GetPhone(),
		Name:        s.GetNickname(),
		IsActivated: s.GetIsActivated(),
		IsDeleted:   !s.GetIsActivated(),
	}
}

// toProtoStaffFromUser
func toProtoStaffFromUser(t domain.UserIntf) *pb.Staff {
	if t == nil {
		return nil
	}
	return &pb.Staff{
		// 员工 ID
		StaffId: t.GetUserID(),
		// 手机号
		Phone: t.GetPhone(),
		// 员工姓名
		Name:        t.GetNickname(),
		IsActivated: t.GetIsActivated(),
		IsDeleted:   !t.GetIsActivated(),
	}
}

// toProtoTenantEntity
func toProtoTenantEntity(entity domain.TenantEntityIntf) *pb.TenantEntity {
	if entity == nil {
		return nil
	}
	return &pb.TenantEntity{
		// 租户 ID
		TenantId: entity.GetTenantID(),
		// 租户名称
		Name: entity.GetStoreName(),
		// 地址
		Address: toProtoAddress(entity.GetProvince(), entity.GetCity(), entity.GetDistrict(), entity.GetStreet()),
		// 联系人姓名
		ContactName: entity.GetContactName(),
		// 联系人电话
		ContactPhone: entity.GetContactPhone(),
		// 营业执照地址
		BusinessLicenseUrl: entity.GetBusinessLicenseUrl(),
		// logo 地址
		LogoUrl: entity.GetLogoUrl(),
		// 社会信用代码
		SocialCreditCode: entity.GetSocialCreditCode(),
		SafePhone:        entity.GetSafePhone(),
		CreateAt:         timestamppb.New(entity.GetCreatedAt()),
	}
}

// toProtoTenantReviewStatus
func toProtoTenantReviewStatus(status int32) pb.TenantReviewStatus {
	switch status {
	case domain.TenantReviewStatusReviewing:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEWING
	case domain.TenantReviewStatusReviewSuccess:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_SUCCESS
	case domain.TenantReviewStatusReviewFailed:
		return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_REVIEW_FAIL
	}
	return pb.TenantReviewStatus_TENANT_REVIEW_STATUS_UNSET
}

// toProtoSubscriptionTimeline
func toProtoSubscriptionTimeline(t domain.SubscriptionTimelineIntf) *pb.TenantSubscriptionTimeline {
	if t == nil {
		return nil
	}
	return &pb.TenantSubscriptionTimeline{
		TenantId:  t.GetTenantID(),
		StartTime: timestamppb.New(t.GetStartTime()),
		EndTime:   timestamppb.New(t.GetEndTime()),
		Years:     t.GetYears(),
	}
}

// toPagePrivilege
func toPagePrivilege(p pb.PagePrivilege) string {
	switch p {
	case pb.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL:
		return "2"
	case pb.PagePrivilege_PAGE_PRIVILEGE_REVIEW:
		return "3"
	case pb.PagePrivilege_PAGE_PRIVILEGE_TENANT:
		return "4"
	case pb.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION:
		return "5"
	case pb.PagePrivilege_PAGE_PRIVILEGE_LIST:
		return "6"
	}
	return ""
}

func toProtoPrivilegeGroup(t domain.PrivilegePolicyIntf) *pb.PrivilegeGroup {
	return &pb.PrivilegeGroup{
		PrivilegeId:   t.GetPrivilegeID(),
		PrivilegeName: t.GetPrivilegeName(),
		Remark:        t.GetRemark(),
	}
}

func toProtoSystemUser(t domain.UserIntf) *pb.SystemUser {
	return &pb.SystemUser{
		UserId:    t.GetUserID(),
		Phone:     t.GetPhone(),
		Nickname:  t.GetNickname(),
		IsDeleted: !t.GetIsActivated(),
		Remark:    t.GetRemark(),
		RoleType:  strconv.Itoa(int(t.GetRoleType())),
	}
}

// toPrivilegePage
func toPrivilegePage(d domain.PrivilegePolicyPageIntf) pb.PagePrivilege {
	switch d.GetPresetID() {
	case "1":
		return pb.PagePrivilege_PAGE_PRIVILEGE_UNSET
	case "2":
		return pb.PagePrivilege_PAGE_PRIVILEGE_OPERATIONAL
	case "3":
		return pb.PagePrivilege_PAGE_PRIVILEGE_REVIEW
	case "4":
		return pb.PagePrivilege_PAGE_PRIVILEGE_TENANT
	case "5":
		return pb.PagePrivilege_PAGE_PRIVILEGE_SUBSCRIPTION
	case "6":
		return pb.PagePrivilege_PAGE_PRIVILEGE_LIST
	}
	return pb.PagePrivilege_PAGE_PRIVILEGE_INVALID
}

// toCustomerPagination
func toCustomerPagination(p *pb.Pagination) *customerv1.Pagination {
	if p == nil {
		return nil
	}

	return &customerv1.Pagination{
		Offset: p.GetOffset(),
		Size:   p.GetSize(),
	}
}

// toAppMenu
func toAppMenu(d domain.MenuIntf) *pb.Menu {
	if d == nil {
		return &pb.Menu{}
	}

	return &pb.Menu{
		MenuId:  d.GetMenuID().String(),
		Title:   d.GetTitle(),
		Path:    d.GetPath(),
		Sort:    d.GetSort(),
		Visible: d.GetVisible(),
	}
}

// toMenuPath
func toMenuPath(p string) string {
	switch p {
	case "2":
		return "/"
	case "3":
		return "/examine"
	case "4":
		return "/check"
	case "5":
		return "/tissuelist"
	case "6":
		return "/advancepayment"
	}
	return ""
}
