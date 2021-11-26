package model

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"twentyfour.com/server/util"
)

type Leaderboard struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	ProfileID string             `bson:"profile_id" json:"profile_id"`
	Name      string             `json:"name"`
	Score     int64              `json:"score"`
}

func FindLeaderboard(redisClient *redis.Client, client *mongo.Client, filter bson.D) ([]*Leaderboard, error) {
	var leaderboard []*Leaderboard
	leaderboardString, err := redisClient.Get(context.TODO(), util.LeaderboardRedisKey).Result()
	if err != nil {
		if err == redis.Nil || leaderboardString == "" {
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

			jsonLeaderboard, err := json.Marshal(leaderboard)
			if err != nil {
				return nil, err
			}

			err = redisClient.Set(context.TODO(), util.LeaderboardRedisKey, string(jsonLeaderboard), 24*time.Hour).Err()
			if err != nil {
				return nil, err
			}

			return leaderboard, nil
		} else {
			return nil, err
		}
	}

	err = json.Unmarshal([]byte(leaderboardString), &leaderboard)
	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

func InsertLeaderboard(body []byte, redisClient *redis.Client, client *mongo.Client) (string, error) {
	var leaderboardInterface []interface{}
	if err := json.Unmarshal(body, &leaderboardInterface); err != nil {
		return "", err
	}

	collection := client.Database("config").Collection("leaderboard")
	_, err := collection.InsertMany(context.TODO(), leaderboardInterface)
	if err != nil {
		return "", err
	}

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
		return "", err
	}

	var leaderboard []*Leaderboard

	for cur.Next(context.TODO()) {
		var board Leaderboard
		err = cur.Decode(&board)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
			return "", err
		}

		leaderboard = append(leaderboard, &board)
	}

	leadByte, err := json.Marshal(leaderboard)
	if err != nil {
		return "", err
	}

	err = redisClient.Set(context.TODO(), util.LeaderboardRedisKey, string(leadByte), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return "Data inserted", nil
}
