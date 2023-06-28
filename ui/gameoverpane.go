package ui

import (
	"image/color"

	"github.com/prizelobby/reverset-raiders/core"
)

type GameOverPane struct {
	Winner core.Alignment
}

func (g *GameOverPane) Draw(screen *ScaledScreen) {
	screen.DrawRect(0, 0, 960, 480, color.Black)
	screen.DrawTextCenteredAt("Game Over", 48.0, 480, 120, color.White)
	winner := "You win!"
	if g.Winner == core.WEST {
		winner = "You lose"
	}
	screen.DrawTextCenteredAt(winner, 32.0, 480, 300, color.White)
	screen.DrawTextCenteredAt("Return to main", 24.0, 480, 400, color.White)
}
