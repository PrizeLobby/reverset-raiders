package core

import (
	"fmt"
	"math/rand"
)

//go:generate stringer -type Species
type Species int

const (
	NO_SPECIES Species = iota
	Duck
	Tortoise
	Capybara
)

//go:generate stringer -type CreatureColor
type CreatureColor int

const (
	NO_COLOR CreatureColor = iota
	Red
	Green
	Blue
)

//go:generate stringer -type Alignment
type Alignment int

const (
	EAST Alignment = 1
	WEST Alignment = -1
)

func (a Alignment) Opposite() Alignment {
	return Alignment(-int(a))
}

type Creature struct {
	X         int
	Y         int
	Alignment Alignment
	Id        int
	Power     int
	Color     CreatureColor
	Species   Species
	Removed   bool
}

func (c *Creature) Name() string {
	return c.Color.String() + c.Species.String()
}

func (c *Creature) String() string {
	return fmt.Sprintf("{Creature x %d y %d align %s power %d color %s species %s removed %t}", c.X, c.Y, c.Alignment, c.Power, c.Color, c.Species, c.Removed)
}

func (c *Creature) ApplyEffect(e Effect) bool {
	if c.Removed {
		return false
	}

	if c.X < 0 || c.X >= MAP_WIDTH {
		return false
	}

	if e.Targets == TILE && (e.X != c.X || e.Y != c.Y) {
		return false
	}

	if (e.ColorCondition == NO_COLOR || e.ColorCondition == c.Color) && (e.SpeciesCondition == NO_SPECIES || e.SpeciesCondition == c.Species) {
		c.Power += e.Value
		return true
	}
	return false
}

func ShuffledCreatureBag(random *rand.Rand) []*Creature {
	var b []*Creature = make([]*Creature, 0, 9)
	for i := 1; i < 4; i++ {
		for j := 1; j < 4; j++ {
			b = append(b, &Creature{Color: CreatureColor(i), Species: Species(j), Power: 5, Removed: false})
		}
	}

	random.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return b
}

func GetInitialRandomCreatures(a Alignment, random *rand.Rand) []*Creature {
	creatures := make([]*Creature, 0, 27)
	for i := 0; i < 3; i++ {
		creatures = append(creatures, ShuffledCreatureBag(random)...)
	}

	y := 2
	dy := 2
	x := -1
	dx := -2

	if a == WEST {
		y = 2*MAP_HEIGHT - 2
		dy = -2
		x = MAP_WIDTH
		dx = 2
	}

	for i := 0; i < 27; i++ {
		creatures[i].Alignment = a
		creatures[i].X = x
		creatures[i].Y = y
		x += dx
		y += dy
		// strictly speaking if map width is even,
		// these checks will need to be updated because the
		// west player's row parity will be wrong
		// TODO: Account for possible even width maps
		if y > 2*(MAP_HEIGHT-1) {
			y = 0
		}
		if y < 0 {
			y = 2 * (MAP_HEIGHT - 1)
		}
	}

	return creatures
}
