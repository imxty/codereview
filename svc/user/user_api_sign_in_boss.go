package user

import (
	"context"
	gerr "errors"
	"sort"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrSystemUserNotFound
	ErrSystemUserNotFound = 5015
	// ErrSysPasswordWrong
	ErrSysPasswordWrong = 5016
)

// 登录金姆平台运营
func (u *UserAPIHandler) SignInBoss(ctx context.Context, req *pb.SignInBossRequest, rsp *pb.SignInBossResponse) error {
	// 1.验证request
	err := validateSignInBossRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	// 通过手机号查询用户
	user, err := u.userStore.GetSystemUserByPhone(ctx, req.GetPhone())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if user == nil {
		return errors.Errorf(ErrSystemUserNotFound, "system user not found by phone %s", req.GetPhone())
	}

	// 生成hash密码
	hashedPassword := generateHashSHA256(req.GetPassword())
	if user.GetHashedPassword() != hashedPassword {
		return errors.Error(ErrSysPasswordWrong, "password wrong")
	}

	// 获取菜单
	menus, err := u.userStore.GetMenus(ctx, user.GetPrivilegeID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 如果没有可访问的目录，不能登录
	if len(menus) < 1 {
		return errors.Errorf(ErrHasNoPrivileges, "system user [%s] has no privileges", user.GetUserID())
	}

	resultMenus := make([]*pb.Menu, len(menus))
	for k, v := range menus {
		resultMenus[k] = toAppMenu(v)
	}

	sort.Slice(resultMenus, func(i, j int) bool {
		return resultMenus[i].GetSort() < resultMenus[j].GetSort()
	})

	group, err := u.userStore.GetPrivilegeGroupById(ctx, user.GetPrivilegeID())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	if group == nil {
		return errors.Errorf(ErrPrivilegeGroupNotExist, "group[%s] not exist", user.GetPrivilegeID())
	}

	rsp.Menus = resultMenus
	rsp.SystemUser = toAppSystemUser(user)
	rsp.SystemUser.PrivilegeGroup = &pb.PrivilegeGroup{
		PrivilegeId: group.GetPrivilegeID(),
	}

	return nil
}

// 验证request
func validateSignInBossRequest(req *pb.SignInBossRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
