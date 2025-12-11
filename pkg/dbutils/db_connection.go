package dbutils

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 实现事务控制接口 Tx
type txDbCtxKey struct{}

type Tx interface {
	// 开启事务
	BeginTx(ctx context.Context) context.Context
	// 提交事务
	CommitTx(ctx context.Context)
	// 回滚事务
	RollbackTx(ctx context.Context)
}

// Connection 数据库连接
type Connection struct {
	db *gorm.DB
}

// NewConnection 创建数据库连接
func NewConnection(db *gorm.DB) *Connection {
	// 设置日志级别为Info
	db.Config.Logger = logger.Default.LogMode(logger.Info)
	return &Connection{
		db: db,
	}
}

// GetConnection 获取连接
func (h *Connection) GetConnection(ctx context.Context) *gorm.DB {
	if v := ctx.Value(txDbCtxKey{}); v != nil {
		return v.(*gorm.DB)
	}
	return h.db
}

// BeginTx 开启事务
func (h *Connection) BeginTx(ctx context.Context) context.Context {
	txDB := h.db.Begin()
	// 从原始的 DB 开启事务
	return context.WithValue(ctx, txDbCtxKey{}, txDB)
}

// CommitTx 提交事务
func (h *Connection) CommitTx(ctx context.Context) {
	h.GetConnection(ctx).Commit()
}

// RollbackTx 回滚事务
func (h *Connection) RollbackTx(ctx context.Context) {
	h.GetConnection(ctx).Rollback()
}
