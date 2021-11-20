package services

import (
	"math/rand"
	"time"

	"github.com/shogo82148/go-shuffle"
)

const (
	maxShuffle = 10
	minShuffle = 1
)

func GetCardSet() *[]int {
	var set []int
	for i := 1; i < 10; i++ {
		set = append(set, i)
		// set = append(set, i)
		// set = append(set, i)
		// set = append(set, i)
	}
	return &set
}

func Get4Card(cardSet []int) (*[]int, *[]int) {
	cards, currSet := cardSet[0:4], cardSet[4:]
	return &cards, &currSet
}

func Shuffle(cardSet *[]int) *[]int {
	rand.Seed(time.Now().UnixNano())
	cards := *cardSet
	for i := 0; i < (rand.Intn(maxShuffle-minShuffle+1) + minShuffle); i++ {
		shuffle.Ints(cards)
	}
	return &cards
}
