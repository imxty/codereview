package token

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRefreshTokenExpire = errors.New("refresh token is expired")
)

// RefreshTokenDetails refresh_token详情
type RefreshTokenDetails struct {
	// refresh_token
	RefreshToken string
	// 用户信息
	UserInfo map[string]interface{}
	// 过期时间(UTC时间)
	ExpiredAt time.Time
}

// RefreshTokenStore refresh_token处理接口
type RefreshTokenStore interface {
	// 通过user_id生成refresh_token
	CreateRefreshToken(ctx context.Context, userInfo map[string]interface{}) (*RefreshTokenDetails, error)
	// 验证refresh_token是否存在并返回结果
	GetRefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenDetails, error)
	// 删除token,用token作为key
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	// 续订refreshToken
	RenewRefreshToken(ctx context.Context, refreshToken string) error
}
