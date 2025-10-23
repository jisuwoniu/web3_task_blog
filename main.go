package main

import (
	"log"
	"web3_task_blog/internal/routes"
	"web3_task_blog/internal/utils"
)

func main() {
	// 初始化配置
	utils.InitConfig()

	// /HTTP
	r := routes.SetupRoutes()
	log.Println("Server is running on :8080")
	log.Fatal(r.Run(":8080"))
}