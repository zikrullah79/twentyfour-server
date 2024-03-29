package play

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math/rand"
	"net/http"
	"twentyfour.com/server/db"
	"twentyfour.com/server/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{}
var mongoClient *mongo.Client
var redisClient *redis.Client

func ConnectDb() {
	mongoClient = db.GetMongoClient()
	err := mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the mongo database", err)
	} else {
		log.Println("Connected to Mongodb")
	}

	redisClient, err = db.RedisCl()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database", err)
	} else {
		log.Println("Connected to Redis")
	}
}

// func Join(c *gin.Context) {
// 	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	defer conn.Close()

// 	for {
// 		msgType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			// log.Fatal(err)
// 			break
// 		}
// 		log.Printf("Received : %s", msg)
// 		if err := conn.WriteMessage(msgType, msg); err != nil {
// 			// log.Fatal(err)
// 			break
// 		}

// 	}
// }

func Join(c *gin.Context, r *model.Room) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	// defer conn.Close()
	playerType := model.RoomMaster

	if len(r.Players) > 0 {
		log.Println("roomMember")
		playerType = model.RoomMember
	}
	newPlayer := &model.Player{
		Id:         uint(rand.Uint64()),
		PlayerType: playerType,
		State:      model.PlayerJoining,
		Point:      0,
		Conn:       conn,
		Room:       r,
		Send:       make(chan []byte, 256),
	}

	newPlayer.Room.Register <- newPlayer
	go newPlayer.WritePlayerUpdate()
	go newPlayer.ReadPlayerUpdate()
}

func GetLeaderboard(c *gin.Context) {
	leaderboard, err := model.FindLeaderboard(redisClient, mongoClient, bson.D{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, leaderboard)
}

func PostLeaderboard(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		log.Fatal(err)
	}

	message, err := model.InsertLeaderboard(body, redisClient, mongoClient)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": message})
}
