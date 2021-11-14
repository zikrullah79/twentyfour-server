package play

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"twentyfour.com/server/db"
	"twentyfour.com/server/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{}
var mongoClient *mongo.Client

func ConnectMongodb() {
	mongoClient = db.GetMongoClient()
	err := mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected")
	}
}

func Join(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// log.Fatal(err)
			break
		}
		log.Printf("Received : %s", msg)
		if err := conn.WriteMessage(msgType, msg); err != nil {
			// log.Fatal(err)
			break
		}
	}
}

func GetLeaderboard(c *gin.Context) {
	leaderboard, err := model.FindLeaderboard(mongoClient, bson.D{})
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}
	c.IndentedJSON(http.StatusOK, leaderboard)
}
