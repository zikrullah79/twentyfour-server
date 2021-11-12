package model

const (
	GameLog = iota
	PlayerLog
)

const (
	EndGame = iota
	PlayerPointing
	LastPlayer
	KeyCorrect
	KeyUncorrect
	KnowTheKey
	Unresolve
	FalseUnresolve
	PointIncrease
	AFK
	PlayerJoining
)

const (
	NewQuestion = iota
	GameStarted
	GameStopped
	GameNotStarted
)

type LogHistory struct {
	Type            int
	PlayerLogData   *PlayerLogData
	GamePlayLogData *GameLogData
}
type GameLogData struct {
	Question []*int
}
type PlayerLogData struct {
	Id  int64
	Key int
}
