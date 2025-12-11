package user

import (
	"context"
	gerr "errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinmukeji/go-pkg/v2/crypto/hash"
	notificationpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/user/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/rs/xid"
)

const (
	// 默认商户数量
	DefaultTenantLimit = 2000
)

// 注册组织管理请求
func (u *UserAPIHandler) SignUpOrganization(ctx context.Context, req *pb.SignUpOrganizationRequest, rsp *pb.SignUpOrganizationResponse) error {
	// 验证 request
	err := validateSignUpOrganizationRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 查询手机号是否注册
	ou, err := u.userStore.GetOrganizationByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if ou != nil {
		// 手机号存在
		return errors.Errorf(ErrPhoneHasBeenUsed, "organization phone[%s] has been used", req.GetPhone())
	}

	// 验证手机短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationpb.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		TxId:    req.GetTxId(),
		SmsCode: req.GetSmsCode(),
		Action:  notificationpb.TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION,
	})
	if err != nil {
		return err
	}
	// 生成组织 ID
	exist := true
	times := 0
	oid := ""
	for exist {
		// 如果循环 500 次还有冲突则报错
		if times == 500 {
			return errors.Error(codes.InvalidOperation, "generate organization id failed")
		}
		oid = generateOrganizationID()
		exist, err = u.userStore.CheckOrganizationIDExist(ctx, oid)
		if err != nil {
			return errors.Error(codes.DataAccessFailed, err.Error())
		}
		times++
	}

	// 创建组织
	organizationUser := &domain.User{
		UserID:                   xid.New().String(),
		OrganizationID:           oid,
		Username:                 req.GetUsername(),
		Nickname:                 req.GetUsername(),
		Phone:                    req.GetPhone(),
		OrganizationContactPhone: req.GetPhone(),
		RoleType:                 domain.RoleTypeOrganization,
		HashedPassword:           generateHashSHA256(req.GetPlainPassword()),
		TenantLimit:              DefaultTenantLimit,
	}
	// 创建组织
	err = u.userStore.CreateUser(ctx, organizationUser)
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	// 返回组织信息
	org := &pb.Organization{
		// 组织 ID
		OrganizationId: organizationUser.OrganizationID,
		// 名称
		Name: organizationUser.Username,
		// 手机号
		Phone: organizationUser.Phone,
	}
	rsp.Organization = org
	return nil
}

// 生成组织 ID
func generateOrganizationID() string {
	rand.NewSource(time.Now().UnixNano()) // 设置随机种子，确保每次运行都有不同的随机数
	randomInt := rand.Intn(1000000)       // 生成 0 到 999999 之间的随机整数
	return fmt.Sprintf("N%06d", randomInt)
}

// 验证 request
func validateSignUpOrganizationRequest(req *pb.SignUpOrganizationRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPlainPassword() == "" {
		return gerr.New("password should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	if req.GetUsername() == "" {
		return gerr.New("username should not be empty")
	}
	return nil
}

// 生成 SHA256 加密的密码
func generateHashSHA256(pwd string) string {
	pwdByte := hash.SHA256String(pwd)
	return hash.HexString(pwdByte)
}
