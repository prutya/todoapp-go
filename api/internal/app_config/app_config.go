package app_config

import "os"

type appConfig struct {
	AuthSecret        string
	AuthExpirySeconds int64
	BcryptCost        int
	DatabaseUrl       string
	Mode              string
}

func New() appConfig {
	return appConfig{
		AuthSecret:        os.Getenv("TODOAPP_AUTH_SECRET"),
		AuthExpirySeconds: 10 * 24 * 60 * 60,
		BcryptCost:        14,
		DatabaseUrl:       os.Getenv("TODOAPP_DATABASE_URL"),
		Mode:              os.Getenv("TODOAPP_MODE"),
	}
}
