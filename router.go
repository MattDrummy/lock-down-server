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

func logIn(c *gin.Context) {
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
	time := c.Query("timestamp")

	if len(time) != 0 {
		timestamp, _ := strconv.Atoi(time)
		var data User
		users.Find(bson.M{"timestamp": timestamp}).One(&data)
		if data.Timestamp == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "nothing found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"user": data,
			})
		}

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
	time := c.Query("timestamp")
	owner := c.Query("owner")

	if len(time) != 0 {
		timestamp, _ := strconv.Atoi(time)
		var data Game
		games.Find(bson.M{"timestamp": timestamp}).One(&data)
		if data.Timestamp == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "nothing found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"game": data,
			})
		}
	} else if len(owner) != 0 {
		var data Game
		games.Find(bson.M{"owner": owner}).One(&data)
		if data.Timestamp == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "nothing found",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"game": data,
			})
		}

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
	if len(check) > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "This username is already taken.",
		})
	} else {
		err = users.Insert(&User{
			Username:  username,
			Email:     email,
			Timestamp: timestamp,
			Password:  password,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed insert",
			})
		} else {
			var data []User
			users.Find(bson.M{"timestamp": timestamp}).All(&data)
			c.JSON(http.StatusOK, gin.H{
				"users": data,
			})
		}
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
			"error": "Each user may only have one game in the lobby for their username.  An existing game exists under this username, would you like to delete that game and start again?  Clicking OK will reload the page, and you will have to recreate the game.",
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
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed insert",
			})
		} else {
			var data []Game
			err = games.Find(bson.M{"timestamp": timestamp}).All(&data)
			c.JSON(http.StatusOK, gin.H{
				"games": data,
			})

		}
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
		c.JSON(http.StatusOK, gin.H{
			"error": "failed remove",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{})

	}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed remove",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{})

	}
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

// PUT ROUTES

func updateUser(c *gin.Context)  {
	timestamp, _ := strconv.Atoi(c.Param("timestamp"))
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	update := bson.M{
		"username": username,
		"email": email,
		"password": password,
	}
	change := bson.M{
		"$set": update,
	}
	mongo := os.Getenv("MONGODB_URI")
	db := os.Getenv("DATABASE_NAME")
	session, err := mgo.Dial(mongo)
	if err != nil {
		log.Println(err)
	}
	users := session.DB(db).C("user")
	users.Update(bson.M{"timestamp": timestamp}, change)
	var data User
	users.Find(bson.M{"timestamp": timestamp}).One(&data)
	c.JSON(http.StatusOK, gin.H{
		"user": data,
	})
}

// INDEX HANDLER

func indexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
