package h5

import (
	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/h5/v1"
	productpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/product/v1"
	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
)

type H5APIHandler struct {
	reportAPI  reportpb.ReportAPIService
	productAPI productpb.ProductAPIService
	userAPI    userpb.UserAPIService
	s3Domain   string
}

var _ pb.ReportAPIHandler = (*H5APIHandler)(nil)
var _ pb.ProductAPIHandler = (*H5APIHandler)(nil)

func NewH5APIHandler(reportAPI reportpb.ReportAPIService,
	productAPI productpb.ProductAPIService,
	userAPI userpb.UserAPIService,
	s3Domain string) *H5APIHandler {
	return &H5APIHandler{
		reportAPI:  reportAPI,
		productAPI: productAPI,
		userAPI:    userAPI,
		s3Domain:   s3Domain,
	}
}
