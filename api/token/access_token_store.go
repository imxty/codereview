package token

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccessTokenExpire    = errors.New("access token is expired")
	ErrAccessTokenKickedOut = errors.New("access token has been kicked out")
)

// AccessTokenDetails access_token详情
type AccessTokenDetails struct {
	// access_token
	AccessToken string
	// 用户信息
	UserInfo map[string]interface{}
	// 过期时间(UTC时间)
	ExpiredAt time.Time
}

// AccessTokenStore access_token处理接口
type AccessTokenStore interface {
	// 创建access_token
	CreateAccessToken(ctx context.Context, userInfo map[string]interface{}) (*AccessTokenDetails, error)
	// 删除access_token
	DeleteAccessToken(ctx context.Context, accessToken string) error
	// 检测access_token是否存在并返回结果
	GetAccessTokenDetails(ctx context.Context, accessToken string) (*AccessTokenDetails, error)
	// 用于web，续订accessToken
	RenewAccessToken(ctx context.Context, accessToken string) error
}
