package service

import (
	"github.com/gin-gonic/gin"
	"log"
	_ "net/http"
	"restapi_server/common"
)

func NewService() *UserService {
	router := NewUserRouter()
	service := &UserService{httpServer: common.NewHttpServer(common.HttpConfig{Router: router, Addr: ":8080"})}
	return service
}

// 实现 Service 接口
type UserService struct {
	httpServer common.HttpServer
}

func (s UserService) Start() (err error) {
	return s.httpServer.Start()
}

func (s UserService) Stop() (err error) {
	err = s.httpServer.Stop()
	if err != nil {
		log.Fatalf("%v", err)
	}
	return err
}

// 请求参数 Bind
type State struct {
	State string `json:"state" binding:"required"`
}

type User struct {
	Name string `json:"name" binding:"required"`
}

// 自定义 HandlerFunc
type HandlerFuncWithError func(c *gin.Context) *APIException

// 错误处理统一收敛到 wrapper 装饰器，wrapper 返回 gin.HandlerFunc
func wrapper(handler HandlerFuncWithError) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			c.JSON(err.Code, err)
			return
		}
	}
}

type UserRouter struct {
	handler HttpHandler
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		handler: HttpHandler{},
	}
}

func (r *UserRouter) Route(e *gin.Engine) {
	group := e.Group("/users")

	// 获取用户信息
	group.GET("/", wrapper(r.handler.GetUsers))
	// 创建用户信息
	group.POST("/", wrapper(r.handler.CreateUser))
	// 获取用户关系
	group.GET("/:user_id/relationships", wrapper(r.handler.GetRelationship))
	// 更新用户关系
	group.PUT("/:user_id/relationships/:other_user_id", wrapper(r.handler.UpdateRelationship))
}
