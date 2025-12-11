package client

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/google/uuid"
)

// TemplateLanguage 模块语言
type TemplateLanguage int32

const (
	// SimpleChinese 简体中文
	SimpleChinese TemplateLanguage = 0
)

const (
	DefaultNationCode = "+86"
)

// TemplateAction 短信验证码类型
type TemplateAction int32

const (
	// 组织注册
	TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION TemplateAction = 2
	// 组织登陆
	TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION TemplateAction = 3
	// 组织重置密码
	TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD TemplateAction = 4
	// 组织绑定手机号
	TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE TemplateAction = 5
	// 组织账单
	TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING TemplateAction = 6
	// 商户创建成功
	TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS TemplateAction = 7
	// 商户创建失败
	TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL TemplateAction = 8
	// 商户绑定手机号
	TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE TemplateAction = 9
	// 商户登陆
	TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT TemplateAction = 10
	// 商户修改手机号
	TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE TemplateAction = 11
	// 商户重置密码
	TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD TemplateAction = 12
	// 员工手机号登录
	TemplateAction_TEMPLATE_ACTION_SIGNIN_APP TemplateAction = 13
	// 组织修改手机号
	TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE TemplateAction = 14
	// 金姆平台登录
	TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS TemplateAction = 15
	// 重置员工 App 登录密码
	TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD = 16
	// 审核通知
	TemplateAction_TEMPLATE_ACTION_REVIEW_NOTIFICATION TemplateAction = 17
)

const (
	// DysmsAPIEndpoint 阿里云短信接口URL
	DysmsAPIEndpoint = "http://dysmsapi.aliyuncs.com"
)

// simpleChineseTemplateAction 国内简体模版
var simpleChineseTemplateAction = map[TemplateAction]string{
	// 组织注册
	TemplateAction_TEMPLATE_ACTION_SIGNUP_ORGANIZATION: "SMS_498750744",
	// 组织登录
	TemplateAction_TEMPLATE_ACTION_SIGNIN_ORGANIZATION: "SMS_498705718",
	// 组织重置密码
	TemplateAction_TEMPLATE_ACTION_RESET_ORGANIZATION_PASSWORD: "SMS_498850657",
	// 组织绑定手机号
	TemplateAction_TEMPLATE_ACTION_BIND_ORGANIZATION_PHONE: "SMS_498745704",
	// 组织账单
	TemplateAction_TEMPLATE_ACTION_ORGANIZATION_BILLING: "SMS_498725729",
	// 商户创建成功
	TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS: "SMS_498770673",
	// 商户创建失败
	TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL: "SMS_498765666",
	// 商户绑定手机号
	TemplateAction_TEMPLATE_ACTION_BIND_TENANT_PHONE: "SMS_498795691",
	// 商户登陆
	TemplateAction_TEMPLATE_ACTION_SIGNIN_TENANT: "SMS_498705718",
	// 商户修改手机号
	TemplateAction_TEMPLATE_ACTION_MODIFY_TENANT_PHONE: "SMS_498780700",
	// 商户重置密码
	TemplateAction_TEMPLATE_ACTION_RESET_TENANT_PASSWORD: "SMS_498850657",
	// 员工手机号登录
	TemplateAction_TEMPLATE_ACTION_SIGNIN_APP: "SMS_498705718",
	// 组织修改手机号
	TemplateAction_TEMPLATE_ACTION_MODIFY_ORGANIZATION_PHONE: "SMS_498780700",
	// 审核通知
	TemplateAction_TEMPLATE_ACTION_REVIEW_NOTIFICATION: "SMS_498740699",
	// 登录金姆平台
	TemplateAction_TEMPLATE_ACTION_SIGNIN_BOSS: "SMS_498705718",
	// 重置员工 App 登录密码
	TemplateAction_TEMPLATE_ACTION_RESET_APP_PASSWORD: "SMS_498850657",
}

// SendSmsReply 发送短信返回
type SendSmsReply struct {
	Code      string `json:"Code,omitempty"`      // 状态码
	Message   string `json:"Message,omitempty"`   // 状态码的描述
	RequestID string `json:"RequestId,omitempty"` // 请求 ID
	BizID     string `json:"BizId,omitempty"`     // 发送回执 ID
}

