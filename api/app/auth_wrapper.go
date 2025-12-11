package app

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	gerr "errors"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jinmukeji/huimaibao-service/api"
	jjwt "github.com/jinmukeji/huimaibao-service/api/jwt"
	"github.com/jinmukeji/huimaibao-service/api/token"
	"github.com/jinmukeji/huimaibao-service/api/token/tokenstore"
	"github.com/jinmukeji/plat-pkg/v4/auth/jwt"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/jinmukeji/plat-pkg/v4/micro/meta"
	"go-micro.dev/v4/server"
)

const (
	// app_api请求service
	DeskApp = "com.shangyikangyou.huimaibao.api.app"
)

const (
	ContextUserID      = "user_id"
	ContextAccessToken = "access_token"
	ContextAddress     = "address"
)

var (
	// access token 白名单
	// 白名单中的 API 不验证 Access Token
	whiteList = map[string]bool{
		"UserAPI.RefreshAccessToken":                  true,
		"UserAPI.RefreshRefreshToken":                 true,
		"UserAPI.SignInByPassword":                    true,
		"NotificationAPI.SendVerificationCode":        true,
		"NotificationAPI.VerifyPhoneVerificationCode": true,
		"UserAPI.GetAndroidUpdateInfo":                true,
		"UserAPI.GetIOSUpdateInfo":                    true,
		"UserAPI.SignInBySmsCode":                     true,
		"UserAPI.ResetAppStaffPassword":               true,
		"UserAPI.GetSystemUpdateInfo":                 true,
	}
)

type AuthWrapper struct {
	tokenStore *token.TokenStore
	// jwt公钥内容
	jwtPublicKey *rsa.PublicKey
	// redis
	redisCli *redis.Client
}

func NewAuthWrapper(store *token.TokenStore, key *rsa.PublicKey) *AuthWrapper {
	return &AuthWrapper{
		tokenStore:   store,
		jwtPublicKey: key,
	}
}

// 验证jwt
func (wrapper *AuthWrapper) Auth(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if req.Service() != DeskApp {
			return errors.Error(codes.InvalidRequest, "invalid Request")
		}

		// 清理 Hijack Headers
		ctx = meta.Delete(ctx, ContextUserID)
		ctx = meta.Delete(ctx, ContextAccessToken)

		// 获取jwt
		bearToken, ok := meta.Get(ctx, "Authorization")
		if !ok {
			return errors.Error(api.ErrInvalidJWT, "invalid jwt")
		}
		strArr := strings.Split(bearToken, " ")
		if len(strArr) != 2 {
			return errors.Error(api.ErrInvalidJWT, "invalid jwt")
		}
		jtjwt := strArr[1]
		// 验证jwt
		opt := jwt.VerifyOption{
			MaxExpInterval: jjwt.DefaultMaxExpirationInterval,
			GetPublicKeyFunc: func(iss string) *rsa.PublicKey {
				// 判断issuer是否相同
				if iss != jjwt.SmIssuer {
					return nil
				}
				return wrapper.jwtPublicKey
			},
		}
		claims := &jjwt.JwtClaims{}
		// 获取token详情
		store := *(wrapper.tokenStore)
		// 验证jwt
		ok, err := jwt.RSAVerifyCustomJWT(jtjwt, opt, claims)
		if ok {
			// 判断是否是白名单的请求
			if _, ok := whiteList[req.Method()]; ok {
				return fn(ctx, req, resp)
			}
			if claims.AccessToken == "" {
				return errors.Error(codes.PermissionDenied, "invalid jwt: missing access_token")
			}
			// 获取地址
			address := base64.StdEncoding.EncodeToString([]byte(claims.Address))
			// 从claims中获取信息
			tk := claims.AccessToken
			// 获取token详情
			details, err := store.GetAccessTokenDetails(ctx, tk)
			if err != nil {
				if gerr.Is(err, token.ErrAccessTokenKickedOut) {
					return errors.Error(api.ErrKickedOut, err.Error())
				}
				return errors.Error(api.ErrInvalidAccessToken, api.ErrorChineseMsg(api.ErrInvalidAccessToken))
			}
			ctx = meta.Set(ctx, ContextUserID, details.UserInfo[tokenstore.UserID].(string))
			ctx = meta.Set(ctx, ContextAccessToken, details.AccessToken)
			ctx = meta.Set(ctx, ContextAddress, address)
			return fn(ctx, req, resp)
		} else {
			if err != nil {
				return errors.Error(codes.PermissionDenied, err.Error())
			}
			return nil
		}
	}
}

// 设置jwt公钥
func (wrapper *AuthWrapper) SetPublicKey(key *rsa.PublicKey) {
	wrapper.jwtPublicKey = key
}

// 设置redis
func (wrapper *AuthWrapper) SetRedisDsn(dsn string, db int) {
	c := redis.NewClient(&redis.Options{
		Addr: dsn,
		DB:   db,
	})
	wrapper.redisCli = c
}
