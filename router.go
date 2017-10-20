package main

import (
  "log"
  "net/http"
  "os"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"

)

type User struct {
  ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
  Username string
  Email string
  Timestamp int32
}

func getUsers(c *gin.Context)  {
  log.Println("GET at api/v1/users")
  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  users := session.DB(db).C("user")
  var data []User
  users.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "users": data,
  })
}

func indexHandler(c *gin.Context)  {
  c.JSON(http.StatusOK, gin.H{
    "status": "success",
  })
}