// TemplateParam 短信模版变量
type TemplateParam struct {
	Code string `json:"code,omitempty"` // 对应模板中的 ${code}
}

// TemplateEnvParam 短信模板变量
type TemplateEnvParam struct {
	Code string `json:"code,omitempty"` // 对应模板中的 ${code}
	Env  string `json:"env,omitempty"`  // 对应模板中的 ${env}
}

// TemplateTidParam 短信模板变量
type TemplateTidParam struct {
	TId string `json:"tid,omitempty"` // 对应模版中的 ${tid}
	Env string `json:"env,omitempty"` // 对应模版中的 ${env}
}

// TemplateEnv 短信模板变量
type TemplateEnv struct {
	Env string `json:"env,omitempty"` // 对应模板中的 ${env}
}

// TemplateReviewNotification 短信模板变量
type TemplateReviewNotification struct {
	Issue  string `json:"issue,omitempty"`  // 对应模板中的 ${issue}
	Number string `json:"number,omitempty"` // 对应模板中的 ${number}
}

// AliyunSMSClient 阿里云 SMS 客户端
type AliyunSMSClient struct {
	AccessKeyID     string
	AccessKeySecret string
	Client          *dysmsapi.Client
	Env             string
}

// NewAliyunSMSClient 生成 SMSClient
func NewAliyunSMSClient(accessKeyID, accessKeySecret string, Env string) (*AliyunSMSClient, error) {
	if accessKeyID == "" {
		return nil, errors.New("AccessKeyId should be not empty")
	}
	if accessKeySecret == "" {
		return nil, errors.New("SecretAccessKey should be not empty")
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-qingdao", accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &AliyunSMSClient{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		Client:          client,
		Env:             Env,
	}, nil
}

// SendSms 发送短信
func (client *AliyunSMSClient) SendSms(phoneNumber, signName string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	// 添加环境
	templateParamJSON := client.convertTemplateParam(templateParam, templateAction, 0)
	// 转化成阿里云短信模版
	templateCode := simpleChineseTemplateAction[templateAction]
	// 转化成阿里云所需的手机号
	phoneNumber = joinPhoneNumber(phoneNumber, DefaultNationCode)

	sortQueryStringTmp, errGetSortQueryString := client.getSortQueryStringTmp(phoneNumber, templateCode, templateParamJSON, signName)

	if errGetSortQueryString != nil {
		return false, errGetSortQueryString
	}
	// 去除第一个多余的 & 符号
	sortedQueryString := sortQueryStringTmp[1:]
	// * HTTPMethod + “&” + specialUrlEncode(“/”) + ”&” + specialUrlEncode(sortedQueryString)
	stringToSign := fmt.Sprintf("GET&%s&%s", specialURLEncode("/"), specialURLEncode(sortedQueryString))

	sign, errGetSign := getSign(client.AccessKeySecret, stringToSign)
	if errGetSign != nil {
		return false, errGetSign
	}
	// 签名最后也要做特殊 URL 编码
	signature := specialURLEncode(sign)
	return sendDySmsAPI(signature, sortQueryStringTmp)
}

// convertTemplateParam 转化成阿里的所需的模板参数（number 是群发短信的参数数量）
func (client *AliyunSMSClient) convertTemplateParam(templateParam map[string]string, templateAction TemplateAction, number int) string {
	switch templateAction {
	case TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_FAIL:
		params := &TemplateTidParam{
			TId: templateParam["code"],
			Env: client.Env,
		}
		jsonTemplateParam, _ := json.Marshal(params)
		return string(jsonTemplateParam)
	case TemplateAction_TEMPLATE_ACTION_TENANT_CREATE_SUCCESS:
		params := &TemplateEnv{
			Env: client.Env,
		}
		jsonTemplateParam, _ := json.Marshal(params)
		return string(jsonTemplateParam)
	case TemplateAction_TEMPLATE_ACTION_REVIEW_NOTIFICATION:
		params := make([]*TemplateReviewNotification, number)
		for i := 0; i < number; i++ {
			params[i] = &TemplateReviewNotification{
				Issue:  templateParam["issue"],
				Number: templateParam["number"],
			}
		}
		jsonTemplateParam, _ := json.Marshal(params)
		return string(jsonTemplateParam)
	default:
		params := &TemplateParam{
			Code: templateParam["code"],
		}
		jsonTemplateParam, _ := json.Marshal(params)
		return string(jsonTemplateParam)
	}
}

