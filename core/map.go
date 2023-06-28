package core

import "math/rand"

type MapCoord struct {
	X, Y int
}

type Map struct {
	Tiles [][]*Tile
}

type Tile struct {
	X             int
	Y             int
	ObverseEffect Effect
	ReverseEffect Effect
	Reversed      bool
	HasCreature   bool
}

func (t *Tile) GetActiveEffect() Effect {
	if t.Reversed {
		return t.ReverseEffect
	} else {
		return t.ObverseEffect
	}
}

func RandomTile(x, y int, random *rand.Rand) *Tile {
	return &Tile{
		X:             x,
		Y:             y,
		Reversed:      false,
		ObverseEffect: RandomEffect(x, y, random),
		ReverseEffect: RandomEffect(x, y, random),
		HasCreature:   false,
	}
}

const MAP_HEIGHT = 3
const MAP_WIDTH = 5

func NewMap(random *rand.Rand) *Map {
	tiles := make([][]*Tile, MAP_WIDTH)
	for i := 0; i < MAP_WIDTH; i++ {
		j := i % 2
		tiles[i] = make([]*Tile, 2*MAP_HEIGHT-1+j)
		for ; j < 2*MAP_HEIGHT; j += 2 {
			tiles[i][j] = RandomTile(i, j, random)
		}
	}
	return &Map{
		Tiles: tiles,
	}
}
