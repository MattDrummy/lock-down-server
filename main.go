package main

import (
  "io"
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
  gin.DisableConsoleColor()
  f, err := os.Create("gin.log")
  if err != nil {
    log.Println(err)
  }
  gin.DefaultWriter = io.MultiWriter(f)

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
  router.GET("/socket.io/", gin.WrapH(socket))
  router.POST("/socket.io/", gin.WrapH(socket))
  log.Println("server ready on port: " + port)
  log.Println("CORS allowed on " + client)
  router.Run(":" + port)

}
