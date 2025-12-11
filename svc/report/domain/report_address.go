package domain

type ReportAddress struct {
	ReportID    string
	Mac         string
	CountryCode int
	Country     string
	Province    string
	City        string
	District    string
	FullAddress string
	Ip          string
	Latitude    string
	Longitude   string
	Remark      string
}

type ReportAddressMapper interface {
	ToDomainReportAddressMapper
	FromDomainReportAddressMapper
}

type ToDomainReportAddressMapper interface {
	ToDomainReportAddress() ReportAddressIntf
}

type FromDomainReportAddressMapper interface {
	FromDomainReportAddress(ReportAddressIntf)
}

type ReportAddressIntf interface {
	GetReportID() string
	GetMac() string
	GetCountryCode() int
	GetCountry() string
	GetProvince() string
	GetCity() string
	GetDistrict() string
	GetFullAddress() string
	GetIp() string
	GetLatitude() string
	GetLongitude() string
	GetRemark() string
}

var _ ReportAddressIntf = (*ReportAddress)(nil)

func (r *ReportAddress) GetReportID() string {
	if r == nil {
		return ""
	}
	return r.ReportID
}

func (r *ReportAddress) GetMac() string {
	if r == nil {
		return ""
	}
	return r.Mac
}

func (r *ReportAddress) GetCountryCode() int {
	if r == nil {
		return 0
	}
	return r.CountryCode
}

func (r *ReportAddress) GetCountry() string {
	if r == nil {
		return ""
	}
	return r.Country
}

func (r *ReportAddress) GetProvince() string {
	if r == nil {
		return ""
	}
	return r.Province
}

func (r *ReportAddress) GetCity() string {
	if r == nil {
		return ""
	}
	return r.City
}

func (r *ReportAddress) GetDistrict() string {
	if r == nil {
		return ""
	}
	return r.District
}

func (r *ReportAddress) GetFullAddress() string {
	if r == nil {
		return ""
	}
	return r.FullAddress
}

func (r *ReportAddress) GetIp() string {
	if r == nil {
		return ""
	}
	return r.Ip
}

func (r *ReportAddress) GetLatitude() string {
	if r == nil {
		return ""
	}
	return r.Latitude
}

func (r *ReportAddress) GetLongitude() string {
	if r == nil {
		return ""
	}
	return r.Longitude
}

func (r *ReportAddress) GetRemark() string {
	if r == nil {
		return ""
	}
	return r.Remark
}
