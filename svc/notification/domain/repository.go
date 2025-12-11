package domain

import (
	"context"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

type NotificationRepository interface {
	dbutils.Tx

	// GetSmsByTxId 通过 txid 获取短信
	GetSmsByTxId(ctx context.Context, txId string) (SmsIntf, error)
	// GetLatestSms 获取最新短信
	GetLatestSms(ctx context.Context, phone string) (SmsIntf, error)
	// GetLatestSmsByAction 获取最新短信
	GetLatestSmsByAction(ctx context.Context, phone, action string) (SmsIntf, error)
	// CreateSmsRecord 创建短信记录
	CreateSmsRecord(ctx context.Context, sms SmsIntf) error
}
