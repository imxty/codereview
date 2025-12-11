package app

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AppAPIHandler) SignInBySmsCode(ctx context.Context, req *pb.SignInBySmsCodeRequest, rsp *pb.SignInBySmsCodeResponse) error {
	// 1.验证request
	err := validateSignInBySmsCodeRequest(req)
	if err != nil {
		return errors.Error(api.ErrInvalidRequest, err.Error())
	}

	// 发送获取员工信息的请求
	getStaffRsp, err := s.userAPI.SignInAppBySmsCode(ctx, &userpb.SignInAppBySmsCodeRequest{
		Phone:   req.GetPhone(),
		SmsCode: req.GetSmsCode(),
		TxId:    req.GetTxId(),
	})
	if err != nil {
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 构造和token关联的信息
	tenantId := getStaffRsp.GetStaff().GetTenant().GetTenantId()
	userInfo := map[string]interface{}{
		tokenstore.UserID:   getStaffRsp.GetStaff().GetStaffId(),
		tokenstore.TenantID: tenantId,
	}
	// 生成token
	at, err := s.tokenStore.CreateAccessToken(ctx, userInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrCreateAccessToken, err, api.ErrorMsg(err))
	}
	rt, err := s.tokenStore.CreateRefreshToken(ctx, userInfo)
	if err != nil {
		return errors.ErrorWithCause(api.ErrCreateRefreshToken, err, api.ErrorMsg(err))
	}

	// 返回数据
	rsp.Staff = toAppStaff(getStaffRsp.GetStaff(), s.s3Domain)
	// 返回token
	rtExpireAt := timestamppb.New(rt.ExpiredAt)
	rsp.RefreshToken = &pb.TokenDetails{
		Token:       rt.RefreshToken,
		ExpiredTime: rtExpireAt,
	}

	atExpireAt := timestamppb.New(at.ExpiredAt)
	rsp.AccessToken = &pb.TokenDetails{
		Token:       at.AccessToken,
		ExpiredTime: atExpireAt,
	}
	return nil
}

// 验证request
func validateSignInBySmsCodeRequest(req *pb.SignInBySmsCodeRequest) error {
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
