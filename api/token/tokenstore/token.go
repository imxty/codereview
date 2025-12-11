package tokenstore

import (
	"context"
	gerr "errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jinmukeji/huimaibao-service/api/token"
)

const (
	UserID   = "user_id"
	TenantID = "tenant_id"
)

// access_token 生成与存储
type tokenStore struct {
	client           *redis.Client
	atExpireInterval time.Duration
	rtExpireInterval time.Duration
}

var _ token.TokenStore = (*tokenStore)(nil)

// NewTokenStore 创建一个access_token操作对象
func NewTokenStore(redisDSN string, db int, atInterval time.Duration, rtInterval time.Duration) token.TokenStore {
	c := redis.NewClient(&redis.Options{
		Addr: redisDSN,
		DB:   db,
	})
	pong, err := c.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)
	return &tokenStore{
		client:           c,
		atExpireInterval: atInterval,
		rtExpireInterval: rtInterval,
	}
}

// 通过user_id生成access_token
// [at:xxxxxxxxx] : userID
func (store *tokenStore) CreateAccessToken(ctx context.Context, userInfo map[string]interface{}) (*token.AccessTokenDetails, error) {

	// 生成uuid
	atID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	// 生成存入redis的key
	atKey := store.buildAtKey(atID.String())

	expiredAt := time.Now().Add(store.atExpireInterval).UTC()
	_, err = store.client.HSet(ctx, atKey, userInfo).Result()
	// 查看结果
	if err != nil {
		return nil, err
	}
	_, err = store.client.Expire(ctx, atKey, store.atExpireInterval).Result()
	// 查看结果
	if err != nil {
		return nil, err
	}
	// 返回结果
	details := &token.AccessTokenDetails{
		AccessToken: atID.String(),
		UserInfo:    userInfo,
		ExpiredAt:   expiredAt,
	}
	return details, nil
}

// 验证access_token是否存在
// 1.获取指定 AT，如果过期,返回nil, err
// 2.黑名单是否存在
func (store *tokenStore) GetAccessTokenDetails(ctx context.Context, accessToken string) (*token.AccessTokenDetails, error) {
	// 获取拼接后的access_token key
	atKey := store.buildAtKey(accessToken)
	// 1.获取指定 AT，如果过期,返回nil, err
	// 获取userID
	userInfo, err := store.client.HGetAll(ctx, atKey).Result()
	if err != nil {
		return nil, err
	}
	info := make(map[string]interface{}, len(userInfo))
	for k, v := range userInfo {
		info[k] = v
	}
	// 获取过期时间
	issuedAt := time.Now()
	remainingInterval, err := store.client.TTL(ctx, atKey).Result()
	if err != nil {
		return nil, err
	}
	// 如果已经过期(-2是过期，-1是未设置过期时间)
	if remainingInterval <= time.Duration(-2) {
		return nil, token.ErrAccessTokenExpire
	}
	if remainingInterval == time.Duration(-1) {
		return nil, gerr.New("invalid access token")
	}
	// 获取创建时间
	expire := issuedAt.Add(remainingInterval).Add(-store.atExpireInterval)

	atDetails := &token.AccessTokenDetails{
		UserInfo:    info,
		AccessToken: accessToken,
		ExpiredAt:   issuedAt.Add(remainingInterval),
	}
	// 拼接黑名单的key
	blackUserID := store.getBlockedUserKey(userInfo[UserID])
	// 检测是否存在黑名单
	result, err := store.client.Get(ctx, blackUserID).Int64()
	if gerr.Is(err, redis.Nil) {
		// 如果黑名单上没有记录
		return atDetails, nil
	}
	if err != nil {
		return nil, err
	}
	// 获取退出时间
	kickedOutAt := time.Unix(result, 0)
	// 如果创建时间在退出时间之前，说明已经无效
	if expire.Before(kickedOutAt) {
		return nil, token.ErrAccessTokenKickedOut
	}
	return atDetails, nil

}

// 删除token
// 用户退出
func (store *tokenStore) DeleteAccessToken(ctx context.Context, accessToken string) error {
	// 获取拼接后的access_token
	atKey := store.buildAtKey(accessToken)
	// 删除token
	num, err := store.client.Del(ctx, atKey).Result()
	if err != nil {
		return err
	}
	// 如果没有对任何数据造成影响，则说明token无效
	if num == 0 {
		return gerr.New("invalid access token")
	}
	return nil
}

// 获取token存入redis的key
func (store *tokenStore) buildAtKey(at string) string {
	// 拼接key
	atKey := fmt.Sprintf("at:%s", at)
	return atKey
}

// 获取token存入redis的key
func (store *tokenStore) getBlockedUserKey(userID string) string {
	// 拼接黑名单的key
	blackUserID := fmt.Sprintf("bl:%s", userID)
	return blackUserID
}

// 通过user_id生成refresh_token
func (store *tokenStore) CreateRefreshToken(ctx context.Context, userInfo map[string]interface{}) (*token.RefreshTokenDetails, error) {
	// 生成uuid
	rt, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	rtKey := store.buildRtKey(rt.String())
	expire := time.Now().Add(store.rtExpireInterval).UTC()
	// 插入数据
	_, err = store.client.HSet(ctx, rtKey, userInfo).Result()
	if err != nil {
		return nil, err
	}
	_, err = store.client.Expire(ctx, rtKey, store.rtExpireInterval).Result()
	// 查看结果
	if err != nil {
		return nil, err
	}
	// 返回结果
	details := &token.RefreshTokenDetails{
		RefreshToken: rt.String(),
		UserInfo:     userInfo,
		ExpiredAt:    expire,
	}
	return details, nil
}

