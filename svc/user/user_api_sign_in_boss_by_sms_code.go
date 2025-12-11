package user

import (
	"context"
	gerr "errors"
	"sort"

	notificationv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	ErrHasNoPrivileges = 5049
)

// 验证码登录金姆平台
func (u *UserAPIHandler) SignInBossBySmsCode(ctx context.Context, req *pb.SignInBossBySmsCodeRequest, rsp *pb.SignInBossBySmsCodeResponse) error {
	err := validateSignInBossBySmsCodeRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 验证手机短信
	_, err = u.notificationAPI.VerifyPhoneVerificationCode(ctx, &notificationv1.VerifyPhoneVerificationCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
		Action:  notificationv1.TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS,
	})
	if err != nil {
		return err
	}

	// 获取系统用户
	us, err := u.userStore.GetSystemUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if us == nil {
		return errors.Errorf(ErrSystemUserNotFound, "user[%s] not found", req.GetPhone())
	}

	// 获取菜单
	menus, err := u.userStore.GetMenus(ctx, us.GetPrivilegeID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果没有可访问的目录，不能登录
	if len(menus) < 1 {
		return errors.Errorf(ErrHasNoPrivileges, "system user [%s] has no privileges", us.GetUserID())
	}

	resultMenus := make([]*pb.Menu, len(menus))
	for k, v := range menus {
		resultMenus[k] = toAppMenu(v)
	}

	sort.Slice(resultMenus, func(i, j int) bool {
		return resultMenus[i].GetSort() < resultMenus[j].GetSort()
	})

	group, err := u.userStore.GetPrivilegeGroupById(ctx, us.GetPrivilegeID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if group == nil {
		return errors.Errorf(ErrPrivilegeGroupNotExist, "group[%s] not exist", us.GetPrivilegeID())
	}

	rsp.SystemUser = toAppSystemUser(us)
	rsp.Menus = resultMenus
	rsp.SystemUser.PrivilegeGroup = &pb.PrivilegeGroup{
		PrivilegeId: group.GetPrivilegeID(),
	}

	return nil
}

// 验证 request
func validateSignInBossBySmsCodeRequest(req *pb.SignInBossBySmsCodeRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetSmsCode() == "" {
		return gerr.New("sms_code should not be empty")
	}
	if req.GetTxId() == "" {
		return gerr.New("tx_id should not be empty")
	}
	return nil
}
