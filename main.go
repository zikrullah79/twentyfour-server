package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"twentyfour.com/server/db"
	"twentyfour.com/server/model"
	"twentyfour.com/server/play"
)

func main() {
	mongoClient := db.GetMongoClient()
	err := mongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected")
	}

	leaderboard := model.GetLeaderboard(mongoClient, bson.D{})
	for _, board := range leaderboard {
		log.Println(board.ProfileId, board.Name, board.Score)
	}

	router := gin.Default()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/join", play.Join)
	router.POST("/create")
	router.Run()
}
