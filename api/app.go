package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

type UserDTOShowInputUri struct {
	Id string `uri:"id" binding:"required"`
}

type UserDTOShowInputHeaders struct {
	Authorization string `uri:"Authorization"`
}

type SessionDTOCreateInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		row := db.QueryRow(`
			SELECT 1
			FROM users
			WHERE login=$1
		`, inputJson.Login)

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
			config.BcryptCost,
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
		var paramsUri UserDTOShowInputUri

		if err := ctx.ShouldBindUri(&paramsUri); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		var paramsHeaders UserDTOShowInputHeaders

		if err = ctx.ShouldBindHeader(&paramsHeaders); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		token, tokenErr := jwt.Parse(
			paramsHeaders.Authorization,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(config.AuthSecret), nil
			},
		)

		if tokenErr != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if !(paramsUri.Id == "current" || paramsUri.Id == claims["user_id"]) {
			fmt.Println(paramsUri.Id)
			ctx.JSON(http.StatusForbidden, gin.H{})
			return
		}

		var user UserDTO
		row := db.QueryRow(`
			SELECT id, login, locale, roles
			FROM users
			WHERE id=$1
		`, claims["user_id"])

		dbErr := row.Scan(&user.Id, &user.Login, &user.Locale, pq.Array(&user.Roles))

		if dbErr != nil {
			panic(dbErr)
		}

		ctx.JSON(http.StatusOK, user)
	})

	router.POST("/sessions", func(ctx *gin.Context) {
		var inputJson SessionDTOCreateInput
		if err := ctx.ShouldBindJSON(&inputJson); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		var id uuid.UUID
		var passwordDigest string
		row := db.QueryRow(`
			SELECT id, password_digest
			FROM users
			WHERE login=$1
		`, inputJson.Login)

		switch err = row.Scan(&id, &passwordDigest); err {
		case sql.ErrNoRows:
			// TODO: Prevent a timing attack
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		case nil:
		default:
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(passwordDigest),
			[]byte(inputJson.Password),
		)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		currentTimeUnix := time.Now().Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": id,
			"exp":     currentTimeUnix + config.AuthExpirySeconds,
			"iat":     currentTimeUnix,
		})

		tokenString, tokenErr := token.SignedString([]byte(config.AuthSecret))

		if tokenErr != nil {
			panic(tokenErr)
		}

		ctx.JSON(http.StatusCreated, gin.H{"jwt": tokenString})
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
