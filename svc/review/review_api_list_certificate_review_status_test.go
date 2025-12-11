package review

import (
	"context"
	"testing"

	pb "github.com/jinmukeji/huimaibao-proto/gen/go/huimaibao/biz/review/v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
	store "github.com/jinmukeji/huimaibao-service/store/review"
	utils "github.com/jinmukeji/huimaibao-service/svc/testing"

	"github.com/stretchr/testify/suite"
)

type ListCertificateReviewStatusTestSuite struct {
	suite.Suite
	hdl *ReviewAPIHandler
}

func (suite *ListCertificateReviewStatusTestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open(dbutils.GetDsn(GetConn())), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	conn := dbutils.NewConnection(db)
	utils.InitData(*flagFedSource, *flagSourceUser, *flagSourcePassword, *flagSourceDBName, reviewFile)
	s := store.NewReviewStore(conn)

	suite.hdl = NewReviewAPIHandler(s, nil, nil, nil)
}

// TestListNotes
func (suite *ListCertificateReviewStatusTestSuite) TestListCertificateReviewStatus() {
	ctx := context.Background()
	req := &pb.ListCertificateReviewStatusRequest{
		TenantIds: []string{tenantId},
	}
	resp := &pb.ListCertificateReviewStatusResponse{}
	err := suite.hdl.ListCertificateReviewStatus(ctx, req, resp)
	suite.Assert().NoError(err)
}

func (suite *ListCertificateReviewStatusTestSuite) TearDownSuite() {}

func TestListCertificateReviewStatusTestSuite(t *testing.T) {
	suite.Run(t, new(ListCertificateReviewStatusTestSuite))
}
