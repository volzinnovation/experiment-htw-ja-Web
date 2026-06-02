package wumpus

import (
	"fmt"
	"slices"
	"sort"
)

type Cave struct {
	exits map[int][]int
}

type Hazard string

const (
	HazardWumpus Hazard = "Wumpus"
	HazardPit    Hazard = "Pit"
	HazardBats   Hazard = "Bats"
)

func NewCave() Cave {
	exits := map[int][]int{
		1:  {2, 5, 8},
		2:  {1, 3, 10},
		3:  {2, 4, 12},
		4:  {3, 5, 14},
		5:  {1, 4, 6},
		6:  {5, 7, 15},
		7:  {6, 8, 17},
		8:  {1, 7, 9},
		9:  {8, 10, 18},
		10: {2, 9, 11},
		11: {10, 12, 19},
		12: {3, 11, 13},
		13: {12, 14, 20},
		14: {4, 13, 15},
		15: {6, 14, 16},
		16: {15, 17, 20},
		17: {7, 16, 18},
		18: {9, 17, 19},
		19: {11, 18, 20},
		20: {13, 16, 19},
	}
	return Cave{exits: exits}
}

func (c Cave) Exits(room int) ([]int, error) {
	exits, ok := c.exits[room]
	if !ok {
		return nil, fmt.Errorf("invalid room %d", room)
	}
	return append([]int(nil), exits...), nil
}

func (c Cave) HasTunnel(from, to int) bool {
	exits, ok := c.exits[from]
	if !ok {
		return false
	}
	for _, exit := range exits {
		if exit == to {
			return true
		}
	}
	return false
}

func (c Cave) ReachableFrom(start int) []int {
	if _, ok := c.exits[start]; !ok {
		return nil
	}
	seen := map[int]struct{}{start: {}}
	queue := []int{start}
	for len(queue) > 0 {
		room := queue[0]
		queue = queue[1:]
		for _, exit := range c.exits[room] {
			if _, ok := seen[exit]; !ok {
				seen[exit] = struct{}{}
				queue = append(queue, exit)
			}
		}
	}
	return sortedKeys(seen)
}

func (c Cave) AdjacentHazards(room int, setup Setup) []Hazard {
	exits, ok := c.exits[room]
	if !ok {
		return nil
	}
	var hazards []Hazard
	seen := map[Hazard]bool{}
	for _, exit := range exits {
		for _, hazard := range hazardsInRoom(exit, setup) {
			if !seen[hazard] {
				seen[hazard] = true
				hazards = append(hazards, hazard)
			}
		}
	}
	return hazards
}

func hazardsInRoom(room int, setup Setup) []Hazard {
	var hazards []Hazard
	if setup.Wumpus == room {
		hazards = append(hazards, HazardWumpus)
	}
	if slices.Contains(setup.Pits, room) {
		hazards = append(hazards, HazardPit)
	}
	if slices.Contains(setup.Bats, room) {
		hazards = append(hazards, HazardBats)
	}
	return hazards
}

func sortedKeys(values map[int]struct{}) []int {
	keys := make([]int, 0, len(values))
	for value := range values {
		keys = append(keys, value)
	}
	sort.Ints(keys)
	return keys
}

// mutate4go-manifest-begin
// {"version":1,"tested_at":"2026-06-02T10:00:40-05:00","module_hash":"f0e1660602bb2be33b1673d4e39b2b9cf3b88e41f0656707cca8b4942c141d79","functions":[{"id":"func/NewCave","name":"NewCave","line":21,"end_line":45,"hash":"682215dabf2f6e64840ef529aa393c7d6ca3bf4410c70c3d42ea20f6dfc59004"},{"id":"func/Cave.Exits","name":"Cave.Exits","line":47,"end_line":53,"hash":"c92e8c40813bb0a35e13ba03f6548e8664cb53593de31fe0c870dcfb81be9925"},{"id":"func/Cave.HasTunnel","name":"Cave.HasTunnel","line":55,"end_line":66,"hash":"eecfc59abe2742483633cf907329d282d27bad29d7483b6f0cf638a4e0103e82"},{"id":"func/Cave.ReachableFrom","name":"Cave.ReachableFrom","line":68,"end_line":85,"hash":"6a0cfa86b9149f3767e010f2867027382c2e8a3751a4bb62a749be5dfaeeafb6"},{"id":"func/Cave.AdjacentHazards","name":"Cave.AdjacentHazards","line":87,"end_line":103,"hash":"19af123b7236f2bd5440113862bf6371c3909a2e55774d1a477a2b833f94cd76"},{"id":"func/hazardsInRoom","name":"hazardsInRoom","line":105,"end_line":117,"hash":"4d0d4e47f15b4acfeb6c636780791a5e143322d704fc1b0fbe06aff0bd21bf1e"},{"id":"func/sortedKeys","name":"sortedKeys","line":119,"end_line":126,"hash":"95fefa5db632b9f596faaf10d8e48f657a28ac79912c2c096edce6b134416e86"}]}
// mutate4go-manifest-end
