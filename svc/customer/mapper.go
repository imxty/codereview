package customer

import (
	gerr "errors"
	"time"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/customer/v1"
	"github.com/jinmukeji/huimaibao-service/svc/customer/domain"
)

// proto->domain Gender
func toDomainGender(gender pb.Gender) (int32, error) {
	switch gender {
	case pb.Gender_GENDER_INVALID:
		return domain.ErrTipInvalidGender, gerr.New("invalid gender")
	case pb.Gender_GENDER_UNSET:
		return domain.UnsetGender, gerr.New("unset gender")
	case pb.Gender_GENDER_MALE:
		return domain.Male, nil
	case pb.Gender_GENDER_FEMALE:
		return domain.Female, nil
	}
	return domain.ErrTipInvalidGender, gerr.New("invalid gender")
}

// domain->proto Customer
func toProtoCustomer(d domain.CustomerIntf, initial string, age int32) *pb.Customer {
	if d == nil {
		return nil
	}

	b := d.GetBirthday()
	birthday := &pb.Date{
		Year:  int32(b.Year()),
		Month: int32(b.Month()),
		Day:   int32(b.Day()),
	}

	return &pb.Customer{
		// 常客 ID
		CustomerId: d.GetCustomerID(),
		// 员工 ID
		StaffId: d.GetStaffID(),
		// 常客昵称
		Nickname: d.GetNickname(),
		// 常客首字母
		Initial: initial,
		// 常客性别
		Gender: toProtoGender(d.GetGender()),
		// 常客手机号
		Phone: d.GetPhone(),
		// 常客生日
		Birthday: birthday,
		// 常客年龄
		Age: age,
		// 常客身高
		Height: d.GetHeight(),
		// 常客体重
		Weight: d.GetWeight(),
		// 常客既往病史
		Pmh: d.GetPmh(),
		// 常客其他备注
		Remarks: d.GetRemark(),
	}
}

// domain->proto Gender
func toProtoGender(gender int32) pb.Gender {
	switch gender {
	case domain.UnsetGender:
		return pb.Gender_GENDER_UNSET
	case domain.Male:
		return pb.Gender_GENDER_MALE
	case domain.Female:
		return pb.Gender_GENDER_FEMALE
	}
	return pb.Gender_GENDER_INVALID
}

// proto->domain
func toDomainBirthday(birthday *pb.Date) time.Time {
	return mapProtoDateToTime(birthday)
}

// mapProtoDateToTime
func mapProtoDateToTime(pd *pb.Date) time.Time {
	if pd == nil {
		return time.Now()
	}
	return time.Date(int(pd.GetYear()), time.Month(pd.GetMonth()), int(pd.GetDay()), 0, 0, 0, 0, time.UTC)
}

// toProtoCustomerLastStatus
func toProtoCustomerLastStatus(d domain.CustomerLastStatusIntf) *pb.CustomerLastStatus {
	if d == nil {
		return nil
	}

	updateTime := d.GetUpdatedAt()
	updateTime = time.Date(updateTime.Year(), updateTime.Month(), updateTime.Day(), 0, 0, 0, 0, time.Local)
	nowTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	return &pb.CustomerLastStatus{
		OverdueCount: int32(nowTime.Sub(updateTime).Hours() / 24),
	}
}
