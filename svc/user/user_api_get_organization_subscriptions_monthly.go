package user

import (
	"context"
	gerr "errors"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *UserAPIHandler) GetOrganizationSubscriptionsMonthly(ctx context.Context, req *pb.GetOrganizationSubscriptionsMonthlyRequest, rsp *pb.GetOrganizationSubscriptionsMonthlyResponse) error {
	err := validateGetOrganizationSubscriptionsMonthlyRequest(req)
	if err != nil {
		return errors.Error(codes.InvalidRequest, err.Error())
	}

	ts, count, err := u.userStore.GetOrganizationSubscriptionsMonthly(ctx, req.GetOrganizationId(), req.GetStartTime().AsTime(), req.GetEndTime().AsTime(), req.GetPagination().GetSize(), req.GetPagination().GetOffset())
	if err != nil {
		return errors.Error(codes.DataAccessFailed, err.Error())
	}

	results := make([]*pb.TenantSubscriptionMonthly, len(ts))
	for k, v := range ts {
		results[k] = &pb.TenantSubscriptionMonthly{
			Date:        timestamppb.New(v.GetDate()),
			TenantCount: v.GetTenantCount(),
			YearCount:   v.GetYearCount(),
		}
	}

	rsp.Subscriptions = results
	rsp.TotalCount = count

	return nil
}

func validateGetOrganizationSubscriptionsMonthlyRequest(req *pb.GetOrganizationSubscriptionsMonthlyRequest) error {
	if req.GetOrganizationId() == "" {
		return gerr.New("organization_id should not be empty")
	}
	if req.GetStartTime() == nil {
		return gerr.New("start_time should not be nil")
	}
	if req.GetEndTime() == nil {
		return gerr.New("end_time should not be nil")
	}
	if req.GetPagination() == nil {
		return gerr.New("pagination should not be nil")
	}
	return nil
}
