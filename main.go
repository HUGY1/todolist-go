// todolist 主程序入口
// 基于 Gin 和 GORM 的 Todo List 后端 API 服务
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"todolist/config"
	"todolist/models"
	"todolist/routes"
	"todolist/utils"
)

func main() {
	utils.InitJWT()
	// 初始化日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("========== Todo List API 服务启动 ==========")

	// 初始化数据库连接
	config.InitDB()

	// 自动迁移：根据 Todo 模型创建或更新 list 表结构
	// AutoMigrate 会创建缺失的表和列，不会删除已存在的列
	config.GetDB().AutoMigrate(&models.File{}, &models.Todo{})
	if err := config.GetDB().AutoMigrate(&models.Todo{}, &models.User{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库表 list,users 已就绪")

	// 配置路由
	router := routes.SetupRouter()

	// fakeData := mock.CreateFakeData()
	// for _, todo := range fakeData {
	// 	config.GetDB().Create(&todo)
	// }
	// log.Println("数据库表 fakeData 假数据已插入")

	// 在 goroutine 中启动 HTTP 服务，以便主协程可以监听退出信号
	// 默认监听 8080 端口
	go func() {
		log.Println("HTTP 服务监听端口: 3000")
		if err := router.Run(":3000"); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待退出信号（Ctrl+C 或 kill）
	// 优雅关闭：收到信号后关闭数据库连接再退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("收到退出信号，正在关闭服务...")
	config.CloseDB()
	log.Println("服务已安全退出")
}
