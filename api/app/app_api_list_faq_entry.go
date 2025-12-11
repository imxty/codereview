package app

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/api/app/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/api"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
)

func (s *AppAPIHandler) ListFAQEntry(ctx context.Context, req *pb.ListFAQEntryRequest, rsp *pb.ListFAQEntryResponse) error {
	// 发送请求
	listFAQEntryRsp, err := s.userAPI.ListFAQEntry(ctx, &userpb.ListFAQEntryRequest{})
	// 如果错误不为空
	if err != nil {
		// 返回错误，先获取code，然后返回对应的msg
		return errors.Error(api.GetSrvErrorCode(err), api.ErrorMsg(err))
	}

	// 返回数据
	appFaqs := make([]*pb.FAQEntry, len(listFAQEntryRsp.GetFaqs()))
	for k, v := range listFAQEntryRsp.GetFaqs() {
		appFaqs[k] = &pb.FAQEntry{
			Query:  v.GetQuery(),
			Answer: v.GetAnswer(),
		}
	}

	rsp.Faqs = appFaqs

	return nil
}
