package tenant

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/tenant/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TenantAPIHandler) SignInTenantAdmin(ctx context.Context, req *pb.SignInTenantAdminRequest, rsp *pb.SignInTenantAdminResponse) error {
	// 验证request
	err := validateSignInTenantAdminRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 登录商户后台
	signInRsp, err := s.userAPI.SignInTenantAdmin(ctx, &userpb.SignInTenantAdminRequest{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	rsp.Entity = toAppTenantEntity(signInRsp.GetEntity(), s.s3Domain)
	rsp.ReportSharingSwitchStatus = rsp.GetReportSharingSwitchStatus()

	// 生成token
	tenantId := signInRsp.GetEntity().GetTenantId()
	userInfo := map[string]interface{}{
		tokenstore.UserID: tenantId,
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
	return nil
}

// 验证request
func validateSignInTenantAdminRequest(req *pb.SignInTenantAdminRequest) error {
	if req.GetPhone() == "" {
		return gerr.New("phone should not be empty")
	}
	if req.GetPassword() == "" {
		return gerr.New("password should not be empty")
	}
	return nil
}
