package model

const (
	GameLog = iota
	PlayerLog
)

const (
	EndGame = iota
	KeyCorrect
	KeyUncorrect
	FalseUnresolve
	PointIncrease
	AFK
	PlayerJoining
)
const (
	NewQuestion = iota
	GameStopped
	GameNotStarted
)
const (
	PlayerPointing = iota
	ClaimUnresolve
	StartGame
	ClaimSolution
	AnswerTheQuestion
)

const (
	PlayerPointed = iota
	LastPlayer
	KnowTheSolution
	Thinking
	Unresolve
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
	Id  uint
	Key int
}
