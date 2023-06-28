package core

import (
	"math/rand"
	"time"
)

type GameMove struct {
	First  MapCoord
	Second MapCoord
}

type Game struct {
	Map           *Map
	EastCreatures []*Creature
	WestCreatures []*Creature
	EastHealth    int
	WestHealth    int
	CurrentTurn   Alignment
	Rand          *rand.Rand
	Seed          int64
	AllCoords     []MapCoord
}

func NewGameWithSeed(seed int64) *Game {
	s := rand.NewSource(seed)
	random := rand.New(s)

	allCoords := make([]MapCoord, 0, MAP_HEIGHT*MAP_WIDTH)
	for i := 0; i < MAP_WIDTH; i++ {
		for j := i % 2; j < 2*MAP_HEIGHT; j += 2 {
			allCoords = append(allCoords, MapCoord{i, j})
		}
	}

	return &Game{
		Map:           NewMap(random),
		EastCreatures: GetInitialRandomCreatures(EAST, random),
		WestCreatures: GetInitialRandomCreatures(WEST, random),
		EastHealth:    50,
		WestHealth:    50,
		CurrentTurn:   EAST,
		Rand:          random,
		Seed:          seed,
		AllCoords:     allCoords,
	}
}

func NewGame() *Game {
	return NewGameWithSeed(time.Now().UnixNano())
}

type GameEventType int

const (
	MOVE GameEventType = iota
	WARP
	APPLY_EFFECT
	UPDATE_POWER
	DEAL_DAMAGE
	REVERSE_TILE
	DEATH
	GAME_OVER
)

type GameEvent struct {
	EventType      GameEventType
	SourceX        int
	SourceY        int
	SourceCreature *Creature
	TargetX        int
	TargetY        int
	TargetCreature *Creature
	Value          int
	Effect         Effect
}

