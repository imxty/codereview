package user

import (
	"context"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

// 获取FAQ
func (u *UserAPIHandler) ListFAQEntry(ctx context.Context, req *pb.ListFAQEntryRequest, rsp *pb.ListFAQEntryResponse) error {
	// 获取faq
	faqs, err := u.userStore.ListFaqs(ctx)
	if err != nil {
		return errors.Errorf(codes.DataAccessFailed, "list faq failed[%s]", err.Error())
	}

	// 构建返回数据
	faqEntry := make([]*pb.FAQEntry, len(faqs))
	for k, v := range faqs {
		faqEntry[k] = &pb.FAQEntry{
			// 问题
			Query: v.GetQuestion(),
			// 回答
			Answer: v.GetAnswer(),
		}
	}

	rsp.Faqs = faqEntry
	return nil
}
