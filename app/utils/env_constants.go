package utils

import (
	"fmt"
	"os"
)

func GetHTTPListenAddress() string {
	addr := os.Getenv("HTTP_LISTEN_ADDRESS")
	if addr == "" {
		if port := os.Getenv("PORT"); port != "" {
			addr = ":" + port
		} else {
			addr = ":8080"
		}
	}
	return addr
}

//	func GetJwtSecret() string {
//		return os.Getenv("JWT_SECRET")
//	}

func GetDatabaseURL() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}
	return dsn
}

func GetHostOrigin() string {
	return os.Getenv("HOST_ORIGIN")
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}
