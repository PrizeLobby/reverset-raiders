package ui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/res"
)

type EffectSprite struct {
	x int
	y int
	I int
}

func NewEffectSprite(x, y int) *EffectSprite {
	return &EffectSprite{x: x, y: y, I: 1}
}

func (s *EffectSprite) Draw(screen *ScaledScreen) {
	if s.I < 1 || s.I > 15 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	screen.DrawImage(res.GetImage("buffeffect"+strconv.Itoa(s.I)), op)
}
