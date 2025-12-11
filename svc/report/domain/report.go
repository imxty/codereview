package domain

import "time"

const (
	// HandLeft 左手
	HandLeft = 2
	// HandRight 右手
	HandRight = 3
)

const (
	TongueFaceReportStatusNotGet  = 0
	TongueFaceReportStatusSuccess = 1
	TongueFaceReportStatusFailed  = 2
)

const (
	// GenderUnset 未设置
	GenderUnset = 1
	// GenderMale 男性
	GenderMale = 2
	// GenderFemale 女性
	GenderFemale = 3
)

const (
	CustomerTypeBoth     = 0
	CustomerTypeCustomer = 1
	CustomerTypeTemp     = 2
)

type Report struct {
	ReportID               string
	TenantID               string
	OrganizationID         string
	StaffID                string
	DeviceID               string
	IsCustomer             bool
	CustomerID             string
	Age                    string
	Gender                 int32
	Hand                   int32
	Remarks                string
	StaffRemarks           string
	HeartRate              int32
	Spo                    int32
	Risks                  string
	C0                     int32
	C1                     int32
	C2                     int32
	C3                     int32
	C4                     int32
	C5                     int32
	C6                     int32
	C7                     int32
	PhysicalDialectics     string
	DirtyDialectic         string
	DirtyDialecticFlag     int32
	DietaryAdvice          string
	MeasurementJudgment    string
	F0                     int32
	F1                     int32
	F2                     int32
	F3                     int32
	F4                     int32
	IsStressState          bool
	IsReportComplete       bool
	StressState            string
	TenantTreatmentRevID   string
	FaceImageURL           string
	TongueImageURL         string
	TongueFaceReportStatus int32
	TongueFaceReport       string
	Rev                    int32
	CreatedAt              time.Time
}

type ReportMapper interface {
	ToDomainReportMapper
	FromDomainReportMapper
}

type ToDomainReportMapper interface {
	ToDomainReport() ReportIntf
}

type FromDomainReportMapper interface {
	FromDomainReport(ReportIntf)
}

type ReportIntf interface {
	GetReportID() string
	GetTenantID() string
	GetOrganizationID() string
	GetStaffID() string
	GetDeviceID() string
	GetIsCustomer() bool
	GetCustomerID() string
	GetAge() string
	GetGender() int32
	GetHand() int32
	GetRemarks() string
	GetStaffRemarks() string
	GetHeartRate() int32
	GetSpo() int32
	GetRisks() string
	GetC0() int32
	GetC1() int32
	GetC2() int32
	GetC3() int32
	GetC4() int32
	GetC5() int32
	GetC6() int32
	GetC7() int32
	GetPhysicalDialectics() string
	GetDirtyDialectic() string
	GetDirtyDialecticFlag() int32
	GetDietaryAdvice() string
	GetMeasurementJudgment() string
	GetF0() int32
	GetF1() int32
	GetF2() int32
	GetF3() int32
	GetF4() int32
	GetIsStressState() bool
	GetIsReportComplete() bool
	GetStressState() string
	GetTenantTreatmentRevID() string
	GetFaceImageURL() string
	GetTongueImageURL() string
	GetTongueFaceReportStatus() int32
	GetTongueFaceReport() string
	GetRev() int32
	GetCreatedAt() time.Time
}

var _ ReportIntf = (*Report)(nil)

func (r *Report) GetReportID() string {
	if r == nil {
		return ""
	}
	return r.ReportID
}

func (r *Report) GetTenantID() string {
	if r == nil {
		return ""
	}
	return r.TenantID
}

func (r *Report) GetOrganizationID() string {
	if r == nil {
		return ""
	}
	return r.OrganizationID
}

func (r *Report) GetStaffID() string {
	if r == nil {
		return ""
	}
	return r.StaffID
}

func (r *Report) GetDeviceID() string {
	if r == nil {
		return ""
	}
	return r.DeviceID
}

func (r *Report) GetIsCustomer() bool {
	if r == nil {
		return false
	}
	return r.IsCustomer
}

func (r *Report) GetCustomerID() string {
	if r == nil {
		return ""
	}
	return r.CustomerID
}

