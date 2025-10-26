package main

import (
	"log"
	"web3_task_blog/internal/repository"
	"web3_task_blog/internal/routes"
	"web3_task_blog/internal/utils"
)

func main() {
	// 初始化配置
	utils.InitConfig()
	repository.InitDB()
	// 自动迁移数据库表
	err := repository.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	// 设置路由
	r := routes.SetupRoutes()
	log.Println("Server is running on :8080")
	log.Fatal(r.Run(":8080"))
}
