package jwt

import (
	"time"

	"github.com/jinmukeji/plat-pkg/v4/auth/jwt"
)

const (
	// DefaultMaxExpirationInterval 默认最大的过期时间间隔（10分钟）
	DefaultMaxExpirationInterval = 10 * time.Minute
	// 慧脉宝JWT issuer
	SmIssuer = "JinmuHealth"
)

// JwtClaims 慧脉宝药jwt的claims
type JwtClaims struct {
	jwt.StandardClaims

	AccessToken string `json:"access_token"`
	Address     string `json:"address"`
}
