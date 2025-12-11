import https from '../utils/https';

function raw(method, request) {
    const d = {
        method: method,
        service: "com.jinmuhealth.huimai.api.organization",
        request: request
    }
    return d
}

//注册组织管理
export function SignUpRequest(params) {
    return https.post('', raw("UserAPI.SignUp", params))
}

//账号登录组织管理
export function SignInByUsernameRequest(params) {
    return https.post('', raw("UserAPI.SignInByUsername", params))
}

//刷新access_token
export function RefreshAccessTokenRequest(params) {
    return https.post('', raw("UserAPI.RefreshAccessToken", params))
}

//手机号登录组织管理
export function SignInByPhoneRequest(params) {
    return https.post('', raw("UserAPI.SignInByPhone", params))
}

//发送手机短信
export function SendPhoneVerificationCode(params) {
    return https.post('', raw("NotificationAPI.SendPhoneVerificationCode", params))
}

//验证手机验证码
export function VerifyPhoneVerificationCodeRequest(params) {
    return https.post('', raw("NotificationAPI.VerifyPhoneVerificationCode", params))
}

//重置密码
export function ResetOrganizationPasswordRequest(params) {
    return https.post('', raw("UserAPI.ResetOrganizationPassword", params))
}

//判断组织用户名是否存在
export function CheckOrganizationUsernameExistRequest(params) {
    return https.post('', raw("UserAPI.CheckOrganizationUsernameExist", params))
}

//判断组织手机号是否存在
export function CheckOrganizationPhoneExistRequest(params) {
    return https.post('', raw("UserAPI.CheckOrganizationPhoneExist", params))
}

//创建商户
export function CreateTenantRequest(params) {
    return https.post('', raw("UserAPI.CreateTenant", params))
}

//获取组织概况数据
export function GetOrganizationTenantSummaryRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationTenantSummary", params))
}

//获取组织测量数据
export function GetOrganizationMeasurementSummaryRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationMeasurementSummary", params))
}

//获取商户排行请求
export function ListTopTenantRankRequest(params) {
    return https.post('', raw("UserAPI.ListTopTenantRank", params))
}

//获取组织报告统计相应请求
export function SearchReportsRequestRequest(params) {
    return https.post('', raw("ReportAPI.SearchReports", params))
}

//获取组织vip概况请求
export function GetOrganizationCustomerCompareSummaryRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationCustomerCompareSummary", params))
}

//获取组织测量概况请求
export function GetOrganizationTenantMeasurementSummaryRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationTenantMeasurementSummary", params))
}

//查询商户数据对比请求
export function SearchOrganizationTenantCompareRequest(params) {
    return https.post('', raw("UserAPI.SearchOrganizationTenantCompare", params))
}

//获取商户排行请求
export function SearchCustomerLastStatusRequest(params) {
    return https.post('', raw("CustomerAPI.SearchCustomerLastStatus", params))
}

//获取商户列表请求
export function ListTenantsRequest(params) {
    return https.post('', raw("UserAPI.ListTenants", params))
}

//获取商户模版请求
export function GetTenantStencilRequest(params) {
    return https.post('', raw("UserAPI.GetTenantStencil", params))
}

//激活商户请求
export function ActivateTenantRequest(params) {
    return https.post('', raw("UserAPI.ActivateTenant", params))
}

//停用商户请求
export function DisableTenantRequest(params) {
    return https.post('', raw("UserAPI.DisableTenant", params))
}

//搜索商户请求
export function SearchTenantsRequest(params) {
    return https.post('', raw("UserAPI.SearchTenants", params))
}

//创建商户模版请求
export function CreateTenantStencilResponse(params) {
    return https.post('', raw("UserAPI.CreateTenantStencil", params))
}

//批量激活商户请求
export function BatchActivateTenantsRequest(params) {
    return https.post('', raw("UserAPI.BatchActivateTenants", params))
}

//获取商户链接
export function GetTenantInvitationLinkRequest(params) {
    return https.post('', raw("UserAPI.GetTenantInvitationLink", params))
}

//获取商品列表请求
export function ListProductsRequest(params) {
    return https.post('', raw("ProductAPI.ListProducts", params))
}

//获取症候名称和key的映射关系请求
export function GetSymptomKeyMapRequest(params) {
    return https.post('', raw("ProductAPI.GetSymptomKeyMap", params))
}

