package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"twentyfour.com/server/util"
)

type Leaderboard struct {
	ProfileId string `json:"profile_id"`
	Name      string `json:"name"`
	Score     int64  `json:"score"`
}

func FindLeaderboard(redisClient *redis.Client, client *mongo.Client, filter bson.D) ([]*Leaderboard, error) {
	var leaderboard []*Leaderboard
	leaderboardString, err := redisClient.Get(context.TODO(), util.LeaderboardRedisKey).Result()
	fmt.Printf("leaderboardstring %s\n", leaderboardString)
	if err != nil {
		if err == redis.Nil || leaderboardString == "" {
			fmt.Println("disini")
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
