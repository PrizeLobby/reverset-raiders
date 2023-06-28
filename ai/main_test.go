package ai_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/prizelobby/reverset-raiders/ai"
	"github.com/prizelobby/reverset-raiders/core"
)

func TestReverseGame(t *testing.T) {
	game := core.NewGameWithSeed(1)
	agent := ai.NewAgent(1, 0)
	if err := GamesAreEqual(game, agent.Game); err != nil {
		t.Fatalf(err.Error())
	}

	_, e := agent.MakeMove()
	agent.ReverseEvents(e)
	if err := GamesAreEqual(game, agent.Game); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestReverseGame2(t *testing.T) {
	s := rand.NewSource(1)
	random := rand.New(s)

	game := core.NewGameWithSeed(2)
	agent := ai.NewAgent(2, 2)

	for i := 0; i < 10; i++ {
		legalMoves := game.GenerateLegalMoves()
		r := random.Intn(len(legalMoves))
		game.AcceptMove(legalMoves[r])
		agent.AcceptMove(legalMoves[r])
	}
	if err := GamesAreEqual(game, agent.Game); err != nil {
		t.Fatalf(err.Error())
	}

	_, e := agent.MakeMove()
	agent.ReverseEvents(e)
	if err := GamesAreEqual(game, agent.Game); err != nil {
		t.Fatalf(err.Error())
	}
}

func GamesAreEqual(g1, g2 *core.Game) error {
	if g1.CurrentTurn != g2.CurrentTurn {
		return errors.New("Wrong player turn")
	}

	if g1.WestHealth != g2.WestHealth {
		return errors.New("West Health different")
	}
	if g1.EastHealth != g2.EastHealth {
		return errors.New("East Health different")
	}
	for _, c := range g1.AllCoords {
		t1 := g1.Map.Tiles[c.X][c.Y]
		t2 := g2.Map.Tiles[c.X][c.Y]
		if t1.Reversed != t2.Reversed {
			return fmt.Errorf("tile at %d %d has different values for Reversed", c.X, c.Y)
		}
		if t1.HasCreature != t2.HasCreature {
			return fmt.Errorf("tile at %d %d has different values for HasCreature", c.X, c.Y)
		}
		if !EffectsAreEqual(t1.ReverseEffect, t2.ReverseEffect) {
			return fmt.Errorf("tile at %d %d has different values for ReverseEffect", c.X, c.Y)
		}

		if !EffectsAreEqual(t1.ObverseEffect, t2.ObverseEffect) {
			return fmt.Errorf("tile at %d %d has different values for ObverseEffect", c.X, c.Y)
		}
	}
	for i := 0; i < len(g1.EastCreatures); i++ {
		c1 := g1.EastCreatures[i]
		c2 := g2.EastCreatures[i]
		if c1.X != c2.X {
			return fmt.Errorf("east creature at index %d has different x values %d %d", i, c1.X, c2.X)
		}
		if c1.Y != c2.Y {
			return fmt.Errorf("east creature at index %d has different y values %d %d", i, c1.Y, c2.Y)
		}
		if c1.Removed != c2.Removed {
			return fmt.Errorf("east creature at index %d has different Removed values %t %t", i, c1.Removed, c2.Removed)
		}
		if c1.Power != c2.Power {
			return fmt.Errorf("east creature at index %d has different power values %d %d", i, c1.Power, c2.Power)
		}
	}
	for i := 0; i < len(g1.WestCreatures); i++ {
		c1 := g1.WestCreatures[i]
		c2 := g2.WestCreatures[i]
		if c1.X != c2.X {
			return fmt.Errorf("west creature at index %d has different x values %d %d", i, c1.X, c2.X)
		}
		if c1.Y != c2.Y {
			return fmt.Errorf("west creature at index %d has different y values %d %d", i, c1.Y, c2.Y)
		}
		if c1.Removed != c2.Removed {
			return fmt.Errorf("west creature at index %d has different Removed values %t %t", i, c1.Removed, c2.Removed)
		}
		if c1.Power != c2.Power {
			return fmt.Errorf("west creature at index %d has different power values %d %d", i, c1.Power, c2.Power)
		}
	}

	return nil
}

func EffectsAreEqual(e1, e2 core.Effect) bool {
	return e1.X == e2.X && e1.Y == e2.Y && e1.Targets == e2.Targets && e1.SpeciesCondition == e2.SpeciesCondition && e1.ColorCondition == e2.ColorCondition && e1.Value == e2.Value
}
