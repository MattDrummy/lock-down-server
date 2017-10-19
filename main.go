package main

import (
  "io"
  "os"
  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
)

func main()  {
  gin.DisableConsoleColor()
  f, _ := os.Create("gin.log")
  gin.DefaultWriter = io.MultiWriter(f)
  
}
