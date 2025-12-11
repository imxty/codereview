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

// 验证码登录
func (s *BossAPIHandler) SignInBySms(ctx context.Context, req *pb.SignInBySmsRequest, rsp *pb.SignInBySmsResponse) error {
	err := validateSignInBySmsRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	signRsp, err := s.userAPI.SignInBossBySmsCode(ctx, &userv1.SignInBossBySmsCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.SystemUser = toApiSystemUser(signRsp.GetSystemUser())
	rsp.Menus = rsp.GetMenus()

	// 生成 token
	sysUser := signRsp.GetSystemUser()
	userId := rsp.GetSystemUser().GetUserId()
	userInfo := map[string]interface{}{
		tokenstore.UserID: userId,
	}
	at, err := s.tokenStore.CreateAccessToken(ctx, userInfo)
	if err != nil {
		return errors.Error(api.ErrCreateAccessToken, api.ErrorMsg(err))
	}
	rt, err := s.tokenStore.CreateRefreshToken(ctx, userInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrCreateRefreshToken, err, api.ErrorMsg(err))
	}

	atExpiredAt := timestamppb.New(at.ExpiredAt)
	rtExpiredAt := timestamppb.New(rt.ExpiredAt)

	rsp.AccessToken = &pb.TokenDetails{
		Token:       at.AccessToken,
		ExpiredTime: atExpiredAt,
	}
	rsp.RefreshToken = &pb.TokenDetails{
		Token:       rt.RefreshToken,
		ExpiredTime: rtExpiredAt,
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

// 验证 request
func validateSignInBySmsRequest(req *pb.SignInBySmsRequest) error {
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
