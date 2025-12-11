package organization

import (
	"context"
	"crypto/rsa"
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
	OrganizationApp = "com.shangyikangyou.huimaibao.api.organization"
)

const (
	ContextUserID      = "user_id"
	ContextAccessToken = "access_token"
)

var (
	// access token 白名单
	// 白名单中的 API 不验证 Access Token
	whiteList = map[string]bool{
		"UserAPI.RefreshAccessToken":                  true,
		"UserAPI.UploadImage":                         true,
		"UserAPI.SignInByUsername":                    true,
		"UserAPI.SignInByPhone":                       true,
		"UserAPI.SignUp":                              true,
		"UserAPI.ReCreateTenant":                      true,
		"UserAPI.GetUnAuthTenantRevision":             true,
		"UserAPI.ResetOrganizationPassword":           true,
		"NotificationAPI.SendPhoneVerificationCode":   true,
		"NotificationAPI.VerifyPhoneVerificationCode": true,
		"UserAPI.CheckOrganizationPhoneExist":         true,
		"UserAPI.CheckOrganizationUsernameExist":      true,
		"UserAPI.CreateTenant":                        true,
		"UserAPI.GetTenantStencil":                    true,
		"UserAPI.GetTenantStencilByID":                true,
		"UserAPI.GetWxSignature":                      true,
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
		if req.Service() != OrganizationApp {
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
		smjwt := strArr[1]
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
		ok, err := jwt.RSAVerifyCustomJWT(smjwt, opt, claims)
		if ok {
			// 判断是否是白名单的请求
			if _, ok := whiteList[req.Method()]; ok {
				return fn(ctx, req, resp)
			}
			if claims.AccessToken == "" {
				return errors.Error(codes.PermissionDenied, "invalid jwt: missing access_token")
			}
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
