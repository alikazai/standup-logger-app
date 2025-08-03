package utils

import "os"

func GetHTTPListenAddress() string {
	return os.Getenv("HTTP_LISTEN_ADDRESS")
}

//	func GetJwtSecret() string {
//		return os.Getenv("JWT_SECRET")
//	}
//
//	func GetMongoDbName() string {
//		return os.Getenv("MONGO_DB_NAME")
//	}
//
//	func GetMongoDbUrl() string {
//		return os.Getenv("MONGO_DB_URL")
//	}
func GetHostOrigin() string {
	return os.Getenv("HOST_ORIGIN")
}

func GetEnvironment() string {
	return os.Getenv("ENVIRONMENT")
}
