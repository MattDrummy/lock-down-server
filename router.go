package main

import (
  "log"
  "net/http"
  "os"
  "time"
  "strconv"
  "fmt"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
  "github.com/dgrijalva/jwt-go"
  "github.com/SlyMarbo/gmail"
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
  Password string `json:"password"`
  Timestamp int32 `json:"timestamp"`
  OperatorPassword string `json:"operatorPassword"`
  OperatorPort string `json:"operatorPort"`
  OperativePort string `json:"operativePort"`
  OperativeLocation string `json:"operativeLocation"`
}

// EMAIL

func emailHandler(c *gin.Context)  {
  email := gmail.Compose(c.PostForm("subject"), c.PostForm("body"))
  email.From = os.Getenv("EMAIL_SENDER")
  email.Password = os.Getenv("EMAIL_PASSWORD")
  email.ContentType = "text/html; charset=utf-8"
  email.AddRecipient(c.PostForm("recipient"))
  if err := email.Send(); err != nil {
    c.JSON(http.StatusOK, gin.H{
      "error": err,
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
      "message": "sent",
    })
  }
}

// JWT

func createToken(claims *jwt.MapClaims) string {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  secretKey := os.Getenv("HMAC_SECRET_KEY")
  hmacSampleSecret := []byte(secretKey)
  tokenString, err := token.SignedString(hmacSampleSecret)
  if err != nil {
    log.Println(err)
  }
  return tokenString

}

func signJWT(c *gin.Context)  {
  claims := &jwt.MapClaims{
    "username": c.PostForm("username"),
    "email": c.PostForm("email"),
    "password": c.PostForm("password"),
  }
  tokenString := createToken(claims)
  c.JSON(http.StatusOK, gin.H{
    "tokenString": tokenString,
  })
}

func verifyJWT(c *gin.Context) {
  tokenString := c.PostForm("tokenString")
  secretKey := os.Getenv("HMAC_SECRET_KEY")
  hmacSampleSecret := []byte(secretKey)

  token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }
    return hmacSampleSecret, nil
  })


  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    c.JSON(http.StatusOK, gin.H{
      "claims": claims,
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
      "error": err,
    })
  }
}

// GET ROUTES

func getUsers(c *gin.Context)  {
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

// POST ROUTES

func postUser(c *gin.Context)  {
  username := c.PostForm("username")
  email := c.PostForm("email")
  password := c.PostForm("password")
  timestamp := int32(time.Now().Unix())

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
  password := c.PostForm("password")
  timestamp := int32(time.Now().Unix())
  operatorPassword := c.PostForm("operatorPassword")
  operatorPort := c.PostForm("operatorPort")
  operativePort := c.PostForm("operativePort")
  operativeLocation := c.PostForm("operativeLocation")

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
    Password: password,
    Timestamp: timestamp,
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

// DELETE ROUTES

func deleteUser(c *gin.Context)  {
  time, _ := strconv.Atoi(c.Param("time"))

  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  users := session.DB(db).C("user")
  err = users.Remove(bson.M{"timestamp":time})
  if err != nil {
    log.Println(err)
  }

  c.JSON(http.StatusOK, gin.H{
    "message": "deleted",
  })
}

func deleteGame(c *gin.Context){
  time, _ := strconv.Atoi(c.Param("time"))

  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  games := session.DB(db).C("game")
  err = games.Remove(bson.M{"timestamp":time})
}

// PUT ROUTES

func patchUser(c *gin.Context){
  username := c.PostForm("username")
  email := c.PostForm("email")
  password := c.PostForm("password")
  time, _ := strconv.Atoi(c.Param("time"))

  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  users := session.DB(db).C("user")

  update := bson.M{
    "username": username,
    "email": email,
    "password": password,
  }
  change := bson.M{"$set": update}
  users.Update(bson.M{"timestamp": time}, change)

  var data []User
  users.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "users": data,
  })
}

func patchGame(c *gin.Context){
  owner := c.PostForm("owner")
  room := c.PostForm("room")
  password := c.PostForm("password")
  operatorPassword := c.PostForm("operatorPassword")
  operatorPort := c.PostForm("operatorPort")
  operativePort := c.PostForm("operativePort")
  operativeLocation := c.PostForm("operativeLocation")
  time, _ := strconv.Atoi(c.Param("time"))

  mongo := os.Getenv("MONGODB_URI")
  db := os.Getenv("DATABASE_NAME")
  session, err := mgo.Dial(mongo)
  if err != nil {
    log.Println(err)
  }
  games := session.DB(db).C("game")

  update := bson.M{
    "owner": owner,
    "room": room,
    "password": password,
    "operatorPassword": operatorPassword,
    "operatorPort": operatorPort,
    "operativePort": operativePort,
    "operativeLocation": operativeLocation,

  }
  change := bson.M{"$set": update}
  games.Update(bson.M{"timestamp": time}, change)

  var data []Game
  games.Find(nil).All(&data)
  c.JSON(http.StatusOK, gin.H{
    "games": data,
  })
}


// INDEX HANDLER

func indexHandler(c *gin.Context)  {
  c.JSON(http.StatusOK, gin.H{
    "status": "success",
  })
}
