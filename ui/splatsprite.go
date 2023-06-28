package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/res"
)

type SplatSprite struct {
	X       int
	Y       int
	Img     *ebiten.Image
	Power   int
	Transp  float32
	Removed bool
}

func (c *SplatSprite) MoveTo(x, y int) {
	c.X = x
	c.Y = y
}

func NewSplatSprite(x, y, p int) *SplatSprite {
	r := random.Intn(4) + 1
	img := res.GetImage("splat" + strconv.Itoa(r))

	return &SplatSprite{
		X:       x,
		Y:       y,
		Img:     img,
		Power:   p,
		Transp:  1.0,
		Removed: false,
	}
}

func (c *SplatSprite) Update() {

}

func (c *SplatSprite) Draw(screen *ScaledScreen) {
	if c.Removed {
		return
	}

	opts := &ebiten.DrawImageOptions{}

	opts.ColorScale.Scale(c.Transp, c.Transp, c.Transp, c.Transp)
	opts.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.Img, opts)
	screen.DrawTextCenteredAt(strconv.Itoa(c.Power), 20, c.X+32, c.Y+25, color.White)
}
