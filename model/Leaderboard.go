package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Leaderboard struct {
	ProfileId string `json:"profile_id"`
	Name      string `json:"name"`
	Score     int64  `json:"score"`
}

func FindLeaderboard(client *mongo.Client, filter bson.D) ([]*Leaderboard, error) {
	var leaderboard []*Leaderboard
	collection := client.Database("config").Collection("leaderboard")

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var board Leaderboard
		err = cur.Decode(&board)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
			return nil, err
		}

		leaderboard = append(leaderboard, &board)
	}

	return leaderboard, nil
}

//func InsertLeaderboard(client *mongo.Client, leaderboard []interface{}) interface{} {
//	collection := client.Database("config").Collection("leaderboard")
//
//	result, err := collection.InsertMany(context.TODO(), leaderboard)
//
//	cur, err := collection.Find(context.TODO(), filter)
//	if err != nil {
//		log.Fatal("Error on Finding all the documents", err)
//	}
//
//	for cur.Next(context.TODO()) {
//		var board Leaderboard
//		err = cur.Decode(&board)
//		if err != nil {
//			log.Fatal("Error on Decoding the document", err)
//		}
//
//		leaderboard = append(leaderboard, &board)
//	}
//
//	return leaderboard
//}