func (g *Game) AcceptMove(move GameMove) []GameEvent {
	events := make([]GameEvent, 0, 100)

	// we don't actually check if the move is valid (ie at least 1 subaction must be valid)
	for _, m := range []MapCoord{move.First, move.Second} {
		if m.X != -1 && m.Y != -1 {
			g.Map.Tiles[m.X][m.Y].Reversed = !g.Map.Tiles[m.X][m.Y].Reversed
			events = append(events, GameEvent{
				EventType: REVERSE_TILE,
				SourceX:   m.X,
				SourceY:   m.Y,
			})
		}
	}

	if g.CurrentTurn == EAST {
		for i := 0; i < len(g.EastCreatures); i++ {
			creature := g.EastCreatures[i]
			if creature.Removed {
				continue
			}
			creatureStartX := creature.X
			creatureStartY := creature.Y
			if creature.X < 0 || creature.X >= MAP_WIDTH {
				creature.X += 1

				events = append(events, GameEvent{
					EventType:      MOVE,
					SourceX:        creatureStartX,
					SourceY:        creatureStartY,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.Y,
				})

				// if we're still off the map, we don't need to do any other calculations
				if creature.X < 0 || creature.X >= MAP_WIDTH {
					continue
				}
				// otherwise we need to check enemies, etc
			} else {
				tile := g.Map.Tiles[creature.X][creature.Y]

				if tile.Reversed {
					creature.Y += 1
				} else {
					creature.Y -= 1
				}
				creature.X += 1

				tile.Reversed = !tile.Reversed
				tile.HasCreature = false

				events = append(events, GameEvent{
					EventType:      MOVE,
					SourceX:        creatureStartX,
					SourceY:        creatureStartY,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.Y,
				})
				events = append(events, GameEvent{
					EventType: REVERSE_TILE,
					SourceX:   creatureStartX,
					SourceY:   creatureStartY,
				})
			}

			if creature.Y < creature.X%2 {
				events = append(events, GameEvent{
					EventType:      WARP,
					SourceX:        creature.X,
					SourceY:        creature.Y,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        2*MAP_HEIGHT - 2 + (creature.X % 2),
				})
				creature.Y = 2*MAP_HEIGHT - 2 + (creature.X % 2)
			}
			if creature.Y > 2*MAP_HEIGHT-2+(creature.X%2) {
				events = append(events, GameEvent{
					EventType:      WARP,
					SourceX:        creature.X,
					SourceY:        creature.Y,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.X % 2,
				})
				creature.Y = creature.X % 2
			}

			// if we move off the map into opponent territory
			if creature.X >= MAP_WIDTH {
				creature.Removed = true
				g.WestHealth -= creature.Power

				events = append(events, GameEvent{
					EventType:      DEAL_DAMAGE,
					SourceCreature: creature,
					TargetX:        int(WEST),
					Value:          creature.Power,
				})
				if g.WestHealth <= 0 {
					events = append(events, GameEvent{
						EventType: GAME_OVER,
						SourceX:   int(EAST),
						TargetX:   int(WEST),
					})
				}

				continue
			}

			for _, c := range g.WestCreatures {
				if c.Removed {
					continue
				}
				if c.X == creature.X && c.Y == creature.Y {
					if c.Power == creature.Power {
						g.Map.Tiles[creature.X][creature.Y].HasCreature = false
						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        c.X,
							SourceY:        c.Y,
							SourceCreature: c,
						})
						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        creature.X,
							SourceY:        creature.Y,
							SourceCreature: creature,
						})
						c.X = -1000
						c.Removed = true
						creature.X = 1000
						creature.Removed = true

					} else if c.Power > creature.Power {
						c.Power -= creature.Power
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: c,
							Value:          -creature.Power,
						})

						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        creature.X,
							SourceY:        creature.Y,
							SourceCreature: creature,
						})
						creature.X = 1000
						creature.Removed = true
					} else if creature.Power > c.Power {
						creature.Power -= c.Power
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: creature,
							Value:          -c.Power,
						})

						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        c.X,
							SourceY:        c.Y,
							SourceCreature: c,
						})
						c.X = -1000
						c.Removed = true
					}
					break
				}
			}

			if !creature.Removed {
				e := g.Map.Tiles[creature.X][creature.Y].GetActiveEffect()
				for _, cc := range g.EastCreatures {
					if cc.ApplyEffect(e) {
						events = append(events, GameEvent{
							EventType:      APPLY_EFFECT,
							TargetCreature: cc,
							Effect:         e,
						})
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: cc,
							Value:          e.Value,
						})
					}
				}
				g.Map.Tiles[creature.X][creature.Y].HasCreature = true
			}
		}
		g.CurrentTurn = WEST
	} else {
		for i := 0; i < len(g.WestCreatures); i++ {
			creature := g.WestCreatures[i]
			if creature.Removed {
				continue
			}
			creatureStartX := creature.X
			creatureStartY := creature.Y
			if creature.X < 0 || creature.X >= MAP_WIDTH {
				creature.X -= 1
				events = append(events, GameEvent{
					EventType:      MOVE,
					SourceX:        creatureStartX,
					SourceY:        creatureStartY,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.Y,
				})

				// if we're still off the map, we don't need to do any other calculations
				if creature.X < 0 || creature.X >= MAP_WIDTH {
					continue
				}
				// otherwise we need to check enemies, etc

			} else {
				tile := g.Map.Tiles[creature.X][creature.Y]

				if tile.Reversed {
					creature.Y += 1
				} else {
					creature.Y -= 1
				}
				creature.X -= 1
				tile.Reversed = !tile.Reversed
				tile.HasCreature = false

				events = append(events, GameEvent{
					EventType:      MOVE,
					SourceX:        creatureStartX,
					SourceY:        creatureStartY,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.Y,
				})
				events = append(events, GameEvent{
					EventType: REVERSE_TILE,
					SourceX:   creatureStartX,
					SourceY:   creatureStartY,
				})
			}

			if creature.Y < creature.X%2 {
				events = append(events, GameEvent{
					EventType:      WARP,
					SourceX:        creature.X,
					SourceY:        creature.Y,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        2*MAP_HEIGHT - 2 + (creature.X % 2),
				})
				creature.Y = 2*MAP_HEIGHT - 2 + (creature.X % 2)
			}
			if creature.Y > 2*MAP_HEIGHT-2+(creature.X%2) {
				events = append(events, GameEvent{
					EventType:      WARP,
					SourceX:        creature.X,
					SourceY:        creature.Y,
					SourceCreature: creature,
					TargetX:        creature.X,
					TargetY:        creature.X % 2,
				})
				creature.Y = creature.X % 2
			}

			// if we move off the map into opponent territory
			if creature.X < 0 {
				creature.Removed = true
				g.EastHealth -= creature.Power
				events = append(events, GameEvent{
					EventType:      DEAL_DAMAGE,
					SourceCreature: creature,
					TargetX:        int(EAST),
					Value:          creature.Power,
				})
				if g.EastHealth <= 0 {
					events = append(events, GameEvent{
						EventType: GAME_OVER,
						SourceX:   int(WEST),
						TargetX:   int(EAST),
					})
				}
				continue
			}

			for _, c := range g.EastCreatures {
				if c.Removed {
					continue
				}
				if c.X == creature.X && c.Y == creature.Y {
					if c.Power == creature.Power {
						g.Map.Tiles[creature.X][creature.Y].HasCreature = false
						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        c.X,
							SourceY:        c.Y,
							SourceCreature: c,
						})
						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        creature.X,
							SourceY:        creature.Y,
							SourceCreature: creature,
						})

						c.X = 1000
						c.Removed = true
						creature.X = -1000
						creature.Removed = true
					} else if c.Power > creature.Power {
						c.Power -= creature.Power
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: c,
							Value:          -creature.Power,
						})

						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        creature.X,
							SourceY:        creature.Y,
							SourceCreature: creature,
						})
						creature.X = -1000
						creature.Removed = true
					} else if creature.Power > c.Power {
						creature.Power -= c.Power
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: creature,
							Value:          -c.Power,
						})

						events = append(events, GameEvent{
							EventType:      DEATH,
							SourceX:        c.X,
							SourceY:        c.Y,
							SourceCreature: c,
						})
						c.X = 1000
						c.Removed = true
					}
					break
				}
			}

			if !creature.Removed {
				e := g.Map.Tiles[creature.X][creature.Y].GetActiveEffect()
				for _, cc := range g.WestCreatures {
					if cc.ApplyEffect(e) {
						events = append(events, GameEvent{
							EventType:      APPLY_EFFECT,
							TargetCreature: cc,
							Effect:         e,
						})
						events = append(events, GameEvent{
							EventType:      UPDATE_POWER,
							TargetCreature: cc,
							Value:          e.Value,
						})
					}
				}
				g.Map.Tiles[creature.X][creature.Y].HasCreature = true
			}
		}
		g.CurrentTurn = EAST
	}
	return events
}

func (g *Game) GenerateLegalMoves() []GameMove {
	moves := make([]GameMove, 0)
	for i := 0; i < len(g.AllCoords); i++ {
		ci := g.AllCoords[i]
		if g.Map.Tiles[ci.X][ci.Y].HasCreature {
			continue
		}
		/*
			for j := i + 1; j < len(g.AllCoords); j++ {
				cj := g.AllCoords[j]
				if g.Map.Tiles[cj.X][cj.Y].HasCreature {
					continue
				}
				moves = append(moves, GameMove{
					First:  ci,
					Second: cj,
				})
			}*/
		moves = append(moves, GameMove{
			First:  ci,
			Second: MapCoord{-1, -1},
		})
	}

	return moves
}
