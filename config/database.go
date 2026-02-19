// Package config 负责应用配置和数据库连接
package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例，供其他包使用
var DB *gorm.DB

// InitDB 初始化数据库连接
// 使用提供的 MySQL 连接信息建立与数据库的连接
// 先连接 MySQL 服务并自动创建数据库，再连接到指定库
func InitDB() {
	const dbName = "todolist"
	const fileDbName = "filelist"
	// 连接参数（不含数据库名，用于先连接 MySQL 服务）
	// 格式: 用户名:密码@tcp(主机:端口)/?charset=utf8mb4&parseTime=True&loc=Local
	dsnWithoutDB := "root:lgzxp6qg@tcp(test-db-mysql.ns-wzme3ot2.svc:3306)/?charset=utf8mb4&parseTime=True&loc=Local"

	// 第一步：连接 MySQL 服务（不指定数据库）
	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 第二步：执行 CREATE DATABASE IF NOT EXISTS，自动创建数据库
	createSQL := "CREATE DATABASE IF NOT EXISTS `" + dbName + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
	if err := tempDB.Exec(createSQL).Error; err != nil {
		tempDB = nil
		log.Fatalf("创建数据库失败: %v", err)
	}

	log.Printf("数据库 %s 已就绪", dbName)

	// 关闭临时连接
	sqlDB, _ := tempDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// 第三步：连接到指定数据库
	dsn := "root:lgzxp6qg@tcp(test-db-mysql.ns-wzme3ot2.svc:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")
}

// GetDB 获取数据库连接实例
// 供其他包调用以执行数据库操作
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接（用于程序退出时清理资源）
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("获取底层 SQL DB 失败: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已关闭")
		}
	} else {
		log.Println("数据库连接未初始化，无需关闭")
	}
}

// Ping 检查数据库连接是否正常
func Ping() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
