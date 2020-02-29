package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	cors "github.com/rs/cors/wrapper/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/prutya/todoapp-go/internal/app_config"
)

type User struct {
	Id             uuid.UUID
	Login          string
	PasswordDigest string
	Locale         string
	Roles          []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Todo struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserDTOCreateInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserDTO struct {
	Id     uuid.UUID `json:"id"`
	Login  string    `json:"login"`
	Locale string    `json:"locale"`
	Roles  []string  `json:"roles"`
}

func main() {
	config := app_config.New()

	db, err := sql.Open("postgres", config.DatabaseUrl)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/users", func(ctx *gin.Context) {
		var inputJson UserDTOCreateInput

		if err := ctx.ShouldBindJSON(&inputJson); err != nil {
			// TODO: Hide Go error
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		row := db.QueryRow(`SELECT 1 FROM users WHERE login=$1;`, inputJson.Login)

		var userExists int
		switch err = row.Scan(&userExists); err {
		case sql.ErrNoRows:
		case nil:
			ctx.JSON(http.StatusConflict, gin.H{})
			return
		default:
			panic(err)
		}

		var passwordDigest, err = bcrypt.GenerateFromPassword(
			[]byte(inputJson.Password),
			14,
		)

		if err != nil {
			panic(err)
		}

		var user UserDTO
		row = db.QueryRow(`
			INSERT INTO users (login, password_digest)
			VALUES ($1, $2)
			RETURNING id, login, locale, roles
		`, inputJson.Login, passwordDigest)

		err = row.Scan(&user.Id, &user.Login, &user.Locale, pq.Array(&user.Roles))

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusCreated, user)
	})

	router.GET("/users/:id", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Find a user by id if it is the current user
	})

	router.POST("/sessions", func(ctx *gin.Context) {
		// TODO: Create a session
	})

	router.POST("/todos", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Create a todo for current user
	})

	router.GET("/todos", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Return todos belonging to the current user
	})

	router.GET("/todos/:id", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Find a todo by id if it belongs to the current user
	})

	router.PATCH("/todos/:id", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Update a todo by id if it belongs to the current user
	})

	router.DELETE("/todos/:id", func(ctx *gin.Context) {
		// TODO: Authenticate
		// TODO: Delete a todo by id if it belongs to the current user
	})

	router.Run()
}
