package config

import "time"

var (
	// JWT 配置
	JWTSecret     = []byte("my-first-jwt") // 生产环境请使用环境变量
	JWTExpiration = time.Hour              // Token 有效期 1h

)
