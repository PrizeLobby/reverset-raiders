package ai

import (
	"math"
	"math/rand"

	"github.com/prizelobby/reverset-raiders/core"
)

type Agent struct {
	Game   *core.Game
	Random *rand.Rand
}

func NewAgent(GameSeed int64, RandomSeed int64) *Agent {
	s := rand.NewSource(RandomSeed)
	random := rand.New(s)

	return &Agent{Game: core.NewGameWithSeed(GameSeed), Random: random}
}

func (a *Agent) Reset() {

}

func (a *Agent) MakeMove() (core.GameMove, []core.GameEvent) {

	moves := a.Game.GenerateLegalMoves()
	move := moves[0]
	best := -10000
	for _, m := range moves {
		e := a.Game.AcceptMove(m)
		val := -a.NegaMax(5, -10000, 100000)
		//fmt.Printf("value %d, current best %d\n", val, best)
		if val > best {
			move = m
			best = val
		}
		a.ReverseEvents(e)
	}

	e := a.Game.AcceptMove(move)
	return move, e
}

func (a *Agent) ReverseEvents(events []core.GameEvent) {
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		if e.EventType == core.WARP {
			e.SourceCreature.X = e.SourceX
			e.SourceCreature.Y = e.SourceY
			if e.TargetX >= 0 && e.TargetX < core.MAP_WIDTH {
				a.Game.Map.Tiles[e.TargetX][e.TargetY].HasCreature = false
			}
		} else if e.EventType == core.MOVE {
			e.SourceCreature.X = e.SourceX
			e.SourceCreature.Y = e.SourceY
			if e.SourceX >= 0 && e.SourceX < core.MAP_WIDTH {
				a.Game.Map.Tiles[e.SourceX][e.SourceY].HasCreature = true
				if e.SourceCreature.Removed {
					e.SourceCreature.Removed = false
				}
			}
			if e.TargetX >= 0 && e.TargetX < core.MAP_WIDTH {
				if e.TargetY >= e.TargetX%2 && e.TargetY < 2*core.MAP_HEIGHT-1+(e.TargetX%2) {
					a.Game.Map.Tiles[e.TargetX][e.TargetY].HasCreature = false
				}
			}
		} else if e.EventType == core.DEAL_DAMAGE {
			if e.TargetX == int(core.WEST) {
				a.Game.WestHealth += e.Value
			} else if e.TargetX == int(core.EAST) {
				a.Game.EastHealth += e.Value
			}
		} else if e.EventType == core.UPDATE_POWER {
			e.TargetCreature.Power -= e.Value
		} else if e.EventType == core.DEATH {
			e.SourceCreature.X = e.SourceX
			e.SourceCreature.Y = e.SourceY
			e.SourceCreature.Removed = false
			a.Game.Map.Tiles[e.SourceX][e.SourceY].HasCreature = true
		} else if e.EventType == core.REVERSE_TILE {
			a.Game.Map.Tiles[e.SourceX][e.SourceY].Reversed = !a.Game.Map.Tiles[e.SourceX][e.SourceY].Reversed
		}
	}
	a.Game.CurrentTurn = a.Game.CurrentTurn.Opposite()
}

func (a *Agent) AcceptMove(move core.GameMove) []core.GameEvent {
	return a.Game.AcceptMove(move)
}

func EvalHealth(h int) int {
	return int(math.Sqrt(float64(h * 100)))
}

func (a *Agent) evaluation(c *core.Game, depth int) int {
	multiplier := 1
	if c.CurrentTurn == core.WEST {
		multiplier = -1
	}

	if c.EastHealth <= 0 {
		return multiplier * -(10000 - depth)
	} else if c.WestHealth <= 0 {
		return multiplier * (10000 - depth)
	}

	value := 0
	for _, ec := range a.Game.EastCreatures {
		if !ec.Removed {
			value += Larger(3*ec.Power, 0)
		}
	}
	for _, wc := range a.Game.WestCreatures {
		if !wc.Removed {
			value -= Larger(3*wc.Power, 0)
		}
	}

	return multiplier * (EvalHealth(c.EastHealth) - EvalHealth(c.WestHealth) + value)
}

func (a *Agent) NegaMax(depth int, alpha, beta int) int {

	if depth <= 0 || a.Game.WestHealth <= 0 || a.Game.EastHealth <= 0 {
		return a.evaluation(a.Game, depth)
	}
	moves := a.Game.GenerateLegalMoves()

	best := -99999
	for _, move := range moves {
		e := a.Game.AcceptMove(move)

		val := -a.NegaMax(depth-1, -beta, -alpha)
		if val > best {
			best = val
		}

		if val > alpha {
			alpha = val
		}

		if alpha >= beta {
			a.ReverseEvents(e)
			break
		}
		a.ReverseEvents(e)
	}
	return best

}

func Larger(x, y int) int {
	if x > y {
		return x
	}
	return y
}
