package domain

type CertificateReviewNote struct {
	CertificateReviewNoteID string
	Content                 string
	Rev                     int32
}

type CertificateReviewNoteMapper interface {
	ToDomainCertificateReviewNoteMapper
	FromDomainCertificateReviewNoteMapper
}

type ToDomainCertificateReviewNoteMapper interface {
	ToDomainCertificateReviewNote() CertificateReviewNoteIntf
}

type FromDomainCertificateReviewNoteMapper interface {
	FromDomainCertificateReviewNote(CertificateReviewNoteIntf)
}

type CertificateReviewNoteIntf interface {
	GetCertificateReviewNoteID() string
	GetContent() string
	GetRev() int32
}

var _ CertificateReviewNoteIntf = (*CertificateReviewNote)(nil)

func (c *CertificateReviewNote) GetCertificateReviewNoteID() string {
	if c == nil {
		return ""
	}
	return c.CertificateReviewNoteID
}

func (c *CertificateReviewNote) GetContent() string {
	if c == nil {
		return ""
	}
	return c.Content
}

func (c *CertificateReviewNote) GetRev() int32 {
	if c == nil {
		return 0
	}
	return c.Rev
}
