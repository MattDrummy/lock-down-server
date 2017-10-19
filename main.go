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
  router := gin.Default()
  router.Use(cors.New(cors.Config{
    AllowOrigins: []string{os.Getenv("SITE_URL")},
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
  router.Run(":" + os.Getenv("PORT"))

}
