package scene

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/reverset-raiders/ai"
	"github.com/prizelobby/reverset-raiders/core"
	"github.com/prizelobby/reverset-raiders/res"
	"github.com/prizelobby/reverset-raiders/ui"
	"github.com/prizelobby/reverset-raiders/ui/animation"
)

const MAP_START_X = 211
const MAP_START_Y = 50
const MAP_START_X_FLOAT = 211.0
const MAP_START_Y_FLOAT = 50.0

const CONFIRM_BUTTON_X_FLOAT = 825.0
const CONFIRM_BUTTON_Y_FLOAT = 400.0

const EAST_HEALTH_X = 100
const EAST_HEALTH_Y = 360

const WEST_HEALTH_X = 860
const WEST_HEALTH_Y = 360

const HELP_TEXT_X_CENTER = 480
const HELP_TEXT_Y_CENTER = 440

type GameUIState int

const (
	WAITING_FOR_PLAYER_MOVE GameUIState = iota
	WAITING_FOR_OPP_MOVE
	WAITING_FOR_PLAYER_ANIMIMATION
	GAME_OVER
)

type GameScene struct {
	//hoverTileX      int
	//hoverTileY      int
	selectedCoords    []core.MapCoord
	Game              *core.Game
	SwitchSceneFunc   func(string)
	UIState           GameUIState
	EventsToAnimate   []core.GameEvent
	OngoingAnimation  animation.Anim
	CreatureSprites   []*ui.CreatureSprite
	EffectSprites     []*ui.EffectSprite
	TileSprites       [][]*ui.TileSprite
	SplatSprite       *ui.SplatSprite
	CreatureSpriteMap map[*core.Creature]*ui.CreatureSprite
	GameOverPane      *ui.GameOverPane
	MoveChan          chan core.GameMove
	Agent             *ai.Agent
}

func NewGameScene(game *core.Game, f func(string)) *GameScene {
	tileSprites := make([][]*ui.TileSprite, core.MAP_WIDTH)
	for i := 0; i < core.MAP_WIDTH; i++ {
		j := i % 2
		tileSprites[i] = make([]*ui.TileSprite, 2*core.MAP_HEIGHT-1+j)
		for ; j < 2*core.MAP_HEIGHT; j += 2 {
			x, y := HexIndicesToScreenCoord(i, j)
			tileSprites[i][j] = ui.NewTileSprite(x, y, game.Map.Tiles[i][j])
		}
	}

	creatureSprites := make([]*ui.CreatureSprite, 0)
	creatureMap := make(map[*core.Creature]*ui.CreatureSprite)

	ec := game.EastCreatures[0]
	x, y := HexIndicesToScreenCoord(ec.X, ec.Y)
	s := ui.NewCreatureSprite(x+32, y+15, ec)
	creatureSprites = append(creatureSprites, s)
	creatureMap[ec] = s

	wc := game.WestCreatures[0]
	x, y = HexIndicesToScreenCoord(wc.X, wc.Y)
	s2 := ui.NewCreatureSprite(x+32, y+15, wc)
	creatureSprites = append(creatureSprites, s2)
	creatureMap[wc] = s2

	return &GameScene{
		Game:              game,
		SwitchSceneFunc:   f,
		selectedCoords:    make([]core.MapCoord, 0, 3),
		TileSprites:       tileSprites,
		CreatureSprites:   creatureSprites,
		CreatureSpriteMap: creatureMap,
		GameOverPane:      &ui.GameOverPane{},
		MoveChan:          make(chan core.GameMove),
		Agent:             ai.NewAgent(game.Seed, 1),
	}
}

func (g *GameScene) MouseCoordsToTileCoords(x, y float64) (int, int) {
	x -= MAP_START_X_FLOAT
	y -= MAP_START_Y

	xp := math.Floor(x / 100)
	remain := x - (xp * 100)
	xIndex := int(xp)

	var sx, sy int
	if remain > 28 {
		sx = xIndex
		if sx%2 == 0 {
			sy = 2 * int(math.Floor(y/100))
		} else {
			y -= 50
			sy = 2*int(math.Floor(y/100)) + 1
		}
	} else {
		yp := math.Floor(y / 50)
		yInd := int(yp)
		yRemain := y - (yp * 50)

		if xIndex%2 == 1 {
			yInd -= 1
		}

		if yInd%2 == 0 {
			if yRemain < 2*(28-remain) {
				sx = xIndex - 1
				sy = yInd - 1
			} else {
				sx = xIndex
				sy = yInd
			}
		} else {
			if yRemain < 2*remain {
				sx = xIndex
				sy = yInd - 1
			} else {
				sx = xIndex - 1
				sy = yInd
			}
		}

		if xIndex%2 == 1 {
			sy += 1
		}
	}
	return sx, sy
}

