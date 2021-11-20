package util

import (
	"sort"

	"twentyfour.com/server/model"
)

type Players []*model.Player

func (p Players) Len() int           { return len(p) }
func (p Players) Less(i, j int) bool { return p[i].Point < p[j].Point }
func (p Players) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func SortByPoint(r *model.Room) Players {
	var players Players
	for _, v := range r.Players {
		players = append(players, v)
	}
	sort.Sort(players)
	return players
}
