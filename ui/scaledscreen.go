package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
	"golang.org/x/image/math/fixed"
)

type ScaledScreen struct {
	Screen         *ebiten.Image
	scaleFactor    float64
	Etxt           *etxt.Renderer
	debugPrintLoc  fixed.Point26_6
	debugPrintSize int
}

func NewScaledScreen(renderer *etxt.Renderer) *ScaledScreen {
	return &ScaledScreen{
		Etxt:           renderer,
		scaleFactor:    ebiten.DeviceScaleFactor(),
		debugPrintSize: int(16 * ebiten.DeviceScaleFactor()),
	}
}

func (s *ScaledScreen) SetTarget(t *ebiten.Image) {
	s.Screen = t
	s.Etxt.SetTarget(t)
	s.debugPrintLoc = fixed.Point26_6{X: 0, Y: 0}
}

func (s *ScaledScreen) DrawImage(image *ebiten.Image, options *ebiten.DrawImageOptions) {
	options.GeoM.Scale(s.scaleFactor, s.scaleFactor)
	s.Screen.DrawImage(image, options)
}

func (s *ScaledScreen) DrawRect(x, y, w, h float64, color color.Color) {
	xx := float32(x * s.scaleFactor)
	yy := float32(y * s.scaleFactor)
	hh := float32(h * s.scaleFactor)
	ww := float32(w * s.scaleFactor)

	vector.DrawFilledRect(s.Screen, xx, yy, ww, hh, color, false)
}

func (s *ScaledScreen) DrawRectShader(w, h int, shader *ebiten.Shader, opts *ebiten.DrawRectShaderOptions) {
	ww := int(float64(w) * s.scaleFactor)
	hh := int(float64(h) * s.scaleFactor)

	opts.GeoM.Scale(s.scaleFactor, s.scaleFactor)
	s.Screen.DrawRectShader(ww, hh, shader, opts)
}

func (s *ScaledScreen) scaledTextSize(size float64) int {
	return int(size * s.scaleFactor)
}

func (s *ScaledScreen) DrawText(t string, size float64, x, y int, color color.Color) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSizePx(s.scaledTextSize(size))
	s.Etxt.SetAlign(etxt.Top, etxt.Left)
	s.Etxt.Draw(t, xx, yy)
}

func (s *ScaledScreen) DrawTextCenteredAt(t string, size float64, x, y int, color color.Color) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSizePx(s.scaledTextSize(size))
	s.Etxt.SetAlign(etxt.YCenter, etxt.XCenter)
	s.Etxt.Draw(t, xx, yy)
}

func (s *ScaledScreen) DrawTextWithAlign(t string, size float64, x, y int, color color.Color, vAlign etxt.VertAlign, hAlign etxt.HorzAlign) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSizePx(s.scaledTextSize(size))
	s.Etxt.SetAlign(vAlign, hAlign)
	s.Etxt.Draw(t, xx, yy)
}

func (s *ScaledScreen) DebugPrint(str string) {
	s.Etxt.SetSizePx(s.debugPrintSize)
	s.Etxt.SetAlign(etxt.Top, etxt.Left)
	s.Etxt.SetColor(color.White)
	s.debugPrintLoc = s.Etxt.Draw(str+"\n", s.debugPrintLoc.X.Ceil(), s.debugPrintLoc.Y.Ceil())
}
