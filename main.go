package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	port := os.Getenv("PORT")
	client := os.Getenv("CLIENT_URL")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{client},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	socket, err := socketio.NewServer(nil)
	if err != nil {
		log.Println(err)
	}

	// INDEX ROUTE

	router.GET("/", indexHandler)
	router.DELETE("/api/v1/deleteAll", deleteAll)

	// EMAIL

	router.POST("/email", emailHandler)

	// AUTH

	router.POST("/signJWT", signJWT)
	router.POST("/verifyToken", verifyToken)
	router.POST("/logIn", logIn)

	// SOCKETS

	socket.On("connection", socketConnectionHandler)
	socket.On("error", socketErorrHandler)
	router.GET("/socket.io/", gin.WrapH(socket))
	router.POST("/socket.io/", gin.WrapH(socket))

	// USER DB

	router.GET("api/v1/users", getUsers)
	router.POST("api/v1/users", postUser)
	router.DELETE("api/v1/users/:timestamp", deleteUser)
	router.PUT("api/v1/users/:timestamp", updateUser)
	// GAME DB

	router.GET("api/v1/games", getGames)
	router.POST("api/v1/games", postGame)
	router.DELETE("api/v1/games/:timestamp", deleteGame)
	router.PUT("api/v1/games/:timestamp", updateGame)

	// LOG AND RUN

	log.Println("server ready on port: " + port)
	log.Println("CORS allowed on " + client)
	router.Run(":" + port)
}
