package app_config

import "os"

type appConfig struct {
	DatabaseUrl string
	Mode        string
}

func New() appConfig {
	return appConfig{
		DatabaseUrl: os.Getenv("TODOAPP_DATABASE_URL"),
		Mode:        os.Getenv("TODOAPP_MODE"),
	}
}
