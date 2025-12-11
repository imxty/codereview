package token

import "context"

// TokenStore Token处理
type TokenStore interface {
	// access_token处理
	AccessTokenStore
	// refresh_token处理
	RefreshTokenStore
	// 踢出用户
	KickOutUser(ctx context.Context, userID string) error
	// 批量踢出用户
	KickOutUsers(ctx context.Context, userID []string) error
}