func (g *GameScene) IndexOfSelectedCoord(i, j int) int {
	for ind, c := range g.selectedCoords {
		if c.X == i && c.Y == j {
			return ind
		}
	}
	return -1
}

func (g *GameScene) IsValidTileCoordsForTurn(i, j int) bool {
	if i < 0 || i >= core.MAP_WIDTH {
		return false
	}
	if j < 0 || j >= len(g.Game.Map.Tiles[i]) {
		return false
	}
	return !g.Game.Map.Tiles[i][j].HasCreature
}

func (g *GameScene) IsInsideConfirmButton(x, y float64) bool {
	return x > CONFIRM_BUTTON_X_FLOAT &&
		x < CONFIRM_BUTTON_X_FLOAT+80 &&
		y > CONFIRM_BUTTON_Y_FLOAT &&
		y < CONFIRM_BUTTON_Y_FLOAT+60
}

func HexIndicesToScreenCoord(i, j int) (int, int) {
	return 102*i + MAP_START_X, 50*j + MAP_START_Y
}

func (g *GameScene) AnimationForEvent(e core.GameEvent) animation.Anim {
	if e.EventType == core.MOVE {
		//fmt.Printf("process event type move %d %d %d %d %s\n", e.SourceX, e.SourceY, e.TargetX, e.TargetY, e.SourceCreature)

		// having a side effect in this function isn't really great
		// TODO: figure out a better way to do this
		if (e.TargetX == -1 && e.SourceCreature.Alignment == core.EAST) || (e.TargetX == core.MAP_WIDTH && e.SourceCreature.Alignment == core.WEST) {
			if g.CreatureSpriteMap[e.SourceCreature] == nil {
				x, y := HexIndicesToScreenCoord(e.TargetX, e.TargetY)
				s := ui.NewCreatureSprite(x+32, y+15, e.SourceCreature)
				g.CreatureSprites = append(g.CreatureSprites, s)
				g.CreatureSpriteMap[e.SourceCreature] = s
			}
			return nil
		} else if e.TargetX >= 0 && e.TargetX <= core.MAP_WIDTH {
			sx, sy := HexIndicesToScreenCoord(e.SourceX, e.SourceY)
			tx, ty := HexIndicesToScreenCoord(e.TargetX, e.TargetY)
			return animation.NewSpriteMovement(sx+32, sy+15, tx+32, ty+15, 10, g.CreatureSpriteMap[e.SourceCreature])
		}
	} else if e.EventType == core.WARP {
		sx, sy := HexIndicesToScreenCoord(e.SourceX, e.SourceY)
		tx, ty := HexIndicesToScreenCoord(e.TargetX, e.TargetY)

		if e.TargetX < 0 || e.TargetX >= core.MAP_WIDTH {
			return nil
		}

		return animation.NewSpriteMovement(sx+32, sy+15, tx+32, ty+15, 1, g.CreatureSpriteMap[e.SourceCreature])
	} else if e.EventType == core.DEATH {
		return animation.NewDeathAnimation(g.CreatureSpriteMap[e.SourceCreature])
	} else if e.EventType == core.UPDATE_POWER {
		return animation.NewCreatureSpriteUpdatePower(g.CreatureSpriteMap[e.TargetCreature], e.Value)
	} else if e.EventType == core.DEAL_DAMAGE {
		x := EAST_HEALTH_X - 29
		if e.TargetX == int(core.WEST) {
			x = WEST_HEALTH_X - 29
		}
		g.SplatSprite = ui.NewSplatSprite(x, WEST_HEALTH_Y+3, e.Value)
		g.CreatureSpriteMap[e.SourceCreature].Removed = true
		return animation.NewSplatAnimation(g.SplatSprite)
	} else if e.EventType == core.GAME_OVER {
		g.GameOverPane.Winner = core.Alignment(e.SourceX)
		g.UIState = GAME_OVER
	} else if e.EventType == core.APPLY_EFFECT {
		g.EffectSprites = make([]*ui.EffectSprite, 0)
		x, y := HexIndicesToScreenCoord(e.TargetCreature.X, e.TargetCreature.Y)
		s := ui.NewEffectSprite(x+32, y+15)
		g.EffectSprites = append(g.EffectSprites, s)
		return animation.NewEffectSpriteAnimation(s)
	}

	return nil
}