//添加商品请求
export function CreateProductRequest(params) {
    return https.post('', raw("ProductAPI.CreateProduct", params))
}

//更新商品请求
export function UpdateProductRequest(params) {
    return https.post('', raw("ProductAPI.UpdateProduct", params))
}

//删除商品请求
export function DeleteProductRequest(params) {
    return https.post('', raw("ProductAPI.DeleteProduct", params))
}

//获取商品推荐方案列表请求
export function ListTreatmentsRequest(params) {
    return https.post('', raw("ProductAPI.ListTreatments", params))
}

//创建商品推荐方案请求
export function CreateTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.CreateTreatment", params))
}

//获取症候商品请求
export function ListSymptomProductsRequest(params) {
    return https.post('', raw("ProductAPI.ListSymptomProducts", params))
}

//获取商品方案审核状态请求
export function GetTreatmentReviewStatusRequest(params) {
    return https.post('', raw("ProductAPI.GetTreatmentReviewStatus", params))
}

//提交商品推荐方案请求
export function SubmitTreatmentToReviewRequest(params) {
    return https.post('', raw("ProductAPI.SubmitTreatmentToReview", params))
}

//获取商品推荐方案请求
export function GetTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.GetTreatment", params))
}

//删除商品推荐方案请求
export function DeleteTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.DeleteTreatment", params))
}

//更新商品推荐方案请求
export function UpdateTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.UpdateTreatment", params))
}

//取消商品推荐方案审核请求
export function CancelTreatmentReviewRequest(params) {
    return https.post('', raw("ProductAPI.CancelTreatmentReview", params))
}

//发布商品推荐方案请求
export function PublishTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.PublishTreatment", params))
}

//SearchBillingRequest
export function SearchBillingRequest(params) {
    return https.post('', raw("BillingAPI.SearchBilling", params))
}

//获取商品推荐状态
export function GetOrganizationProductRecommendationStatusRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationProductRecommendationStatus", params))
}

//修改商品推荐状态请求
export function ModifyProductRecommendationRequest(params) {
    return https.post('', raw("UserAPI.ModifyProductRecommendation", params))
}

//删除商户请求
export function DeleteTenantRequest(params) {
    return https.post('', raw("UserAPI.DeleteTenant", params))
}

//获取预付款费用请求
export function GetPrestorePriceRequest(params) {
    return https.post('', raw("BillingAPI.GetPrestorePrice", params))
}

//提交商户预付款审核请求
export function SubmitTenantPrestoreAccountRequest(params) {
    return https.post('', raw("BillingAPI.SubmitTenantPrestoreAccount", params))
}

//为商户预付款付费请求
export function PayTenantPrestoreAccountRequest(params) {
    return https.post('', raw("BillingAPI.PayTenantPrestoreAccount", params))
}

//提交实体资质认证请求
export function CommitEntityCertificateRequest(params) {
    return https.post('', raw("UserAPI.CommitEntityCertificate", params))
}

//修改组织密码请求
export function UpdateOrganizationPasswordRequest(params) {
    return https.post('', raw("UserAPI.UpdateOrganizationPassword", params))
}

//修改组织手机号请求
export function UpdateOrganizationPhoneRequest(params) {
    return https.post('', raw("UserAPI.UpdateOrganizationPhone", params))
}


//检测订单支付状态请求
export function CheckBillingPayStatusRequest(params) {
    return https.post('', raw("BillingAPI.CheckBillingPayStatus", params))
}

//获取商品统计数据
export function GetProductStatisticsRequest(params) {
    return https.post('', raw("ProductAPI.GetProductStatistics", params))
}

//获取商户请求
export function GetTenantRequest(params) {
    return https.post('', raw("UserAPI.GetTenant", params))
}

//获取审核通知列表
export function ListReviewNotificationsRequest(params) {
    return https.post('', raw("NotificationAPI.ListReviewNotifications", params))
}
//确认审核通知
export function ConfirmReviewNotificationRequest(params) {
    return https.post('', raw("NotificationAPI.ConfirmReviewNotification", params))
}

//查看账单请求
export function GetBillingRequest(params) {
    return https.post('', raw("BillingAPI.GetBilling", params))
}

//支付账单请求
export function PayBillingRequest(params) {
    return https.post('', raw("BillingAPI.PayBilling", params))
}

