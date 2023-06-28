package scene

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/reverset-raiders/res"
	"github.com/prizelobby/reverset-raiders/ui"
	"github.com/tinne26/etxt"
)

const CENTER = 480
const TITLE_Y_CENTER = 100
const NEW_GAME_Y_CENTER = 300
const CREDITS_Y_CENTER = 400

type MenuScene struct {
	SwitchSceneFunc func(string)
}

func NewMenuScene(switchSceneFunc func(string)) *MenuScene {
	return &MenuScene{
		SwitchSceneFunc: switchSceneFunc,
	}
}

func (m *MenuScene) Update() {
	cursorX, cursorY := ui.AdjustedCursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-300) < 50 {
			m.SwitchSceneFunc("playing")
		}

		if math.Abs(cursorX-CENTER) < 100 && math.Abs(cursorY-400) < 50 {
			m.SwitchSceneFunc("credits")
		}
	}
}

func (m *MenuScene) Draw(scaledScreen *ui.ScaledScreen) {
	scaledScreen.DrawImage(res.GetImage("title"), &ebiten.DrawImageOptions{})
	scaledScreen.DrawTextWithAlign("Reverset Raiders", 48.0, CENTER, TITLE_Y_CENTER, color.Black, etxt.YCenter, etxt.XCenter)
	scaledScreen.DrawTextWithAlign("New Game", 32.0, CENTER, NEW_GAME_Y_CENTER, color.Black, etxt.YCenter, etxt.XCenter)
	scaledScreen.DrawTextWithAlign("Credits", 32.0, CENTER, CREDITS_Y_CENTER, color.Black, etxt.YCenter, etxt.XCenter)
	// NOTE: quit doesn't really make sense for web export

}
