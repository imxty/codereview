package notification

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/notification/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

const (
	// ErrSmsNotFound
	ErrSmsNotFound = 5040
)

func (s *NotificationAPIHandler) CheckTxIdUsage(ctx context.Context, req *pb.CheckTxIdUsageRequest, rsp *pb.CheckTxIdUsageResponse) error {
	err := validateCheckTxIdUsageRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	// 通过 txid 查询短信
	sms, err := s.notificationStore.GetSmsByTxId(ctx, req.GetTxId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if sms == nil {
		return errors.Errorf(ErrSmsNotFound, "sms not found by txid:[%s]", req.GetTxId())
	}

	// 返回 action
	rsp.Action = pb.TemplateAction(pb.TemplateAction_value[sms.GetTemplateAction()])
	return nil
}

// 验证 request
func validateCheckTxIdUsageRequest(req *pb.CheckTxIdUsageRequest) error {
	if req.GetTxId() == "" {
		return gerr.New("tx id should not be empty")
	}
	return nil
}