//获取首次认证失败的商户信息请求
export function GetUnAuthTenantRevisionRequest(params) {
    return https.post('', raw("UserAPI.GetUnAuthTenantRevision", params))
}

//获取首次认证失败的商户信息请求
export function ReCreateTenantRequest(params) {
    return https.post('', raw("UserAPI.ReCreateTenant", params))
}

//修改商户模版
export function ModifyTenantStencilRequest(params) {
    return https.post('', raw("UserAPI.ModifyTenantStencil", params))
}

//获取组织详情
export function GetOrganizationDetailRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationDetail", params))
}

//查询组织是否欠费
export function CheckOrganizationArrearsRequest(params) {
    return https.post('', raw("BillingAPI.CheckOrganizationArrears", params))
}

//通过模板ID获取商户模板
export function GetTenantStencilByIDRequest(params) {
    return https.post('', raw("UserAPI.GetTenantStencilByID", params))
}

//UploadImage 上传图片
export function UploadImageRequest(params) {
    return https.post('', raw("UserAPI.UploadImage", params))
}

//创建抬头
export function CreateInvoiceTitleRequest(params) {
    return https.post('', raw("BillingAPI.CreateInvoiceTitle", params))
}

//获取组织发票抬头
export function ListOrganizationInvoiceTitlesRequest(params) {
    return https.post('', raw("BillingAPI.ListOrganizationInvoiceTitles", params))
}

//修改发票抬头
export function ModifyInvoiceTitleRequest(params) {
    return https.post('', raw("BillingAPI.ModifyInvoiceTitle", params))
}

//删除发票抬头
export function DeleteInvoiceTitleRequest(params) {
    return https.post('', raw("BillingAPI.DeleteInvoiceTitle", params))
}

//获取组织未开票订单请求
export function ListOrganizationUnInvoicedBillingsRequest(params) {
    return https.post('', raw("BillingAPI.ListOrganizationUnInvoicedBillings", params))
}

//创建发票请求
export function CreateInvoiceRequest(params) {
    return https.post('', raw("BillingAPI.CreateInvoice", params))
}

//获取发票请求
export function GetInvoiceRequest(params) {
    return https.post('', raw("BillingAPI.GetInvoice", params))
}

//获取发票请求
export function ListOrganizationInvoicesRequest(params) {
    return https.post('', raw("BillingAPI.ListOrganizationInvoices", params))
}

//批量获取组织商户每月报告统计请求
export function ListOrganizationTenantMonthlyReportCountRequest(params) {
    return https.post('', raw("ReportAPI.ListOrganizationTenantMonthlyReportCount", params))
}

//获取组织每月报告统计请求
export function GetOrganizationMonthlyReportCountRequest(params) {
    return https.post('', raw("ReportAPI.GetOrganizationMonthlyReportCount", params))
}

//分配方案
export function SubmitTreatmentToTenantRequest(params) {
    return https.post('', raw("UserAPI.SubmitTreatmentToTenant", params))
}

//分页查询商户续费明细响应
export function SearchTenantSubscriptionsPaginationRequest(params) {
    return https.post('', raw("UserAPI.SearchTenantSubscriptionsPagination", params))
}

//下架方案
export function RemoveTreatmentRequest(params) {
    return https.post('', raw("ProductAPI.RemoveTreatment", params))
}


//GetOrganizationSubscriptionsMonthly
export function GetOrganizationSubscriptionsMonthlyRequest(params) {
    return https.post('', raw("UserAPI.GetOrganizationSubscriptionsMonthly", params))
}

//SearchTenantSubscriptionsRequest
export function SearchTenantSubscriptionsRequestRequest(params) {
    return https.post('', raw("UserAPI.SearchTenantSubscriptions", params))
}

//CancelTreatmentToTenants
export function CancelTreatmentToTenantsRequest(params) {
    return https.post('', raw("UserAPI.CancelTreatmentToTenants", params))
}

//更改组织联系人手机号
export function UpdateOrganizationContactPhoneRequest(params) {
    return https.post('', raw("UserAPI.UpdateOrganizationContactPhone", params))
}

export function SearchOrganizationTenantDownloadRequest(params) {
    return https.post('', raw("UserAPI.SearchOrganizationTenantDownload", params))
}
