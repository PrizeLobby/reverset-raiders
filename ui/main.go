package ui

import (
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func AdjustedCursorPosition() (float64, float64) {
	cx, cy := ebiten.CursorPosition()
	return float64(cx) / ebiten.DeviceScaleFactor(), float64(cy) / ebiten.DeviceScaleFactor()
}

func CenteredRect(xc, yc, halfH, halfW int) image.Rectangle {
	return image.Rect(xc-halfH, yc-halfH, xc+halfH, xc+halfH)
}

var random *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UnixNano())
	random = rand.New(s)
}
