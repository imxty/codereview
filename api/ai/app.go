package ai

import (
	"net/http"

	"github.com/coze-dev/coze-go"
	"github.com/gin-gonic/gin"
	"github.com/jinmukeji/plat-pkg/v4/rpc/service"
	pweb "github.com/jinmukeji/plat-pkg/v4/web"
	wlog "github.com/jinmukeji/plat-pkg/v4/web/middleware/logger"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/web"
)

type App struct {
	router    *gin.Engine
	api       coze.CozeAPI
	botID     string
	secretKey string
}

func NewApp(svc web.Service, api coze.CozeAPI, botID, secretKey string) *App {
	r := gin.New()
	r.Use(wlog.Logger(service.Logger()))
	r.Use(gin.Recovery())

	// 设置重试次数
	err := svc.Options().Service.Client().Init(
		client.Retries(0),
	)
	die(err)

	app := App{
		router:    r,
		api:       api,
		botID:     botID,
		secretKey: secretKey,
	}
	err = registerHandlers(&app)
	die(err)

	return &app
}

func registerHandlers(app *App) error {
	r := app.router

	// Route handlers
	r.GET("/_health", app.hdlHealth)

	// AI对话
	r.GET("/ai_chat", app.hdlAIChat)

	return nil
}

func (a *App) Handler() http.Handler {
	if a == nil {
		return nil
	}

	return a.router
}

func (a *App) PrintRoutes() {
	pweb.PrintRoutes(a.router)
}

func die(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
