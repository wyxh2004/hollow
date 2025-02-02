package configs

import "hollow/pkg/env"

var (
	// MongoDB配置
	MongoURI = env.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	DBName   = env.GetEnv("DB_NAME", "hollow")

	// JWT配置
	JWTSecret = env.GetEnv("JWT_SECRET", "your-secret-key-please-change-in-production")

	// SMTP配置
	SMTPHost     = env.GetEnv("SMTP_HOST", "smtp.163.com")
	SMTPPort     = env.GetEnvAsInt("SMTP_PORT", 587)
	SMTPUsername = env.GetEnv("SMTP_USERNAME", "wyssixsixsix@163.com")
	SMTPPassword = env.GetEnv("SMTP_PASSWORD", "your-app-password")

	// 服务器配置
	ServerPort = env.GetEnv("SERVER_PORT", "8080")
)
