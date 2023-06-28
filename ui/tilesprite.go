package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/core"
	"github.com/prizelobby/reverset-raiders/res"
)

type TileSprite struct {
	X        int
	Y        int
	Tile     *core.Tile
	Selected bool
	Reversed bool
}

func NewTileSprite(x, y int, tile *core.Tile) *TileSprite {
	return &TileSprite{
		X:    x,
		Y:    y,
		Tile: tile,
	}
}

func (t *TileSprite) Update() {

}

func (t *TileSprite) Draw(screen *ScaledScreen) {
	opts := &ebiten.DrawImageOptions{}
	left := t.X
	top := t.Y

	img := res.GetImage("hexagon2")
	if t.Selected {
		img = res.GetImage("hexagon-selected")
	}

	toReverese := t.Tile.Reversed != t.Selected
	if toReverese {
		opts.GeoM.Scale(1, -1.0)
		opts.GeoM.Translate(0, float64(img.Bounds().Dy()))
	}

	opts.GeoM.Translate(float64(left), float64(top))

	var obColor color.Color = color.White
	var rvColor color.Color = color.White
	if toReverese {
		obColor = color.RGBA{180, 180, 180, 255}
	} else {
		rvColor = color.RGBA{180, 180, 180, 255}
	}

	screen.DrawImage(img, opts)
	offset := 20
	fontSize := 12.0
	if t.Tile.HasCreature {
		offset = 35
		fontSize = 10.0
	}

	screen.DrawTextCenteredAt(t.Tile.ObverseEffect.String(), fontSize, left+65, top+50-offset, obColor)
	screen.DrawTextCenteredAt(t.Tile.ReverseEffect.String(), fontSize, left+65, top+50+offset, rvColor)
}
