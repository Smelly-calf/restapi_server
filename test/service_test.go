package test

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"restapi_server/pkg/model"
	"restapi_server/pkg/service"
	"testing"
)

var eng *httpexpect.Expect

func GetEngine(t *testing.T) *httpexpect.Expect {
	gin.SetMode(gin.TestMode)
	if eng == nil {
		// 加载路由到服务器
		server := httptest.NewServer(service.Router())
		eng = httpexpect.New(t, server.URL)
	}
	return eng
}

func TestRelationshipFields(t *testing.T) {
	e := GetEngine(t)
	resp := e.GET("/users/11/relationships").Expect()
	resp.Status(http.StatusOK)
	body := resp.JSON()
	body.Array().First().Object().Value("state").Equal("matched")
}

func TestUpdateRelationshipToDisliked(t *testing.T) {
	e := GetEngine(t)
	relation := map[string]interface{}{
		"state": "disliked",
	}
	resp := e.PUT("/users/1/relationships/11").WithJSON(relation).Expect()
	body := resp.JSON()
	fmt.Printf("%+v\n",resp)
	fmt.Printf("%+v\n",body)
	resp.Status(http.StatusOK)
	inverse := model.GetOneRelation(11, 1)
	if inverse.State != "liked" {
		t.Errorf("inverse state: %v", inverse.State)
	}
}
