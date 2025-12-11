package main

import (
	"net/http"
	"time"

	"github.com/coze-dev/coze-go"
	mlog "github.com/jinmukeji/go-pkg/v2/log"
	app "github.com/jinmukeji/huimaibao-service/api/ai"
	"github.com/jinmukeji/plat-pkg/v4/rpc"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
)

const (
	// ServiceName 是本微服务的名称
	ServiceName = "ai"
	// ServiceNamespace 是微服务的命名空间
	ServiceNamespace = "com.shangyikangyou.huimaibao.api"
)

const (
	COZE_BASE_URL = "https://api.coze.cn"
)

var (
	log = mlog.StandardLogger()

	// Following values will be set during build.
	// Do NOT manually modify them.

	// ProductVersion is current product version.
	ProductVersion = "(n/a)"
	// GitCommit is the git commit short hash
	GitCommit = "(n/a)"
	// GoVersion is go compiler version `go version`
	GoVersion = "(n/a)"
	// BuildTime is go build time
	BuildTime = "(n/a)"
)

var (
	basePath string
)

type cozeOptions struct {
	PrivateKey          string `json:"private_key" yaml:"private_key"`
	JWTOauthClientID    string `json:"jwt_oauth_client_id" yaml:"jwt_oauth_client_id"`
	JWTOauthPublicKeyID string `json:"jwt_oauth_public_key_id" yaml:"jwt_oauth_public_key_id"`
	BotID               string `json:"bot_id" yaml:"bot_id"`
	SecretKey           string `json:"secret_key" yaml:"secret_key"`
}

func main() {
	opts := service.NewWebOptions(ServiceNamespace, ServiceName)
	opts.ProductVersion = ProductVersion
	opts.GitCommit = GitCommit
	opts.GoVersion = GoVersion
	opts.BuildTime = BuildTime

	svc := service.CreateWeb(opts)

	cozeConfigKey := []string{
		"micro",
		"config",
		"jm",
		"com.shangyikangyou.huimaibao.api.@global",
		"coze",
	}
	cozeOpt := &cozeOptions{}
	err := rpc.YamlConfig().Get(cozeConfigKey...).Scan(cozeOpt)
	if err != nil {
		die(err)
	}

	oauth, err := coze.NewJWTOAuthClient(coze.NewJWTOAuthClientParam{
		ClientID:      cozeOpt.JWTOauthClientID,
		PublicKey:     cozeOpt.JWTOauthPublicKeyID,
		PrivateKeyPEM: cozeOpt.PrivateKey,
	}, coze.WithAuthBaseURL(COZE_BASE_URL))
	if err != nil {
		log.Fatalf("Error creating JWT OAuth client: %v\n", err)
	}

	authCli := coze.NewJWTAuth(oauth, nil)

	httpCli := &http.Client{
		Timeout: 60 * time.Second,
	}

	cozeCli := coze.NewCozeAPI(authCli, coze.WithBaseURL(COZE_BASE_URL), coze.WithHttpClient(httpCli))

	aiApp := app.NewApp(svc, cozeCli, cozeOpt.BotID, cozeOpt.SecretKey)
	svc.Handle("/", aiApp.Handler())

	log.Infof("API Base: %s", basePath)
	aiApp.PrintRoutes()

	// Run Server
	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
