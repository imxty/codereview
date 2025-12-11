package testing

import (
	"math/rand"
	"os"
	"regexp"
	"testing"
	"time"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors"
	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gotest.tools/assert"
)

func ReadFileString(filepath string) (string, error) {
	query, err := ReadFile(filepath)
	data := string(query)
	return data, err
}

func ReadFile(filepath string) ([]byte, error) {
	query, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return query, err
}

func GetErrCode(err error) string {
	code := regexp.MustCompile(`\[errcode:(\d+)\]`)
	errcode := code.FindStringSubmatch(err.Error())
	return errcode[1]
}

// MustRpcError 解析 err 为标准 RpcError. 解析失败返回 nil.
func MustRpcError(err error) *errors.RpcError {
	if v, ok := err.(*errors.RpcError); ok {
		return v
	}
	return nil
}

// AssertRpcErrorCode 断言指定 error 是否为一个 RpcError 且 Code 符合预期
// 如果 err 不是一个 RpcError，则会将 codes.Unknown 与 expected 进行比较.
func AssertRpcErrorCode(t *testing.T, expected codes.Code, err error) {
	c := errors.Code(err)
	assert.Equal(t, expected, c)
}

// randomDuration 返回 [min, max] 之间大小的毫秒精度的 time.Duration
func RandomDuration(min, max int) time.Duration {
	rand.NewSource((time.Now().UnixNano()))
	ms := rand.Intn(max-min+1) + min
	return time.Duration(ms) * time.Millisecond
}

// rpcLatency 随机返回 100~500 毫秒之间的 RPC 调用时延
func RpcLatency() time.Duration {
	return RandomDuration(100, 500)
}

func OpenMySqlDB(addr, user, password, database string) (*gorm.DB, error) {
	cfg := mysqlCfg.NewConfig()
	cfg.MultiStatements = true
	cfg.Addr = addr
	cfg.User = user
	cfg.Net = "tcp"
	cfg.Passwd = password
	cfg.DBName = database
	return gorm.Open(mysql.Open(cfg.FormatDSN()), nil)
}

func InitData(source, user, password, database, file string) {
	db, _ := OpenMySqlDB(source, user, password, database)
	f, err := ReadFileString(file)
	if err != nil {
		panic(err)
	}
	err = db.Exec(f).Error
	if err != nil {
		panic(err)
	}
}