func (r *Report) GetAge() string {
	if r == nil {
		return ""
	}
	return r.Age
}

func (r *Report) GetGender() int32 {
	if r == nil {
		return 0
	}
	return r.Gender
}

func (r *Report) GetHand() int32 {
	if r == nil {
		return 0
	}
	return r.Hand
}

func (r *Report) GetRemarks() string {
	if r == nil {
		return ""
	}
	return r.Remarks
}

func (r *Report) GetStaffRemarks() string {
	if r == nil {
		return ""
	}
	return r.StaffRemarks
}

func (r *Report) GetHeartRate() int32 {
	if r == nil {
		return 0
	}
	return r.HeartRate
}

func (r *Report) GetSpo() int32 {
	if r == nil {
		return 0
	}
	return r.Spo
}

func (r *Report) GetRisks() string {
	if r == nil {
		return ""
	}
	return r.Risks
}

func (r *Report) GetC0() int32 {
	if r == nil {
		return 0
	}
	return r.C0
}

func (r *Report) GetC1() int32 {
	if r == nil {
		return 0
	}
	return r.C1
}

func (r *Report) GetC2() int32 {
	if r == nil {
		return 0
	}
	return r.C2
}

func (r *Report) GetC3() int32 {
	if r == nil {
		return 0
	}
	return r.C3
}

func (r *Report) GetC4() int32 {
	if r == nil {
		return 0
	}
	return r.C4
}

func (r *Report) GetC5() int32 {
	if r == nil {
		return 0
	}
	return r.C5
}

func (r *Report) GetC6() int32 {
	if r == nil {
		return 0
	}
	return r.C6
}

func (r *Report) GetC7() int32 {
	if r == nil {
		return 0
	}
	return r.C7
}

func (r *Report) GetPhysicalDialectics() string {
	if r == nil {
		return ""
	}
	return r.PhysicalDialectics
}

func (r *Report) GetDirtyDialectic() string {
	if r == nil {
		return ""
	}
	return r.DirtyDialectic
}

func (r *Report) GetDirtyDialecticFlag() int32 {
	if r == nil {
		return 0
	}
	return r.DirtyDialecticFlag
}

func (r *Report) GetDietaryAdvice() string {
	if r == nil {
		return ""
	}
	return r.DietaryAdvice
}

func (r *Report) GetMeasurementJudgment() string {
	if r == nil {
		return ""
	}
	return r.MeasurementJudgment
}

func (r *Report) GetF0() int32 {
	if r == nil {
		return 0
	}
	return r.F0
}

func (r *Report) GetF1() int32 {
	if r == nil {
		return 0
	}
	return r.F1
}

func (r *Report) GetF2() int32 {
	if r == nil {
		return 0
	}
	return r.F2
}

func (r *Report) GetF3() int32 {
	if r == nil {
		return 0
	}
	return r.F3
}

func (r *Report) GetF4() int32 {
	if r == nil {
		return 0
	}
	return r.F4
}

func (r *Report) GetIsStressState() bool {
	if r == nil {
		return false
	}
	return r.IsStressState
}

func (r *Report) GetIsReportComplete() bool {
	if r == nil {
		return false
	}
	return r.IsReportComplete
}

func (r *Report) GetStressState() string {
	if r == nil {
		return ""
	}
	return r.StressState
}

func (r *Report) GetTenantTreatmentRevID() string {
	if r == nil {
		return ""
	}
	return r.TenantTreatmentRevID
}

func (r *Report) GetFaceImageURL() string {
	if r == nil {
		return ""
	}
	return r.FaceImageURL
}

func (r *Report) GetTongueImageURL() string {
	if r == nil {
		return ""
	}
	return r.TongueImageURL
}

func (r *Report) GetTongueFaceReportStatus() int32 {
	if r == nil {
		return 0
	}
	return r.TongueFaceReportStatus
}

func (r *Report) GetTongueFaceReport() string {
	if r == nil {
		return ""
	}
	return r.TongueFaceReport
}

func (r *Report) GetRev() int32 {
	if r == nil {
		return 0
	}
	return r.Rev
}

func (r *Report) GetCreatedAt() time.Time {
	if r == nil {
		return time.Time{}
	}
	return r.CreatedAt
}
