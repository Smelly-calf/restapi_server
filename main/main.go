package main

import (
	"restapi_server/pkg/service"
)

// 程序入口
func main() {
	r := service.Router()
	r.Run(":8080")
}