func (g *GameScene) UpdateAnimations() {
	if g.OngoingAnimation != nil {
		if g.OngoingAnimation.IsFinished() {
			g.OngoingAnimation = nil
		} else {
			g.OngoingAnimation.Update()
		}
		return
	}

	if len(g.EventsToAnimate) == 0 {
		if g.Game.CurrentTurn == core.EAST {
			g.UIState = WAITING_FOR_PLAYER_MOVE
		} else {
			g.UIState = WAITING_FOR_OPP_MOVE

			// TODO: this should be moved to immediately after the player clicks the "confirm move"
			// button so we can do the calculations concurrently with the animations
			go func() {
				move, _ := g.Agent.MakeMove()
				g.MoveChan <- move
			}()
		}

	} else {
		event := g.EventsToAnimate[0]
		g.EventsToAnimate = g.EventsToAnimate[1:]
		g.OngoingAnimation = g.AnimationForEvent(event)
	}
}

func (g *GameScene) UpdatePlayerActions() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ui.AdjustedCursorPosition()
		i, j := g.MouseCoordsToTileCoords(cx, cy)
		if index := g.IndexOfSelectedCoord(i, j); index != -1 {
			g.selectedCoords[index] = g.selectedCoords[len(g.selectedCoords)-1]
			g.selectedCoords = g.selectedCoords[:len(g.selectedCoords)-1]
			g.TileSprites[i][j].Selected = !g.TileSprites[i][j].Selected
		} else if g.IsValidTileCoordsForTurn(i, j) && len(g.selectedCoords) < 2 {
			g.selectedCoords = append(g.selectedCoords, core.MapCoord{X: i, Y: j})
			g.TileSprites[i][j].Selected = !g.TileSprites[i][j].Selected
		} else if g.IsInsideConfirmButton(cx, cy) {
			if len(g.selectedCoords) > 0 {
				var move core.GameMove
				move.First = g.selectedCoords[0]
				g.TileSprites[move.First.X][move.First.Y].Selected = false
				if len(g.selectedCoords) > 1 {
					move.Second = g.selectedCoords[1]
					g.TileSprites[move.Second.X][move.Second.Y].Selected = false
				} else {
					move.Second = core.MapCoord{X: -1, Y: -1}
				}
				g.Agent.AcceptMove(move)
				g.EventsToAnimate = g.Game.AcceptMove(move)
				g.selectedCoords = make([]core.MapCoord, 0, 2)
				g.UIState = WAITING_FOR_PLAYER_ANIMIMATION
			}
		}
	}
}

func (g *GameScene) UpdateGameOverActions() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ui.AdjustedCursorPosition()
		if math.Abs(cx-480) < 100 && math.Abs(cy-400) < 50 {
			g.SwitchSceneFunc("menu")
		}
	}
}

