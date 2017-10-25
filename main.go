package main

import (
  "os"
  "log"

  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
  "github.com/gin-contrib/cors"
  "github.com/googollee/go-socket.io"

)

func main()  {
  err := godotenv.Load()
  if err != nil {
    log.Println(err)
  }
  port := os.Getenv("PORT")
  client := os.Getenv("CLIENT_URL")
  router := gin.Default()
  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{client},
    AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowHeaders: []string{"Origin"},
    ExposeHeaders: []string{"Content-Length"},
    AllowCredentials: true,
  }))
  socket, err := socketio.NewServer(nil)
  if err != nil {
    log.Println(err)
  }
  socket.On("connection", socketConnectionHandler)
  socket.On("error", socketErorrHandler)
  router.GET("/", indexHandler)
  router.GET("api/v1/users", getUsers)
  router.GET("api/v1/games", getGames)
  router.POST("api/v1/users", postUser)
  router.POST("api/v1/games", postGame)
  router.DELETE("api/v1/users/:time", deleteUser)
  router.DELETE("api/v1/games/:time", deleteGame)
  router.PUT("api/v1/users/:time", patchUser)
  router.PUT("api/v1/games/:time", patchGame)
  router.GET("/socket.io/", gin.WrapH(socket))
  router.POST("/socket.io/", gin.WrapH(socket))
  log.Println("server ready on port: " + port)
  log.Println("CORS allowed on " + client)
  router.Run(":" + port)
}
