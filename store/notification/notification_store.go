package notification

import (
	"context"
	"errors"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/huimaibao-service/svc/notification/domain"
	"gorm.io/gorm"
)

type NotificationStore struct {
	// 数据库管理
	*dbutils.Connection
}

func NewNotificationStore(management *dbutils.Connection) *NotificationStore {
	return &NotificationStore{
		management,
	}
}

// 实现 Domain 行为
var _ domain.NotificationRepository = (*NotificationStore)(nil)

// GetSmsByTxId 通过 txid 获取短信
func (n *NotificationStore) GetSmsByTxId(ctx context.Context, txId string) (domain.SmsIntf, error) {
	var sm Sms
	err := n.GetConnection(ctx).Model(&Sms{}).Where("tx_id = ?", txId).First(&sm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sm.ToDomainSms(), nil
}

// GetLatestSms 获取最新短信
func (n *NotificationStore) GetLatestSms(ctx context.Context, phone string) (domain.SmsIntf, error) {
	var sm Sms
	db := n.GetConnection(ctx)
	err := db.Model(&Sms{}).Where("phone = ?", phone).Order("created_at desc").First(&sm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sm.ToDomainSms(), nil
}

// GetLatestSmsByAction 获取最新短信
func (n *NotificationStore) GetLatestSmsByAction(ctx context.Context, phone, action string) (domain.SmsIntf, error) {
	var sm Sms
	db := n.GetConnection(ctx)
	err := db.Model(&Sms{}).Where("phone = ? and template_action = ?", phone, action).Order("created_at desc").First(&sm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sm.ToDomainSms(), nil
}

// CreateSmsRecord 创建短信记录
func (n *NotificationStore) CreateSmsRecord(ctx context.Context, sms domain.SmsIntf) error {
	var sm Sms
	sm.FromDomainSms(sms)
	db := n.GetConnection(ctx)
	return db.Model(&Sms{}).Create(&sm).Error
}