// 验证refresh_token是否存在并返回结果
func (store *tokenStore) GetRefreshToken(ctx context.Context, refreshToken string) (*token.RefreshTokenDetails, error) {
	// 获取拼接后的refresh_token key
	rtKey := store.buildRtKey(refreshToken)
	// 1.获取指定 RT，如果过期,返回nil, err
	// 获取userID
	userInfo, err := store.client.HGetAll(ctx, rtKey).Result()
	if err != nil {
		return nil, err
	}
	info := make(map[string]interface{}, len(userInfo))
	for k, v := range userInfo {
		info[k] = v
	}
	// 获取过期时间
	issuedAt := time.Now()
	remainingInterval, err := store.client.TTL(ctx, rtKey).Result()
	if err != nil {
		return nil, err
	}
	// 如果已经过期(-2是过期，-1是未设置过期时间)
	if remainingInterval <= time.Duration(-2) {
		return nil, gerr.New("refresh token is expired")
	}
	if remainingInterval == time.Duration(-1) {
		return nil, gerr.New("invalid refresh token")
	}
	// 获取创建时间
	expire := issuedAt.Add(remainingInterval).Add(-store.rtExpireInterval)

	rtDetails := &token.RefreshTokenDetails{
		UserInfo:     info,
		RefreshToken: refreshToken,
		ExpiredAt:    issuedAt.Add(remainingInterval),
	}
	// 拼接黑名单的key
	blackUserID := store.getBlockedUserKey(userInfo[UserID])
	// 检测是否存在黑名单
	result, err := store.client.Get(ctx, blackUserID).Int64()
	// 如果黑名单上没有记录
	if gerr.Is(err, redis.Nil) {
		return rtDetails, nil
	}
	// 如果是其他错误
	if err != nil {
		return nil, err
	}
	// 获取退出时间
	kickedOutAt := time.Unix(result, 0)
	// 如果创建时间在退出时间之前，说明已经无效
	if expire.Before(kickedOutAt) {
		return nil, gerr.New("refresh token is expired")
	}
	return rtDetails, nil
}

// 删除refresh_token,用token的uuid作为key
func (store *tokenStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	// 获取rt key
	rtKey := store.buildRtKey(refreshToken)
	// 从redis中删除
	result, err := store.client.Del(ctx, rtKey).Result()
	if err != nil {
		return err
	}
	// 如果对0条数据造成影响说明没有该token
	if result == 0 {
		return gerr.New("invalid refresh token")
	}
	return nil
}

// 管理员踢出用户，加入黑名单
func (store *tokenStore) KickOutUser(ctx context.Context, userID string) error {
	// 获取黑名单拼接的用户ID的key
	blackUserID := store.getBlockedUserKey(userID)
	// 存入当前时间，不过期
	_, err := store.client.Set(ctx, blackUserID, time.Now().Unix(), 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// 批量踢出用户
func (store *tokenStore) KickOutUsers(ctx context.Context, userIds []string) error {
	pipe := store.client.TxPipeline()
	t := time.Now().Unix()
	for _, userID := range userIds {
		// 获取黑名单拼接的用户ID的key
		blackUserID := store.getBlockedUserKey(userID)
		// 存入当前时间，不过期
		_, err := pipe.Set(ctx, blackUserID, t, 0).Result()
		if err != nil {
			return err
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 用于web，续订accessToken
func (store *tokenStore) RenewAccessToken(ctx context.Context, accessToken string) error {
	atKey := store.buildAtKey(accessToken)
	// 1.获取指定 AT，如果过期,返回nil
	// 查看是否存在
	_, err := store.client.HGetAll(ctx, atKey).Result()
	if err != nil {
		return err
	}
	// 获取过期时间
	remainingInterval, err := store.client.TTL(ctx, atKey).Result()
	if err != nil {
		return err
	}
	// 如果已经过期(-2是过期，-1是未设置过期时间)
	if remainingInterval <= time.Duration(-2) {
		return token.ErrAccessTokenExpire
	}
	if remainingInterval == time.Duration(-1) {
		return gerr.New("invalid access token")
	}
	// 刷新为一个access_token过期周期
	result := store.client.Expire(ctx, atKey, store.atExpireInterval)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

// 续订refreshToken
func (store *tokenStore) RenewRefreshToken(ctx context.Context, refreshToken string) error {
	rtKey := store.buildRtKey(refreshToken)
	// 1.获取指定 RT，如果过期,返回nil
	// 查看是否存在
	_, err := store.client.HGetAll(ctx, rtKey).Result()
	if err != nil {
		return err
	}
	// 获取过期时间
	remainingInterval, err := store.client.TTL(ctx, rtKey).Result()
	if err != nil {
		return err
	}
	// 如果已经过期(-2是过期，-1是未设置过期时间)
	if remainingInterval <= time.Duration(-2) {
		return token.ErrRefreshTokenExpire
	}
	if remainingInterval == time.Duration(-1) {
		return gerr.New("invalid refresh token")
	}
	// 刷新为一个refresh_token过期周期
	result := store.client.Expire(ctx, rtKey, store.rtExpireInterval)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

// 获取token存入redis的key
func (store *tokenStore) buildRtKey(rt string) string {
	// 拼接key
	rtKey := fmt.Sprintf("rt:%s", rt)
	return rtKey
}
