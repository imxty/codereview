package marshaler

import (
	"strconv"
	"testing"

	"github.com/docker/distribution/context"
	"github.com/jinmukeji/huimaibao-service/pkg/cache"
	"github.com/jinmukeji/huimaibao-service/pkg/cache/store"
	"github.com/stretchr/testify/suite"
)

type MarshalerTestSuite struct {
	suite.Suite
	m *Marshaler
}

func (suite *MarshalerTestSuite) SetupSuite() {
	rs, _ := store.NewRedisStore(":6379", 0, nil)
	suite.m = New(cache.New(rs))
}

// TestSetString  放入string
func (suite *MarshalerTestSuite) TestSetString() {
	ctx := context.Background()

	err := suite.m.Put(ctx, "a", "abc")
	suite.Assert().NoError(err)
	var res string
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("abc", res)
}

// TestSetInt  放入int
func (suite *MarshalerTestSuite) TestSetInt() {
	ctx := context.Background()

	err := suite.m.Put(ctx, "a", 3)
	suite.Assert().NoError(err)
	var res int
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal(3, res)
}

// TestSetMap  放入map
func (suite *MarshalerTestSuite) TestSetMap() {
	ctx := context.Background()
	m := map[string]string{
		"xxx": "bbb",
	}

	err := suite.m.Put(ctx, "a", m)
	suite.Assert().NoError(err)
	var res map[string]string
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("bbb", res["xxx"])
}

// TestHSetMap  写入map特定字段
func (suite *MarshalerTestSuite) TestHSetMap() {
	ctx := context.Background()
	m := map[string]interface{}{
		"xxx": 22,
	}
	err := suite.m.HSet(ctx, "hset", m)
	suite.Assert().NoError(err)
	res, err := suite.m.HGet(ctx, "hset", "xxx")
	val, _ := res.(string)
	v, _ := strconv.Atoi(val)
	suite.Assert().NoError(err)
	suite.Assert().Equal(22, v)
}

// TestSetArray  放入array
func (suite *MarshalerTestSuite) TestSetArray() {
	ctx := context.Background()
	m := []string{"xxx", "bbb"}

	err := suite.m.Put(ctx, "a", m)
	suite.Assert().NoError(err)
	res := make([]string, 2)
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("xxx", res[0])
}

// TestSetStruct  放入结构体
func (suite *MarshalerTestSuite) TestSetStruct() {
	ctx := context.Background()
	type User struct {
		UserID string `json:"user_id"`
		Age    int32  `json:"age"`
	}
	v := &User{
		UserID: "xxx",
		Age:    64,
	}
	v1 := &User{
		UserID: "xxx1",
		Age:    63,
	}

	vv := []*User{v, v1}
	err := suite.m.Put(ctx, "a", vv)
	suite.Assert().NoError(err)
	res := make([]*User, 2)
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("xxx", res[0].UserID)
	suite.Assert().Equal("xxx1", res[1].UserID)
	suite.Assert().Equal(int32(64), res[0].Age)
	suite.Assert().Equal(int32(63), res[1].Age)
}

// TestSetStructNotPointer  放入结构体非指针类型
func (suite *MarshalerTestSuite) TestSetStructNotPointer() {
	ctx := context.Background()
	type User struct {
		UserID string `json:"user_id"`
		Age    int32  `json:"age"`
	}
	v := User{
		UserID: "xxx",
		Age:    64,
	}

	err := suite.m.Put(ctx, "a", v)
	suite.Assert().NoError(err)
	res := new(User)
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("xxx", res.UserID)
	suite.Assert().Equal(int32(64), res.Age)
}

// TestSetComplicatedStruct  放入复杂结构
func (suite *MarshalerTestSuite) TestSetComplicatedStruct() {
	ctx := context.Background()
	type User struct {
		UserID string            `json:"user_id"`
		Age    int32             `json:"age"`
		Map    map[string]string `json:"map"`
	}
	v := User{
		UserID: "xxx",
		Age:    64,
		Map:    map[string]string{"abc": "ddd"},
	}

	err := suite.m.Put(ctx, "a", v)
	suite.Assert().NoError(err)
	res := new(User)
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("xxx", res.UserID)
	suite.Assert().Equal(int32(64), res.Age)
	suite.Assert().Equal("ddd", res.Map["abc"])
}

// TestSetMapWithMap  放入嵌套Map
func (suite *MarshalerTestSuite) TestSetMapWithMap() {
	ctx := context.Background()
	m := map[string]map[string]string{
		"x": {
			"yy": "dd",
		},
	}

	err := suite.m.Put(ctx, "a", m)
	suite.Assert().NoError(err)
	res := make(map[string]map[string]string)
	err = suite.m.Get(ctx, "a", &res)
	suite.Assert().NoError(err)
	suite.Assert().Equal("dd", res["x"]["yy"])
}

// TestClearByTags  测试通过tags删除key
func (suite *MarshalerTestSuite) TestClearByTags() {
	ctx := context.Background()

	err := suite.m.Put(ctx, "a", "cc", store.WithTags([]string{"test", "test1"}))
	suite.Assert().NoError(err)
	err = suite.m.Put(ctx, "b", "zz", store.WithTags([]string{"test", "test1"}))
	suite.Assert().NoError(err)
	err = suite.m.Put(ctx, "c", "dd", store.WithTags([]string{"test", "test1"}))
	suite.Assert().NoError(err)
	err = suite.m.ClearByTags(ctx, []string{"test", "test1"})
	suite.Assert().NoError(err)
}

// TestExist  检测是否存在接口
func (suite *MarshalerTestSuite) TestExist() {
	ctx := context.Background()

	err := suite.m.Put(ctx, "a", "cc", store.WithTags([]string{"test", "test1"}))
	suite.Assert().NoError(err)

	result, err := suite.m.IsExist(ctx, "a")
	suite.Assert().NoError(err)
	suite.Assert().Equal(true, result)
	result, err = suite.m.IsExist(ctx, "v")
	suite.Assert().NoError(err)
	suite.Assert().Equal(false, result)
}

// TestHIncrBy  map指定字段加一
func (suite *MarshalerTestSuite) TestHIncrBy() {
	ctx := context.Background()
	m := map[string]interface{}{
		"abc": 22,
	}
	err := suite.m.HSet(ctx, "hset", m)
	suite.Assert().NoError(err)
	err = suite.m.HIncrBy(ctx, "hset", "abc", 3)
	suite.Assert().NoError(err)
	res, err := suite.m.HGet(ctx, "hset", "abc")
	val, _ := res.(string)
	v, _ := strconv.Atoi(val)
	suite.Assert().NoError(err)
	suite.Assert().Equal(25, v)
}
func TestMarshalerTestSuite(t *testing.T) {
	suite.Run(t, new(MarshalerTestSuite))
}
