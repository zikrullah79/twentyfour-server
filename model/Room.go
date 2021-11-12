package model

import (
	"log"
	"math/rand"
)

type Room struct {
	Id         uint64
	Status     int
	Players    map[*Player]bool
	Broadcast  chan []byte
	Register   chan *Player
	Unregister chan *Player
}

func NewRoom() *Room {
	return &Room{
		Id:         rand.Uint64(),
		Status:     GameNotStarted,
		Players:    make(map[*Player]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
	}
}

func (r *Room) Run() {
	for {
		// log.Println("p")
		select {
		case player := <-r.Register:

			// r.Players[player] = true
			if len(r.Players) < 8 {

				r.Players[player] = true
			}
			// log.Println(len(r.Players))
			// log.Println(player)
			// if len(r.Players) >= 2 {
			// 	Gamelog := &LogHistory{GameStarted, nil, nil}
			// 	res, err := json.Marshal(Gamelog)
			// 	if err != nil {
			// 		log.Println(err.Error())
			// 		// return
			// 	}
			// 	log.Printf("%v : marshal", string(res))
			// 	for players := range r.Players {
			// 		players.Send <- bytes.TrimSpace(res)
			// 	}

			// r.Broadcast <- bytes.TrimSpace([]byte("dddd"))
			// }
		case player := <-r.Unregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.Send)
			}

		case logplay := <-r.Broadcast:
			log.Printf("%v logplay", logplay)
			for player := range r.Players {
				select {
				case player.Send <- logplay:
				default:
					// log.Println("uhuy")
					close(player.Send)
					delete(r.Players, player)
				}

			}
		}
	}
}

// func Uint64() uint64 {
// 	buf := make([]byte, 8)
// 	rand.Read(buf) // Always succeeds, no need to check error
// 	return binary.LittleEndian.Uint64(buf)
// }
