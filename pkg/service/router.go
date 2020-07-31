package service

import (
	"github.com/gin-gonic/gin"
	_ "net/http"
	"restapi_server/pkg/model"
	"strconv"
)

type State struct {
	State string `json:"state" binding:"required"` // todo 验证 state 枚举值范围
}

type User struct {
	Name string `json:"name" binding:"required"`
}

type HandlerFuncWithError func(c *gin.Context) *APIException

// 错误处理统一收敛到 wrapper 装饰器
func wrapper(handler HandlerFuncWithError) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			c.JSON(err.Code, err)
			return
		}
	}
}

// 服务和路由
func Router() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 用户组
	userR := r.Group("/users")

	// 用户信息
	userR.GET("/", wrapper(func(c *gin.Context) *APIException {
		users := model.SelectAllUsers()
		c.JSON(200, users)
		return nil
	}))
	userR.POST("/", wrapper(func(c *gin.Context) *APIException {
		var u User
		err := c.ShouldBindJSON(&u)
		if err != nil {
			return ParameterError("参数 name 非法")
		}
		user := model.InsertUser(u.Name)
		c.JSON(200, user)
		return nil
	}))

	// 用户关系
	userR.GET("/:user_id/relationships", wrapper(func(c *gin.Context) *APIException {
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			return ParameterError("参数 userID 非法")
		}
		res := model.GetRelationshipsByUserID(int(userID))
		c.JSON(200, res)
		return nil
	}))

	userR.PUT("/:user_id/relationships/:other_user_id", wrapper(func(c *gin.Context) *APIException {
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			return ParameterError("参数 userID 非法")
		}
		followerID, err := strconv.ParseInt(c.Param("other_user_id"), 10, 64)
		if err != nil {
			return ParameterError("参数 other_user_id 非法")
		}

		var s State
		err = c.ShouldBindJSON(&s)
		if err != nil {
			return ParameterError("参数 state 非法")
		}
		relation := model.UpdateRelation(int(userID), int(followerID), s.State)
		c.JSON(200, relation)
		return nil
	}))

	return r
}
