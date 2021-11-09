package model

type Player struct {
	Id    uint
	State int
	Point float64
}

type Room struct {
	Id         uint
	Status     int
	Players    map[*Player]bool
	Broadcast  chan []byte
	Register   chan *Player
	Unregister chan *Player
}

func newRoom() *Room {
	return &Room{
		Id:         uint(998),
		Status:     0,
		Players:    make(map[*Player]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
	}
}

func (r *Room) Run() {
	// for{
	// 	select{
	// 		case
	// 	}
	// }
}
