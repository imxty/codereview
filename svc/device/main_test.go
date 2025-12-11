package device

import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/jinmukeji/huimaibao-service/pkg/dbutils"
)

const (
	defaultFedSource      = "localhost:3306"
	defaultSourceUser     = "root"
	defaultSourcePassword = "p@ssw0rd"
	defaultSourceDBName   = "huimaibao"
)

var (
	flagFedSource      = flag.String("mysql.source.address", defaultEnv("DB_SOURCE_ADDR", defaultFedSource), "mysql fed source address")
	flagSourceUser     = flag.String("mysql.source.user", defaultEnv("DB_SOURCE_USER", defaultSourceUser), "mysql fed source user")
	flagSourcePassword = flag.String("mysql.source.password", defaultEnv("DB_SOURCE_PASSWORD", defaultSourcePassword), "mysql fed source password")
	flagSourceDBName   = flag.String("mysql.source.dbname", defaultEnv("DB_SOURCE_NAME", defaultSourceDBName), "mysql fed source db name")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func defaultEnv(env, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}

	return v
}

func GetConn() *dbutils.DBConnection {
	s := make(map[string]string)
	s["parseTime"] = "true"
	v, _ := json.Marshal(s)
	conn := &dbutils.DBConnection{
		Username: *flagSourceUser,
		Password: *flagSourcePassword,
		Net:      "tcp",
		Address:  *flagFedSource,
		Schema:   *flagSourceDBName,
		Params:   string(v),
	}
	return conn
}
