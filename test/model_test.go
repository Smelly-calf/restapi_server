package test

import (
	"restapi_server/pkg/model"
	"testing"
)

func TestUser(t *testing.T) {
	t.Logf("查询user by name：%+v\n", model.SelectUserByName("mike"))
	t.Logf("创建user记录：%+v\n", model.InsertUser("max"))

	t.Logf("查询某条关系：%+v\n", model.DefaultRelationModel.GetOneRelation(1, 11))
	t.Logf("创建一个关系：%+v\n", model.DefaultRelationModel.InsertRelationship(11, 1, "disliked"))
	t.Logf("更改某个关系：%+v\n", model.DefaultRelationModel.UpdateRelation(11, 1, "liked"))
}
