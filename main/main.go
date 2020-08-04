package main

import (
	"github.com/gin-gonic/gin"
	"restapi_server/common"
	"restapi_server/pkg/config"
	"restapi_server/pkg/service"
)

// 程序入口
func main() {
	run()
}

func run() {
	gin.SetMode(gin.DebugMode)
	e := gin.Default()
	r := service.NewUserRouter()
	r.Route(e)

	e.Run(config.RESTPort)
}

// 仿照 tantan-backend-common 项目的 service 启动
func serviceRun() {
	s := service.NewService()
	common.RunService(s).Wait()
}
