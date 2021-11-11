package model

import (
	"math/rand"
)

type Room struct {
	Id         uint
	Status     int
	Players    map[*Player]bool
	Broadcast  chan []byte
	Register   chan *Player
	Unregister chan *Player
}

func NewRoom() *Room {
	return &Room{
		Id:         uint(rand.Uint64()),
		Status:     0,
		Players:    make(map[*Player]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
	}
}

func (r *Room) Run() {
	for {
		select {
		case player := <-r.Register:
			r.Players[player] = true
		case player := <-r.Unregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.Send)
			}

		case logplay := <-r.Broadcast:
			for player := range r.Players {
				select {
				case player.Send <- logplay:
				default:
					close(player.Send)
					delete(r.Players, player)
				}

			}
		}
	}
}
