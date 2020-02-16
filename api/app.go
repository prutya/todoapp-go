package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/prutya/todoapp-go/internal/app_config"
)

func main() {
	config := app_config.New()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	db, err := sql.Open("postgres", config.DatabaseUrl)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	sqlStatement := `
		SELECT login, password_digest, locale, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var login string
	var password_digest string
	var locale string
	var created_at time.Time
	var updated_at time.Time

	row := db.QueryRow(sqlStatement, "4ac5d83a-7f39-486e-aebb-8fdd81fdb29c")
	err = row.Scan(&login, &password_digest, &locale, &created_at, &updated_at)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return
	case nil:
		fmt.Println(login)
		fmt.Println(password_digest)
		fmt.Println(locale)
		fmt.Println(created_at)
		fmt.Println(updated_at)
	default:
		panic(err)
	}

	router.Run()
}