// joinPhoneNumber 拼接国际号码
func joinPhoneNumber(phoneNumber, nationCode string) string {
	// phoneNumber 的 0 是需要去掉再上传到发送平台
	if nationCode != "" && strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = phoneNumber[1:]
	}
	if nationCode == "+86" || nationCode == "86" {
		return phoneNumber
	}
	return dealNationCode(nationCode) + phoneNumber
}

// getSortQueryStringTmp 得到 sortQueryStringTmp
func (client *AliyunSMSClient) getSortQueryStringTmp(phoneNumber, templateCode, templateParamJSON, signName string) (string, error) {
	params := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   uuid.New().String(),
		"AccessKeyId":      client.AccessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         "cn-hangzhou",
		"PhoneNumbers":     phoneNumber,
		"SignName":         signName,
		"TemplateParam":    templateParamJSON,
		"TemplateCode":     templateCode,
	}
	var keys []string

	for k := range params {
		keys = append(keys, k)
	}
	// 根据参数 Key 排序
	sort.Strings(keys)

	var sortQueryStringTmp string

	for _, value := range keys {
		// specialUrlEncode（参数 Key）+ "=" + specialUrlEncode（参数值）
		sortQueryStringTmp = fmt.Sprintf("%s&%s=%s", sortQueryStringTmp, specialURLEncode(value), specialURLEncode(params[value]))
	}
	return sortQueryStringTmp, nil
}

// sendDySmsAPI 发送阿里云短信 API
func sendDySmsAPI(signature, sortQueryStringTmp string) (bool, error) {
	// http://dysmsapi.aliyuncs.com/?Signature=" + signature + sortQueryStringTmp
	str := fmt.Sprintf("%s/?Signature=%s%s", DysmsAPIEndpoint, signature, sortQueryStringTmp)
	resp, err := http.Get(str)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	ssr := &SendSmsReply{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return false, err
	}
	if ssr.Code != "OK" {
		return false, errors.New(ssr.Code)
	}
	return true, nil
}

// specialURLEncode 处理特殊 URL 编码
func specialURLEncode(in string) string {
	// URLEncode后 加号（+）替换成 %20、星号（*）替换成 %2A、%7E 替换回波浪号（~）
	return strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~").Replace(url.QueryEscape(in))
}

// getSign 签名采用 HmacSHA1 算法 + Base64
func getSign(accessSecret, stringToSign string) (string, error) {
	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessSecret)))
	_, err := mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func dealNationCode(nationCode string) string {
	if nationCode == "" {
		return "86"
	}
	if nationCode != "" && strings.HasPrefix(nationCode, "+") {
		return nationCode[1:]
	}
	return nationCode
}

// SendNotification 发送通知
func (c *AliyunSMSClient) SendNotification(phoneNumber []string, signName string, templateAction TemplateAction, templateParam map[string]string) (bool, error) {
	req := dysmsapi.CreateSendBatchSmsRequest()
	sn := make([]string, len(phoneNumber))
	for k := range phoneNumber {
		sn[k] = signName
	}
	snj, err := json.Marshal(sn)
	if err != nil {
		return false, err
	}
	purePhoneJson, err := json.Marshal(phoneNumber)
	if err != nil {
		return false, err
	}
	// 构建参数
	req.TemplateCode = simpleChineseTemplateAction[templateAction]
	req.SignNameJson = string(snj)
	req.PhoneNumberJson = string(purePhoneJson)
	req.TemplateParamJson = c.convertTemplateParam(templateParam, templateAction, len(phoneNumber))
	// 批量发送
	_, err = c.Client.SendBatchSms(req)
	if err != nil {
		return false, err
	}
	return true, nil
}
