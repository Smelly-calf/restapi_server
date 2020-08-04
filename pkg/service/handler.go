package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"restapi_server/pkg/model"
	"strconv"
	"time"
)

const (
	USER         = "user"
	RELATIONSHIP = "relationship"
)

type HttpHandler struct{}

type UserSchema struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type",default:"user"`
}

type RelationSchema struct {
	FollowerID int    `json:"user_id"`
	State      string `json:"state"`
	Type       string `json:"type",default:"relationship"`
}

func (h *HttpHandler) GetUsers(c *gin.Context) *APIException {
	flag := c.Query("flag")
	log.Printf("time:%v, flag:%v", time.Now(), flag)
	if flag == "1" {
		log.Println("sleep 10 second.")
		time.Sleep(10 * time.Second)
	}
	ums, err := model.SelectAllUsers()
	if err != nil {
		return ParameterError("查询发生错误")
	}

	// 转换为 userSchema
	users := make([]UserSchema, 0)
	for _, u := range ums {
		users = append(users, UserSchema{
			ID:   u.ID,
			Name: u.Name,
			Type: USER,
		})
	}
	c.JSON(200, users)
	return nil
}

func (h *HttpHandler) CreateUser(c *gin.Context) *APIException {
	var u User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		return ParameterError("参数 name 非法")
	}
	user, err := model.InsertUser(u.Name)
	if err != nil {
		return ParameterError("参数 name 非法")
	}

	c.JSON(200, UserSchema{
		ID:   user.ID,
		Name: user.Name,
		Type: USER,
	})
	return nil
}

func (h *HttpHandler) GetRelationship(c *gin.Context) *APIException {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return ParameterError("参数 userID 非法")
	}
	models, err := model.DefaultRelationModel.GetRelationshipsByUserID(int(userID))
	if err != nil {
		return ParameterError("参数 userID 非法")
	}

	relations := make([]RelationSchema, 0)
	for _, m := range models {
		relations = append(relations, RelationSchema{
			FollowerID: m.FollowerID,
			State:      m.State,
			Type:       RELATIONSHIP,
		})
	}
	c.JSON(200, relations)
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
	relation, err := model.DefaultRelationModel.UpdateRelation(int(userID), int(followerID), s.State)
	if err != nil {
		return ParameterError("参数非法")
	}

	c.JSON(200, RelationSchema{
		FollowerID: relation.FollowerID,
		State:      relation.State,
		Type:       RELATIONSHIP,
	})
	return nil
}
