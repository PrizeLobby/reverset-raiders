package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/prizelobby/reverset-raiders/core"
	"github.com/prizelobby/reverset-raiders/res"
)

type CreatureSprite struct {
	X       int
	Y       int
	Facing  core.Alignment
	Img     *ebiten.Image
	Power   int
	Rot     int
	Transp  float32
	Removed bool
}

func (c *CreatureSprite) MoveTo(x, y int) {
	c.X = x
	c.Y = y
}

func ColorToSuffix(c core.CreatureColor) string {
	if c == core.Blue {
		return "_blue"
	} else if c == core.Green {
		return "_green"
	} else if c == core.Red {
		return "_red"
	}
	return ""
}

func NewCreatureSprite(x, y int, c *core.Creature) *CreatureSprite {
	var img *ebiten.Image
	if c.Species == core.Capybara {
		img = res.GetImage("capybarasm" + ColorToSuffix(c.Color))
	} else if c.Species == core.Duck {
		img = res.GetImage("ducksm" + ColorToSuffix(c.Color))
	} else if c.Species == core.Tortoise {
		img = res.GetImage("turtlesm" + ColorToSuffix(c.Color))
	}

	return &CreatureSprite{
		X:       x,
		Y:       y,
		Facing:  c.Alignment,
		Img:     img,
		Power:   c.Power,
		Transp:  1.0,
		Removed: false,
	}
}

func (c *CreatureSprite) Update() {

}

func (c *CreatureSprite) Draw(screen *ScaledScreen) {
	if c.Removed {
		return
	}

	opts := &ebiten.DrawImageOptions{}

	if c.Facing == core.WEST {
		opts.GeoM.Scale(-1.0, 1)
		opts.GeoM.Translate(float64(c.Img.Bounds().Dx()), 0)
	}
	opts.ColorScale.Scale(c.Transp, c.Transp, c.Transp, c.Transp)

	opts.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.Img, opts)
	screen.DrawTextCenteredAt(strconv.Itoa(c.Power), 10, c.X+32, c.Y+35, color.White)
}
