package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SlyMarbo/gmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	Timestamp int32         `json:"timestamp"`
}

type Game struct {
	ID                bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Owner             string        `json:"owner"`
	OwnerRole         string        `json:"ownerRole"`
	PublicRoom bool `json:"publicRoom"`
	Timestamp         int32         `json:"timestamp"`
	OperatorPassword  string        `json:"operatorPassword"`
	OperatorPort      string        `json:"operatorPort"`
	OperativePort     string        `json:"operativePort"`
	OperativeLocation string        `json:"operativeLocation"`
}

// EMAIL

func emailHandler(c *gin.Context) {
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

// AUTH

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

func signJWT(c *gin.Context) {
	claims := &jwt.MapClaims{
		"username": c.PostForm("username"),
		"email":    c.PostForm("email"),
		"password": c.PostForm("password"),
		"timestamp": c.PostForm("timestamp"),
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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"claims": claims,
			"valid": token.Valid,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	}
}

// GET ROUTES

func getUsers(c *gin.Context) {
	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	users := session.DB(db).C("user")
	timestamp, _ := strconv.Atoi(c.Query("timestamp"))
	log.Println(timestamp)
	if timestamp != 0 {
		var data User
		users.Find(bson.M{"timestamp": timestamp}).One(&data)
		c.JSON(http.StatusOK, gin.H{
			"user": data,
		})
	} else {
		var data []User
		users.Find(nil).All(&data)
		c.JSON(http.StatusOK, gin.H{
			"users": data,
		})
	}
}

func getGames(c *gin.Context) {
	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	games := session.DB(db).C("game")
	timestamp, _ := strconv.Atoi(c.Query("timestamp"))
	log.Println(timestamp)

	if timestamp != 0 {
		var data Game
		games.Find(bson.M{"timestamp": timestamp}).One(&data)
		c.JSON(http.StatusOK, gin.H{
			"game": data,
		})
	} else {
		var data []Game
		games.Find(nil).All(&data)
		c.JSON(http.StatusOK, gin.H{
			"games": data,
		})
	}
}


// POST ROUTES

func postUser(c *gin.Context) {
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
	var check []User
	users.Find(bson.M{"username": username}).All(&check)
	log.Println(len(check))
	if len(check) > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User already exists, you must create a unique username",
		})
	} else {
		err = users.Insert(&User{
			Username:  username,
			Email:     email,
			Timestamp: timestamp,
			Password:  password,
		})
		if err != nil {
			log.Println(err)
		}
		var data []User
		users.Find(bson.M{"timestamp": timestamp}).All(&data)
		c.JSON(http.StatusOK, gin.H{
			"users": data,
		})

	}
}

func postGame(c *gin.Context) {
	owner := c.PostForm("owner")
	ownerRole := c.PostForm("ownerRole")
	publicRoom, _ := strconv.ParseBool(c.PostForm("publicRoom"))
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

	var check []Game
	games.Find(bson.M{"owner": owner}).All(&check)
	if len(check) > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Each user can only own one game, You must delete your game from the game-lobby before creating a new game",
		})
	} else {
		err = games.Insert(&Game{
			Owner:             owner,
			OwnerRole:         ownerRole,
			PublicRoom: publicRoom,
			Timestamp:         timestamp,
			OperatorPassword:  operatorPassword,
			OperatorPort:      operatorPort,
			OperativePort:     operativePort,
			OperativeLocation: operativeLocation,
		})
		if err != nil {
			log.Println(err)
		}
		var data []Game
		games.Find(bson.M{"timestamp": timestamp}).All(&data)
		c.JSON(http.StatusOK, gin.H{
			"games": data,
		})
	}


}

// DELETE ROUTES

func deleteUser(c *gin.Context) {
	timestamp, _ := strconv.Atoi(c.Param("timestamp"))

	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	users := session.DB(db).C("user")
	err = users.Remove(bson.M{"timestamp": timestamp})
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func deleteGame(c *gin.Context) {
	timestamp, _ := strconv.Atoi(c.Param("timestamp"))

	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	games := session.DB(db).C("game")
	err = games.Remove(bson.M{"timestamp": timestamp})
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

func deleteAll(c *gin.Context)  {
	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	games := session.DB(db).C("game")
	users := session.DB(db).C("user")
	games.DropCollection()
	users.DropCollection()
	c.JSON(http.StatusOK, gin.H{
		"message": "ALL DELETED",
	})
}
// INDEX HANDLER

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
