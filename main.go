package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/core"
	"github.com/prizelobby/reverset-raiders/res"
	"github.com/prizelobby/reverset-raiders/scene"
	"github.com/prizelobby/reverset-raiders/ui"
	"github.com/tinne26/etxt"
)

type GameState int

const (
	MENU GameState = iota
	PLAYING
	CREDITS
)

type EbitenGame struct {
	ScaledScreen *ui.ScaledScreen
	gameState    GameState
	MenuScene    *scene.MenuScene
	CreditsScene *scene.CreditsScene
	GameScene    *scene.GameScene
}

func (g *EbitenGame) SetGameState(s string) {
	if s == "credits" {
		g.gameState = CREDITS
	} else if s == "menu" {
		g.gameState = MENU
	} else if s == "playing" {
		game := core.NewGame()
		g.GameScene = scene.NewGameScene(game, g.SetGameState)
		g.gameState = PLAYING
	}
}

func (g *EbitenGame) Update() error {
	if g.gameState == MENU {
		g.MenuScene.Update()
	} else if g.gameState == CREDITS {
		g.CreditsScene.Update()
	} else if g.gameState == PLAYING {
		g.GameScene.Update()
	}
	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.ScaledScreen.SetTarget(screen)

	//msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	//g.ScaledScreen.DebugPrint(msg)

	if g.gameState == MENU {
		g.MenuScene.Draw(g.ScaledScreen)
	} else if g.gameState == CREDITS {
		g.CreditsScene.Draw(g.ScaledScreen)
	} else if g.gameState == PLAYING {
		g.GameScene.Draw(g.ScaledScreen)
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	panic("use Ebitengine >=v2.5.0")
}

func (g *EbitenGame) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.DeviceScaleFactor()
	canvasWidth := 960 * scale
	canvasHeight := 480 * scale
	return canvasWidth, canvasHeight
}

func main() {
	// create a new text renderer and configure it
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetFont(res.GetFont("Roboto-Medium"))
	txtRenderer.SetAlign(etxt.YCenter, etxt.XCenter)
	txtRenderer.SetSizePx(64)

	scaledScreen := ui.NewScaledScreen(txtRenderer)

	g := &EbitenGame{
		ScaledScreen: scaledScreen,
	}
	g.MenuScene = scene.NewMenuScene(g.SetGameState)
	g.CreditsScene = scene.NewCreditsScene(g.SetGameState)

	ebiten.SetWindowSize(960, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
