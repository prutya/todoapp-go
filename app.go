package main

import (
	"fmt"

	"github.com/prutya/todoapp-go/internal/app_config"
)

func main() {
	config := app_config.New()

	fmt.Println(config.Mode)
}
