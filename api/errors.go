package api

import (
	perr "github.com/jinmukeji/huimaibao-service/pkg/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

var codeToMsg = map[codes.Code]string{
	ErrOK:      "OK",
	ErrUnknown: "Unknown error",

	ErrUserUnauthorized: "User 未授权",

	ErrInvalidRequest:      "Request 数据错误",
	ErrInvalidAccessToken:  "token不合法或者过期，请重新登录",
	ErrInvalidRefreshToken: "不合法的refresh_token",
	ErrKickedOut:           "踢出用户失败",
	ErrInvalidJWT:          "不合法的jwt",
	ErrCreateAccessToken:   "创建access_token失败",
	ErrCreateRefreshToken:  "创建refresh_token失败",
	ErrTokenNotFound:       "Token 不存在",
	// ErrInvalidArgument
	ErrInvalidArgument: "非法的参数",
	// ErrDataAccessFailed
	ErrDataAccessFailed: "数据库访问错误",
	// ErrInvalidOperation 非法的操作
	ErrInvalidOperation: "非法的操作",

	ErrRPCInternal: "RPC请求错误",

	PlatformErrorHeartError:   "测量数据异常,请重新测量",
	PlatformErrorHeartTooLow:  "测量数据异常,请重新测量",
	PlatformErrorHeartTooHigh: "测量数据异常,请重新测量",

	ErrStaffExceedLimit:           "员工数量达到上限",
	ErrStaffPhoneHasBeenUsed:      "员工手机号已被使用",
	ErrPrivilegeGroupNotExist:     "权限组不存在",
	ErrTenantExceedLimit:          "商户数量已达上限",
	ErrStaffHasBeenAdded:          "员工已被其他商户添加",
	ErrStaffHasActivated:          "员工已被激活",
	ErrInvalidTxId:                "不合法的短信凭证",
	ErrTenantNotFound:             "商户不存在",
	ErrStaffNotBelongTenant:       "员工不属于任何商户",
	ErrSystemUserNotFound:         "管理员不存在",
	ErrSysPasswordWrong:           "管理员密码错误",
	ErrOrganizationNotExist:       "组织不存在",
	ErrWrongPassword:              "密码错误",
	ErrOrganizationTenantNotExist: "组织下该商户不存在",
	ErrTenantUnAuth:               "商户未认证",
	ErrWrongOldPassword:           "旧密码错误",
	ErrSamePasswords:              "新旧密码相同",
	ErrStaffHasNotActivated:       "员工被删除",
	ErrPhoneHasBeenUsed:           "手机号已被使用",
	ErrStaffNotExist:              "员工不存在",
	ErrCustomersOutOfUpperLimit:   "常客数量超过上限",
	ErrCustomerPhoneExist:         "常客手机号已存在",
	ErrCustomerNotFound:           "常客不存在",
	ErrSmsNotFound:                "短信不存在",
	ErrSmsRepeatedSending:         "短信重复发送",
	ErrMessageLimitControl:        "短信发送已达上限",
	ErrSmsCodeError:               "短信验证码错误",
	ErrSmsExpired:                 "短信验证码过期",
	ErrReviewPending:              "商户正在审核中",
	ErrStaffPasswordLength:        "密码长度错误",
	ErrTreatmentSubmitExceedLimit: "方案提审超过限制",

	ErrProductNotFound:                  "商品不存在",
	ErrTreatmentNotBelongToOrganization: "方案不属于该组织",

	ErrReportUpdateFailed:      "报告更新异常",
	PlatformErrReportNotEnough: "报告数量不足",
	ErrCustomerAge:             "常客年龄小于16岁",

	ErrGetTreatmentFailed:             "获取方案失败",
	ErrTreatmentExceedLimit:           "方案数量超过上限",
	ErrProductNotBelongToOrganization: "商品不属于该组追",

	ErrReportNotFound:   "报告不存在",
	ErrInitReportFailed: "报告初始化异常",
	ErrTenantDisable:    "商户订阅到期，无法测量",

	ErrReviewIssueNotFound:  "审核不存在",
	ErrReviewTenantNotExist: "审核商户不存在",
	ErrNotificationNotFound: "通知不存在",

	ErrTreatmentUsing:                   "方案使用中，无法下架",
	ErrTenantCanNotActivated:            "商户无法激活",
	ErrTenantActivationTimelineNotFound: "商户激活时间线不存在",
	ErrInvalidTenantPassword:            "商户密码错误",
	ErrStaffDeleted:                     "员工已被删除",
	ErrPrivilegeGroupUsing:              "权限组正在使用中",
	ErrHasNoPrivileges:                  "当前用户没有可访问的权限",
	ErrPlatformError:                    "数据异常,请重新测量",
}

// 错误清单（左右都是闭区间)
// 0-3000 公用错误
const (
	// 错误码定义清单

	// ErrOK OK. Not used.
	ErrOK = 0
	// ErrUnknown Unknown error
	ErrUnknown = 1000

	// ErrUserUnauthorized
	ErrUserUnauthorized = 1001
	// ErrInvalidAccessToken
	ErrInvalidAccessToken = 1002
	// ErrInvalidRefreshToken
	ErrInvalidRefreshToken = 1003
	// ErrInvalidArgument
	ErrInvalidArgument = 1004
	// ErrInvalidRequest
	ErrInvalidRequest = 1006
	// ErrKickedOut
	ErrKickedOut = 1007
	// ErrInvalidJWT
	ErrInvalidJWT = 1008
	// ErrCreateAccessToken
	ErrCreateAccessToken = 1009
	// ErrCreateRefreshToken
	ErrCreateRefreshToken = 1010
	// ErrDataAccessFailed
	ErrDataAccessFailed = 1013
	// ErrInvalidOperation
	ErrInvalidOperation = 1014

	// ErrRPCInternal
	ErrRPCInternal = 3000
	// ErrStaffPasswordLength
	ErrStaffPasswordLength = 3115
	// ErrSignatureError
	ErrSignatureError = 3210

	// ErrTenantCanNotActivated
	ErrTenantCanNotActivated = 5000
	// ErrOrganizationNotExist
	ErrOrganizationNotExist = 5001
	// ErrOrganizationTenantNotExist
	ErrOrganizationTenantNotExist = 5002
	// ErrStaffExceedLimit
	ErrStaffExceedLimit = 5003
	// ErrTenantNotFound
	ErrTenantNotFound = 5004
	// ErrStaffPhoneHasBeenUsed
	ErrStaffPhoneHasBeenUsed = 5005
	// ErrStencilExist
	ErrStencilExist = 5006

	// ErrStencilNotFound
	ErrStencilNotFound = 5007
	// ErrStaffNotExist
	ErrStaffNotExist = 5008
	// ErrTenantActivationTimelineNotFound
	ErrTenantActivationTimelineNotFound = 5009

	// ErrStaffHasNotActivated
	ErrStaffHasNotActivated = 5010
	// ErrPhoneHasBeenUsed
	ErrPhoneHasBeenUsed = 5011

	// ErrWrongOldPassword
	ErrWrongOldPassword = 5012
	// ErrNewPasswordSameWithOld
	ErrSamePasswords = 5013
	// ErrWrongPassword
	ErrWrongPassword = 5014
	// ErrSystemUserNotFound
	ErrSystemUserNotFound = 5015
	// ErrSysPasswordWrong
	ErrSysPasswordWrong = 5016
	// ErrOtpWrong
	ErrOtpWrong = 5017
	// ErrStaffNotBelongTenant
	ErrStaffNotBelongTenant = 5018
	// ErrStaffHasBeenAdded
	ErrStaffHasBeenAdded = 5019
	// ErrStaffHasActivated
	ErrStaffHasActivated = 5020
	// ErrInvalidTxId
	ErrInvalidTxId = 5021

	// ErrInitReportFailed 初始化报告失败
	ErrInitReportFailed = 5022
	// ErrReportNotFound
	ErrReportNotFound = 5023
	// ErrTokenNotFound token不存在
	ErrTokenNotFound = 5024
	// ErrTreatmentSubmitExceedLimit 方案提审超过限制
	ErrTreatmentSubmitExceedLimit = 5025

	// ErrReviewIssueNotFound
	ErrReviewIssueNotFound = 5026
	// ErrGetTreatmentFailed
	ErrGetTreatmentFailed = 5028
	// ErrProductNotBelongToOrganization
	ErrProductNotBelongToOrganization = 5030
	// ErrPublishTreatmentFailed
	ErrPublishTreatmentFailed = 5031
	// ErrSubmitTreatmentTimesNotEnough
	ErrSubmitTreatmentTimesNotEnough = 5032
	// ErrCustomerNotFound 常客不存在
	ErrCustomerNotFound = 5035
	// 常客数量超过上限
	ErrCustomersOutOfUpperLimit = 5036
	// ErrCustomerPhoneExist 常客手机号已存在
	ErrCustomerPhoneExist = 5037
	// ErrNotificationNotFound 通知不存在
	ErrNotificationNotFound = 5038
	// ErrSmsNotFound
	ErrSmsNotFound = 5040
	// ErrSmsRepeatedSending
	ErrSmsRepeatedSending = 5041
	// ErrSmsCodeError
	ErrSmsCodeError = 5042

	// ErrTreatmentUsing
	ErrTreatmentUsing = 5048

	ErrHasNoPrivileges = 5049

	ErrPlatformError = 5051

	ErrInvalidTenantPassword = 5301

	ErrTreatmentExceedLimit = 5302

	ErrReviewTenantNotExist = 5303

	ErrTenantDisable = 5304

	ErrStaffDeleted = 5305

	ErrMessageLimitControl = 5306

	ErrTenantUnAuth = 5307

	ErrSmsExpired = 5308

	ErrTenantExceedLimit = 5309

	ErrCustomerAge   = 5311
	ErrReviewPending = 5312

	ErrPrivilegeGroupUsing = 5347

	ErrPrivilegeGroupNotExist = 5401

	ErrProductNotFound                  = 5029
	ErrTreatmentNotBelongToOrganization = 5050

	ErrReportUpdateFailed      = 3209
	PlatformErrReportNotEnough = 1107

	PlatformErrorHeartError   = 3500
	PlatformErrorHeartTooLow  = 3501
	PlatformErrorHeartTooHigh = 3502
)

// ErrorMsg 根据错误码获得标准错误消息内容
func ErrorMsg(err error) string {
	cd := GetSrvErrorCode(err)
	return codeToMsg[cd]
}

// ErrorChineseMsg 根据错误码获得标准错误消息内容
func ErrorChineseMsg(code codes.Code) string {
	if msg, ok := codeToMsg[code]; ok {
		return msg
	}
	return "数据异常,请稍后再试"
}

// GetSrvErrorCode 获取service返回的code
func GetSrvErrorCode(e error) codes.Code {
	return perr.GetErrCode(e)
}