func (g *GameScene) Update() {
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			for i := 0; i < core.MAP_WIDTH; i++ {
				for j := i % 2; j < 2*core.MAP_HEIGHT; j += 2 {
					fmt.Printf("(%d,%d) reversed %t hasCreature %t\n", i, j, g.Game.Map.Tiles[i][j].Reversed, g.Game.Map.Tiles[i][j].HasCreature)
				}
			}
		}*/

	// TODO: this block should move into the waiting for animation and we should
	// start the goroutine while the player's animations are playing
	select {
	case m := <-g.MoveChan:
		t1 := g.TileSprites[m.First.X][m.First.Y]
		var t2 *ui.TileSprite = nil
		if m.Second.X > 0 {
			t2 = g.TileSprites[m.Second.X][m.Second.Y]
		}

		g.OngoingAnimation = animation.NewTileHighlightAnimation(t1, t2)

		g.EventsToAnimate = g.Game.AcceptMove(m)
		g.UIState = WAITING_FOR_PLAYER_ANIMIMATION
	default:
	}

	if g.UIState == WAITING_FOR_PLAYER_ANIMIMATION {
		g.UpdateAnimations()
	} else if g.UIState == WAITING_FOR_PLAYER_MOVE {
		g.UpdatePlayerActions()
	} else if g.UIState == WAITING_FOR_OPP_MOVE {

	} else if g.UIState == GAME_OVER {
		g.UpdateGameOverActions()
	}
}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	for i := 0; i < core.MAP_WIDTH; i++ {
		for j := i % 2; j < 2*core.MAP_HEIGHT; j += 2 {
			g.TileSprites[i][j].Draw(screen)
		}
	}
	for _, cs := range g.CreatureSprites {
		cs.Draw(screen)
	}

	for _, es := range g.EffectSprites {
		es.Draw(screen)
	}

	cOpts := &ebiten.DrawImageOptions{}
	cOpts.GeoM.Translate(CONFIRM_BUTTON_X_FLOAT, CONFIRM_BUTTON_Y_FLOAT)
	if len(g.selectedCoords) > 0 {
		screen.DrawImage(res.GetImage("confirmmove"), cOpts)
	} else {
		screen.DrawImage(res.GetImage("confirmmoveinactive"), cOpts)
	}

	screen.DrawTextCenteredAt("Player Health", 16, EAST_HEALTH_X, EAST_HEALTH_Y-28, color.RGBA{0xac, 0x32, 0x32, 0xff})
	screen.DrawTextCenteredAt("Enemy Health", 16, WEST_HEALTH_X, WEST_HEALTH_Y-28, color.RGBA{0xac, 0x32, 0x32, 0xff})
	screen.DrawTextCenteredAt(strconv.Itoa(g.Game.EastHealth), 32, EAST_HEALTH_X, EAST_HEALTH_Y, color.RGBA{0xac, 0x32, 0x32, 0xff})
	screen.DrawTextCenteredAt(strconv.Itoa(g.Game.WestHealth), 32, WEST_HEALTH_X, WEST_HEALTH_Y, color.RGBA{0xac, 0x32, 0x32, 0xff})

	if g.UIState == WAITING_FOR_PLAYER_MOVE {
		screen.DrawTextCenteredAt("Select up to 2 tiles to reverse.", 28, HELP_TEXT_X_CENTER, HELP_TEXT_Y_CENTER, color.White)
	} else {
		screen.DrawTextCenteredAt("Waiting for opponent...", 28, HELP_TEXT_X_CENTER, HELP_TEXT_Y_CENTER, color.White)
	}

	eastReserves := "Reserves - Row\n"
	eCount := 0
	for _, ec := range g.Game.EastCreatures {
		if !ec.Removed {
			if ec.X < 0 {
				eCount += 1
				eastReserves += ec.Color.String() + " " + ec.Species.String() + " - " + strconv.Itoa(1+(ec.Y/2)) + "\n"

				if eCount > 10 {
					break
				}
			}
		}
	}
	screen.DrawText(eastReserves, 12, 20, 100, color.White)

	westReserves := "Reserves - Row\n"
	wCount := 0
	for _, wc := range g.Game.WestCreatures {
		if !wc.Removed {
			if wc.X >= core.MAP_WIDTH {
				wCount += 1
				westReserves += wc.Color.String() + " " + wc.Species.String() + " - " + strconv.Itoa(1+(wc.Y/2)) + "\n"

				if wCount > 10 {
					break
				}
			}
		}
	}
	screen.DrawText(westReserves, 12, 825, 100, color.White)

	if g.SplatSprite != nil {
		g.SplatSprite.Draw(screen)
	}

	if g.UIState == GAME_OVER {
		g.GameOverPane.Draw(screen)
	}
}

var RED = []float64{0xac / 255.0, 0x32 / 255.0, 0x32 / 255.0, 1.0}
var GREEN = []float64{0x6a / 255.0, 0xbe / 255.0, 0x30 / 255.0, 1.0}
var BLUE = []float64{0x30 / 255.0, 0x60 / 255.0, 0x82 / 255.0, 1.0}
