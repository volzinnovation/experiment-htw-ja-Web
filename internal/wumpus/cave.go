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
	seen := map[int]bool{start: true}
	queue := []int{start}
	for len(queue) > 0 {
		room := queue[0]
		queue = queue[1:]
		for _, exit := range c.exits[room] {
			if !seen[exit] {
				seen[exit] = true
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

func sortedKeys(values map[int]bool) []int {
	keys := make([]int, 0, len(values))
	for value := range values {
		keys = append(keys, value)
	}
	sort.Ints(keys)
	return keys
}
