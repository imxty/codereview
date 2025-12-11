package dbutils

import (
	"encoding/json"
	"time"

	mysqlCfg "github.com/go-sql-driver/mysql"
)

// DBConnection 数据库连接信息
type DBConnection struct {
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	Net       string `json:"net" yaml:"net"`
	Address   string `json:"address" yaml:"address"`
	Schema    string `json:"schema" yaml:"schema"`
	Collation string `json:"collation" yaml:"collation"`
	Timeout   int32  `json:"timeout" yaml:"timeout"`
	PubKey    string `json:"pub_key" yaml:"pub_key"`
	Params    string `json:"params" yaml:"params"`
}

// GetDsn DBConnection转dsn
func GetDsn(conn *DBConnection) string {
	// 必要
	cfg := mysqlCfg.NewConfig()
	cfg.User = conn.Username
	cfg.Passwd = conn.Password
	cfg.Net = conn.Net
	cfg.Addr = conn.Address
	cfg.DBName = conn.Schema
	cfg.ParseTime = true
	// 可选
	if conn.Collation != "" {
		cfg.Collation = conn.Collation
	}
	if conn.Params != "" {
		var params map[string]string
		err := json.Unmarshal([]byte(conn.Params), &params)
		if err != nil {
			panic(err)
		}
		cfg.Params = params
	}
	if conn.Timeout != 0 {
		cfg.Timeout = time.Duration(conn.Timeout) * time.Second
	}
	if conn.PubKey != "" {
		cfg.ServerPubKey = conn.PubKey
	}
	return cfg.FormatDSN()
}
