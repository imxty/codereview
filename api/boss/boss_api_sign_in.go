package boss

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/boss/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *BossAPIHandler) SignIn(ctx context.Context, req *pb.SignInRequest, rsp *pb.SignInResponse) error {
	err := validateSignInRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	signRsp, err := s.userAPI.SignInBoss(ctx, &userv1.SignInBossRequest{
		Password: req.GetPassword(),
		Phone:    req.GetPhone(),
	})
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	sysUser := signRsp.GetSystemUser()
	// 生成token
	userInfo := map[string]interface{}{
		tokenstore.UserID: sysUser.GetUserId(),
	}
	at, err := s.tokenStore.CreateAccessToken(ctx, userInfo)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorMsg(err))
	}
	rt, err := s.tokenStore.CreateRefreshToken(ctx, userInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrCreateRefreshToken, err, api.ErrorMsg(err))
	}

	atExpireAt := timestamppb.New(at.ExpiredAt)
	rtExpireAt := timestamppb.New(rt.ExpiredAt)

	rsp.AccessToken = &pb.TokenDetails{
		Token:       at.AccessToken,
		ExpiredTime: atExpireAt,
	}
	rsp.RefreshToken = &pb.TokenDetails{
		Token:       rt.RefreshToken,
		ExpiredTime: rtExpireAt,
	}

	rsp.SystemUser = toAppSystemUser(sysUser)
	rsp.SystemUser.PrivilegeGroup = &pb.PrivilegeGroup{
		PrivilegeId: sysUser.GetPrivilegeGroup().GetPrivilegeId(),
	}

	resultMenus := make([]*pb.Menu, len(signRsp.GetMenus()))
	for k, v := range signRsp.GetMenus() {
		resultMenus[k] = &pb.Menu{
			MenuId:  v.GetMenuId(),
			Title:   v.GetTitle(),
			Path:    v.GetPath(),
			Sort:    v.GetSort(),
			Visible: v.GetVisible(),
		}
	}
	rsp.Menus = resultMenus

	return nil
}

// 验证request
func validateSignInRequest(req *pb.SignInRequest) error {
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	return nil
}
