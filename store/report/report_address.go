package report

import (
	"time"

	"github.com/jinmukeji/huimaibao-service/svc/report/domain"
)

type ReportAddress struct {
	ReportID    string    `gorm:"primary_key;column:report_id"`
	Mac         string    `gorm:"column:mac"`
	CountryCode int       `gorm:"column:country_code"`
	Country     string    `gorm:"column:country"`
	Province    string    `gorm:"column:province"`
	City        string    `gorm:"column:city"`
	District    string    `gorm:"column:district"`
	FullAddress string    `gorm:"column:full_address"`
	Ip          string    `gorm:"column:ip"`
	Latitude    string    `gorm:"column:latitude"`
	Longitude   string    `gorm:"column:longitude"`
	Remark      string    `gorm:"column:remark"`
	CreatedAt   time.Time // 创建时间
}

func (r ReportAddress) TableName() string {
	return "report_address"
}

// domain->db
func (r *ReportAddress) FromDomainReportAddress(d domain.ReportAddressIntf) {
	if r == nil || d == nil {
		return
	}

	r.ReportID = d.GetReportID()
	r.Mac = d.GetMac()
	r.CountryCode = d.GetCountryCode()
	r.Country = d.GetCountry()
	r.Province = d.GetProvince()
	r.City = d.GetCity()
	r.District = d.GetDistrict()
	r.FullAddress = d.GetFullAddress()
	r.Ip = d.GetIp()
	r.Latitude = d.GetLatitude()
	r.Longitude = d.GetLongitude()
	r.Remark = d.GetRemark()
}

// db->domain
func (r *ReportAddress) ToDomainReportAddress() domain.ReportAddressIntf {
	if r == nil {
		return nil
	}

	p := domain.ReportAddress{
		ReportID:    r.ReportID,
		Mac:         r.Mac,
		CountryCode: r.CountryCode,
		Country:     r.Country,
		Province:    r.Province,
		City:        r.City,
		District:    r.District,
		FullAddress: r.FullAddress,
		Ip:          r.Ip,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Remark:      r.Remark,
	}
	return &p
}
