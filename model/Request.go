package model

type PlayerAct struct {
	Type int
	Data string
}
type RequestStateGame struct {
	State int
}

func GetRequestStateGame(state int) *RequestStateGame {
	return &RequestStateGame{
		State: state,
	}
}
