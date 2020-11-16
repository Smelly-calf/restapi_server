package test

import (
	"bytes"
	"encoding/json"
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
		e := gin.Default()
		r := service.NewUserRouter()
		r.Route(e)
		server := httptest.NewServer(e)
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
	fmt.Printf("%+v\n", resp)
	fmt.Printf("%+v\n", body)
	resp.Status(http.StatusOK)
	inverse, _ := model.DefaultRelationModel.GetOneRelation(11, 1)
	if inverse.State != "liked" {
		t.Errorf("inverse state: %v", inverse.State)
	}
}

// BenchmarkUpdateRelationshipAPI-8   720  1488407 ns/op
// 接口耗时大约 1.5ms，1s可处理请求量约 666, 即 TPS 约 666 .
func BenchmarkUpdateRelationshipAPI(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	client := &http.Client{}
	rUrl := "http://localhost:8080/users/1/relationships/11"
	data, _ := json.Marshal(service.State{
		State: "disliked",
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest(http.MethodPut, rUrl, bytes.NewBuffer(data))
		request.Header.Set("Content-Type", "application/json")
		client.Do(request)
	}
}

// BenchmarkUpdateRelationshipAPI-8   720  1488407 ns/op
// 接口耗时大约 1.5ms，QPS 720.
func BenchmarkGetRelation(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	client := &http.Client{}
	rUrl := "http://localhost:8080/users/1/relationships"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest(http.MethodGet, rUrl, nil)
		client.Do(request)
	}
}
