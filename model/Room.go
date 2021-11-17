package model

import (
	"log"
	"math/rand"
)

type Room struct {
	Id         uint64
	Status     int
	Players    map[uint]*Player
	Broadcast  chan []byte
	Register   chan *Player
	Unregister chan *Player
}

func NewRoom() *Room {
	return &Room{
		Id:         rand.Uint64(),
		Status:     GameNotStarted,
		Players:    make(map[uint]*Player),
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

				r.Players[player.Id] = player
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
			if _, ok := r.Players[player.Id]; ok {
				delete(r.Players, player.Id)
				close(player.Send)
			}

		case logplay := <-r.Broadcast:
			log.Printf("%v logplay", logplay)
			if r.Status == NewQuestion {
				log.Print("yuhuu")
			}
			for _, player := range r.Players {
				select {
				case player.Send <- logplay:
				default:
					// log.Println("uhuy")
					close(player.Send)
					delete(r.Players, player.Id)
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
