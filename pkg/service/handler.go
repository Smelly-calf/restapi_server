package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"restapi_server/pkg/model"
	"strconv"
	"time"
)

type HttpHandler struct{}

func (h *HttpHandler) GetUsers(c *gin.Context) *APIException {
	flag := c.Query("flag")
	log.Printf("time:%v, flag:%v", time.Now(), flag)
	if flag == "1" {
		log.Println("sleep 10 second.")
		time.Sleep(10 * time.Second)
	}
	users := model.SelectAllUsers()
	c.JSON(200, users)
	return nil
}

func (h *HttpHandler) CreateUser(c *gin.Context) *APIException {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		return ParameterError("参数 name 非法")
	}
	user := model.InsertUser(u.Name)
	c.JSON(200, user)
	return nil
}

func (h *HttpHandler) GetRelationship(c *gin.Context) *APIException {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return ParameterError("参数 userID 非法")
	}
	res := model.DefaultRelationModel.GetRelationshipsByUserID(int(userID))
	c.JSON(200, res)
	return nil
}

func (h *HttpHandler) UpdateRelationship(c *gin.Context) *APIException {
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
	relation := model.DefaultRelationModel.UpdateRelation(int(userID), int(followerID), s.State)
	c.JSON(200, relation)
	return nil
}
