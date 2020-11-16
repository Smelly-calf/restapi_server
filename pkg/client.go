package main

// 使用 go-client 发送请求
import (
	"context"
	"fmt"
	swagger "restapi_server/go-client"
)

var cli = swagger.NewAPIClient(swagger.NewConfiguration())

func AddPet() {
	p := swagger.Pet{
		Id:        1,
		Category:  &swagger.Category{Id: 1, Name: "cat"},
		Name:      "jojo",
		PhotoUrls: []string{"a"},
		Status:    "healthy",
	}
	cli.PetApi.AddPet(context.Background(), p)
}

func ListPet() { // localhost:8080/v2/pet/findByStatus
	pets, resp, err := cli.PetApi.FindPetsByStatus(context.Background(), []string{"healthy"})
	fmt.Println(pets)
	fmt.Println(resp)
	fmt.Println(err)
}

func main() {
	AddPet()
}
