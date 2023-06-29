package ai

import (
	"math"
	"math/rand"

	"github.com/prizelobby/reverset-raiders/core"
)

type Agent struct {
	Game   *core.Game
	Random *rand.Rand
	// note: this should probably be tracked in the game but it is too much
	// work to do right now
	TurnsTaken int
}

func NewAgent(GameSeed int64, RandomSeed int64) *Agent {
	s := rand.NewSource(RandomSeed)
	random := rand.New(s)

	return &Agent{Game: core.NewGameWithSeed(GameSeed), Random: random, TurnsTaken: 0}
}

func (a *Agent) Reset() {

}

func (a *Agent) MakeMove() (core.GameMove, []core.GameEvent) {
	//start := time.Now()

	moves := a.Game.GenerateLegalMoves()

	// start at an offset so that if all the evaluations are the same, we choose
	// a random move instead of the first move in the array
	randomOffset := a.Random.Intn(len(moves))
	move := moves[randomOffset]
	best := -10000
	for i := 0; i < len(moves); i++ {
		m := moves[(i+randomOffset)%len(moves)]

		e := a.Game.AcceptMove(m)

		var val int
		// try making the earlier turns take less time
		if a.TurnsTaken < 2 {
			val = -a.SemiNegaMax(4, -10000, 10000)
		} else {
			val = -a.NegaMax(4, -10000, 10000)
		}

		//too slow
		//val := -a.GuidedNegaMax(m, 4, -10000, 10000)

		//too slow
		//val := -a.SemiNegaMax(5, -10000, 10000)

		//fmt.Printf("value %d, current best %d\n", val, best)
		if val > best {
			move = m
			best = val
		}
		a.ReverseEvents(e)
	}

	e := a.Game.AcceptMove(move)

	//duration := time.Since(start)
	//fmt.Println(duration)

	a.TurnsTaken += 1
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
			value += ec.Power
		}
	}
	for _, wc := range a.Game.WestCreatures {
		if !wc.Removed {
			value -= wc.Power
		}
	}

	return multiplier * (EvalHealth(c.EastHealth) - EvalHealth(c.WestHealth) + value)
}

func (a *Agent) PartMoveNegaMax(depth int, alpha, beta int) int {
	if depth <= 0 || a.Game.WestHealth <= 0 || a.Game.EastHealth <= 0 {
		return a.evaluation(a.Game, depth)
	}
	coords := a.Game.AllCoords

	best := -99999
	for _, coord := range coords {
		move := core.GameMove{First: coord, Second: core.MapCoord{X: -1, Y: -1}}
		if !a.Game.IsMoveLocationsEmpty(move) {
			continue
		}

		e := a.Game.AcceptMove(move)

		val := -a.PartMoveNegaMax(depth-1, -beta, -alpha)
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

func (a *Agent) NegaMax(depth int, alpha, beta int) int {
	if depth <= 0 || a.Game.WestHealth <= 0 || a.Game.EastHealth <= 0 {
		return a.evaluation(a.Game, depth)
	}
	moves := a.Game.AllUncheckedMoves
	best := -99999
	for _, move := range moves {
		if !a.Game.IsMoveLocationsEmpty(move) {
			continue
		}
		e := a.Game.AcceptMove(move)
		val := -a.PartMoveNegaMax(depth-1, -beta, -alpha)
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

// this function name isn't really accurate
func (a *Agent) SemiNegaMax(depth int, alpha, beta int) int {
	if depth <= 0 || a.Game.WestHealth <= 0 || a.Game.EastHealth <= 0 {
		return a.evaluation(a.Game, depth)
	}
	coords := a.Game.AllCoords

	best := -99999
	var bestPartMove core.GameMove
	for _, coord := range coords {
		move := core.GameMove{First: coord, Second: core.MapCoord{X: -1, Y: -1}}
		if !a.Game.IsMoveLocationsEmpty(move) {
			continue
		}
		e := a.Game.AcceptMove(move)
		val := -a.PartMoveNegaMax(depth-1, -beta, -alpha)
		if val > best {
			best = val
			bestPartMove = move
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

	best = -99999
	for _, coord := range coords {
		if bestPartMove.First.X == coord.X && bestPartMove.First.Y == coord.Y {
			continue
		}
		move := core.GameMove{First: bestPartMove.First, Second: coord}
		if !a.Game.IsMoveLocationsEmpty(move) {
			continue
		}
		e := a.Game.AcceptMove(move)
		val := -a.PartMoveNegaMax(depth-1, -beta, -alpha)
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

func (a *Agent) GuidedNegaMax(givenMove core.GameMove, depth int, alpha, beta int) int {
	if depth <= 0 || a.Game.WestHealth <= 0 || a.Game.EastHealth <= 0 {
		return a.evaluation(a.Game, depth)
	}

	// theoretically, the best move for the opponent will involve reversing the given move
	// so we should be able to prune off more of the tree
	moves := a.Game.GenerateNextSteps(givenMove.First)
	moves = append(moves, a.Game.GenerateNextSteps(givenMove.Second)...)
	moves = append(moves, a.Game.GenerateLegalMovesWithExclusions(givenMove)...)

	best := -99999
	for _, move := range moves {
		if !a.Game.IsMoveLocationsEmpty(move) {
			continue
		}
		e := a.Game.AcceptMove(move)
		val := -a.GuidedNegaMax(move, depth-1, -beta, -alpha)
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
