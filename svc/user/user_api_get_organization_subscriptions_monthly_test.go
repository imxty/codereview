package user

import (
	"context"
	"testing"
	"time"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	store "github.com/jinmukeji/huimaibao-service/store/user"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	"github.com/stretchr/testify/suite"
)

type GetOrganizationSubscriptionsMonthlyTestSuite struct {
	suite.Suite
	hdl *UserAPIHandler
}

func (suite *GetOrganizationSubscriptionsMonthlyTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, userFile)
	s := store.NewUserStore(conn)
	suite.hdl = NewUserAPIHandler(s, nil, nil, nil, nil, nil, nil)
}

func (suite *GetOrganizationSubscriptionsMonthlyTestSuite) TestGetOrganizationSubscriptionsMonthly() {
	ctx := context.Background()
	now := time.Now()
	startTime := timestamppb.New(now)
	stopTime := timestamppb.New(now.Add(time.Second * time.Duration(200)))

	req := &userpb.GetOrganizationSubscriptionsMonthlyRequest{
		OrganizationId: organizationId,
		StartTime:      startTime,
		EndTime:        stopTime,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.GetOrganizationSubscriptionsMonthlyResponse{}
	err := suite.hdl.GetOrganizationSubscriptionsMonthly(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *GetOrganizationSubscriptionsMonthlyTestSuite) TestOrganizationIdIsNull() {
	ctx := context.Background()
	t := suite.T()
	req := &userpb.GetOrganizationSubscriptionsMonthlyRequest{
		OrganizationId: organizationIdIsNull,
		StartTime:      nil,
		EndTime:        nil,
		Pagination: &userpb.Pagination{
			Offset: offset,
			Size:   size,
		},
	}
	resp := &userpb.GetOrganizationSubscriptionsMonthlyResponse{}
	err := suite.hdl.GetOrganizationSubscriptionsMonthly(ctx, req, resp)
	utils.AssertRpcErrorCode(t, codes.InvalidRequest, err)
}

func (suite *GetOrganizationSubscriptionsMonthlyTestSuite) TearDownSuite() {

}

func TestGetOrganizationSubscriptionsMonthlyTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrganizationSubscriptionsMonthlyTestSuite))
}
