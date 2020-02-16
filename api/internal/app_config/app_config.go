package app_config

import "os"

type appConfig struct {
	Mode string
}

func New() appConfig {
	return appConfig{
		Mode: os.Getenv("TODOAPP_MODE"),
	}
}
