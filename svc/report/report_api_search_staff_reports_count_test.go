package report

import (
	"context"
	"testing"

	reportpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/report/v1"
	userpb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/user/v1"
	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/report"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"
	usermock "github.com/jinmukeji/huimaibao-service/svc/user/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SearchStaffReportsCountTestSuite struct {
	suite.Suite
	hdl  *ReportAPIHandler
	user *usermock.UserAPIService
}

func (suite *SearchStaffReportsCountTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reportFile)
	s := store.NewReportStore(conn)
	user := &usermock.UserAPIService{}
	suite.user = user
	suite.hdl = NewReportAPIHandler(s, nil, nil, nil, user, nil, "", "")
}

// TestSearchStaffReportsCount
func (suite *SearchStaffReportsCountTestSuite) TestSearchStaffReportsCount() {
	ctx := context.Background()
	suite.user.On("SearchStaffsByNameAndStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(&userpb.SearchStaffsByNameAndStatusResponse{
			Staffs: []*userpb.Staff{
				&userpb.Staff{
					StaffId: staffId,
				},
			},
		}, nil).After(utils.RpcLatency())
	req := &reportpb.SearchStaffReportsCountRequest{
		StartTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: seconds,
			Nanos:   nanos,
		},
		TenantId:  tenantId,
		StaffName: nickname,
		IsActivated: &wrapperspb.BoolValue{
			Value: true,
		},
	}
	resp := &reportpb.SearchStaffReportsCountResponse{}
	err := suite.hdl.SearchStaffReportsCount(ctx, req, resp)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(resp)
}

func (suite *SearchStaffReportsCountTestSuite) TearDownSuite() {
}

func TestSearchStaffReportsCountTestSuite(t *testing.T) {
	suite.Run(t, new(SearchStaffReportsCountTestSuite))
}
