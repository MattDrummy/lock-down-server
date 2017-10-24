package main

import (
  "log"
  "net/http"
  "os"
  "time"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"

)

type User struct {
  ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
  Username string `json:"username"`
  Email string `json:"email"`
  Password string `json:"password"`
  Timestamp int32 `json:"timestamp"`
}

type Game struct {
  ID bson.ObjectId `json:"id" bson:"_id,omitempty"`
  Owner string `json:"owner"`
  Room string `json:"room"`
  Operator string `json:"operator"`
  Operative string `json:"operative"`
  Password string `json:"password"`
  Timestamp int32 `json:"timestamp"`
  OperatorLocation string `json:"operatorLocation"`
  OperatorPassword string `json:"operatorPassword"`
  OperatorPort int32 `json:"operatorPort"`
  OperativePort int32 `json:"operativePort"`
  OperativeLocation string `json:"operativeLocation"`
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

func getGames(c *gin.Context)  {
  log.Println("POST at api/v1/games")
  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  games := session.DB(db).C("game")
  var data []Game
  games.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "games": data,
  })
}

func postUser(c *gin.Context)  {
  log.Println("POST at api/v1/users")
  username := c.PostForm("username")
  email := c.PostForm("email")
  password := c.PostForm("password")
  timestamp := int32(time.Now().Unix())
  log.Println("GET at api/v1/users")
  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  users := session.DB(db).C("user")
  err = users.Insert(&User{
    Username: username,
    Email: email,
    Timestamp: timestamp,
    Password: password,
  })
  if err != nil {
    log.Println(err)
  }
  var data []User
  users.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "users": data,
  })

}

func postGame(c *gin.Context)  {
  owner := c.PostForm("owner")
  room := c.PostForm("room")
  operator := c.PostForm("operator")
  operative := c.PostForm("operative")
  password := c.PostForm("password")
  timestamp := int32(time.Now().Unix())
  operatorLocation := c.PostForm("operatorLocation")
  operatorPassword := c.PostForm("operatorPassword")
  operatorPort := c.PostForm("operatorPort")
  operativePort := c.PostForm("operativePort")
  operativeLocation := c.PostForm("operativeLocation")
  log.Println("GET at api/v1/users")
  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  games := session.DB(db).C("game")
  err = games.Insert(&Game{
    Owner: owner,
    Room: room,
    Operator: operator,
    Operative: operative,
    Password: password,
    Timestamp: timestamp,
    OperatorLocation: operatorLocation,
    OperatorPassword: operatorPassword,
    OperatorPort: operatorPort,
    OperativePort: operativePort,
    OperativeLocation: operativeLocation,
  })
  if err != nil {
    log.Println(err)
  }
  var data []Game
  games.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "games": data,
  })

}

func deleteUser()  {

}

func deleteGame(){

}
func indexHandler(c *gin.Context)  {
  c.JSON(http.StatusOK, gin.H{
    "status": "success",
  })
}
