package scene

import (
	"image/color"
	"os/exec"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/reverset-raiders/ui"
)

type CreditsScene struct {
	SwitchSceneFunc func(string)
}

func NewCreditsScene(f func(string)) *CreditsScene {
	return &CreditsScene{
		SwitchSceneFunc: f,
	}
}

func (c *CreditsScene) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		c.SwitchSceneFunc("menu")
		/*
			if c.EbitengineText.In(cursorX, cursorY) {
				open("https://ebitengine.org/")
			}*/
	}
}

func (c *CreditsScene) Draw(screen *ui.ScaledScreen) {
	screen.DrawTextCenteredAt("Credits", 48, 480, 50, color.White)
	screen.DrawText("ebitengine -- https://ebitengine.org", 32, 50, 100, color.White)
	screen.DrawText("etxt -- https://github.com/tinne26/etxt", 32, 50, 140, color.White)
	screen.DrawText("font: Roboto Medium -- https://github.com/googlefonts/roboto", 32, 50, 180, color.White)

	screen.DrawTextCenteredAt("click anywhere to return", 16, 480, 450, color.White)
}

// todo: will this work for webasm?
func open(url string) error {
	var cmd string
	var args []string
	println("open")
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
