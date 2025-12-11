package customer

import (
	"context"
	gerr "errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/jinmukeji/go-pkg/v2/age"
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	userv1 "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"github.com/mozillazg/go-pinyin"
	"github.com/rs/xid"
)

const (
	// 常客数量超过上限
	ErrCustomersOutOfUpperLimit = 5036
	// ErrCustomerPhoneExist 常客手机号已存在
	ErrCustomerPhoneExist = 5037
)

const (
	// 有效手机号码长度
	validPhoneLength = 11
)

const (
	// 常客数量上限
	customerLimit = 5000
	// ErrCustomerAge
	ErrCustomerAge = 5311
)

func (c *CustomerAPIHandler) AddCustomer(ctx context.Context, req *pb.AddCustomerRequest, rsp *pb.AddCustomerResponse) error {
	err := validateAddCustomerRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	birth := req.GetCustomer().GetBirthday()
	cusAge := getAge(birth.GetYear(), birth.GetMonth(), birth.GetDay())
	if cusAge < 16 {
		return errors.Errorf(ErrCustomerAge, "invalid age[%d]", cusAge)
	}

	// 获取常客总数
	count, err := c.customerStore.GetCustomerCount(ctx, req.GetTenantId())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	// 常客数量超过上限报错
	if count >= customerLimit {
		return errors.Errorf(ErrCustomersOutOfUpperLimit, "customer out of limit tenant:[%s]", req.GetTenantId())
	}
	// 查询添加常客的手机号 已存在返回错误
	phoneExist, err := c.customerStore.CheckCustomerPhoneExist(ctx, req.GetTenantId(), req.GetCustomer().GetPhone(), "+86")
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}
	if phoneExist {
		return errors.Errorf(ErrCustomerPhoneExist, "customer phone exist[%s]", req.GetCustomer().GetPhone())
	}

	// 获取租户组织ID
	getRsp, err := c.userAPI.GetTenantEntity(ctx, &userv1.GetTenantEntityRequest{
		TenantId: req.GetTenantId(),
	})
	if err != nil {
		return err
	}

	// customer 转换 proto -> domain
	dGender, err := toDomainGender(req.GetCustomer().GetGender())
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}
	dBirthday := mapProtoDateToTime(req.GetCustomer().GetBirthday())
	dCustomer := &domain.Customer{
		CustomerID:     xid.New().String(),
		OrganizationID: getRsp.GetEntity().GetOrganizationId(),
		TenantID:       req.GetTenantId(),
		StaffID:        req.GetCustomer().GetStaffId(),
		Nickname:       req.GetCustomer().GetNickname(),
		Initial:        getNicknameInitial(req.GetCustomer().GetNickname()),
		Gender:         dGender,
		Phone:          req.GetCustomer().GetPhone(),
		Birthday:       dBirthday,
		Height:         req.GetCustomer().GetHeight(),
		Weight:         req.GetCustomer().GetWeight(),
		Pmh:            req.GetCustomer().GetPmh(),
		Remark:         req.GetCustomer().GetRemarks(),
	}

	// 在数据库获创建常客
	acErr := c.customerStore.AddCustomer(ctx, req.GetTenantId(), dCustomer)
	if acErr != nil {
		return errors.Errorf(codes.DataAccessFailed, "add customer failed[%s]", acErr.Error())
	}

	// 返回结果
	appCustomer := toProtoCustomer(dCustomer, getNicknameInitial(dCustomer.GetNickname()), int32(age.Age(dCustomer.GetBirthday())))
	rsp.Customer = appCustomer
	return nil
}

// 验证request
func validateAddCustomerRequest(req *pb.AddCustomerRequest) error {
	if req.GetTenantId() == "" {
		return gerr.New("tenant_id should not be empty")
	}
	if req.GetCustomer().GetStaffId() == "" {
		return gerr.New("staff_id should not be empty")
	}
	if req.GetCustomer().GetNickname() == "" || utf8.RuneCountInString(req.GetCustomer().GetNickname()) > 50 {
		return gerr.New("nickname should not be empty")
	}
	if req.GetCustomer().GetPhone() == "" || len(req.GetCustomer().GetPhone()) != validPhoneLength {
		return gerr.New("invalid phone number")
	}
	if req.GetCustomer().GetBirthday() == nil {
		return gerr.New("birthday should not be nil")
	}
	if _, err := toDomainGender(req.GetCustomer().GetGender()); err != nil {
		return err
	}
	if req.GetCustomer().GetHeight() != 0 && (req.GetCustomer().GetHeight() < 50 || req.GetCustomer().GetHeight() > 250) {
		return gerr.New("invalid height")
	}
	if req.GetCustomer().GetWeight() != 0 && (req.GetCustomer().GetWeight() < 30 || req.GetCustomer().GetWeight() > 250) {
		return gerr.New("invalid weight")
	}
	if utf8.RuneCountInString(req.GetCustomer().GetPmh()) > 100 {
		return gerr.New("invalid pmh")
	}
	if utf8.RuneCountInString(req.GetCustomer().GetRemarks()) > 100 {
		return gerr.New("invalid remarks")
	}
	return nil
}

// getAge
func getAge(year, month, day int32) int {
	return age.Age(time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC))
}

// getNicknameInitial 获得昵称的首字母
func getNicknameInitial(nickname string) string {
	nicknameRune := []rune(nickname)
	if len(nicknameRune) == 0 {
		return ""
	}
	nicknamePrefix := nicknameRune[0]
	// 是否为汉字
	if unicode.Is(unicode.Han, nicknamePrefix) {
		a := pinyin.NewArgs()
		a.Style = pinyin.FirstLetter
		nicknameInitial := pinyin.Pinyin(string(nicknamePrefix), a)
		if len(nicknameInitial) == 0 || len(nicknameInitial[0]) == 0 {
			return ""
		}
		return strings.ToUpper(nicknameInitial[0][0])
	}
	// 是否为英文字母
	if unicode.IsLetter(nicknamePrefix) && nicknamePrefix < unicode.MaxASCII {
		return strings.ToUpper(string(nicknamePrefix))
	}
	// 其他字符
	return "#"
}
