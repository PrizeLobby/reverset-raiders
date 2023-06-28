package res

import (
	"bytes"
	"image"
	"log"
	"path"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/sfnt"

	"embed"
)

//go:embed font/* img/*.png
var assets embed.FS

var fonts map[string]*sfnt.Font = make(map[string]*sfnt.Font)

func init() {
	LoadFonts()
}

func LoadFonts() {
	bytes, err := assets.ReadFile("font/Roboto-Medium.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := sfnt.Parse(bytes)

	if err != nil {
		log.Fatal(err)
	}

	fonts["Roboto-Medium"] = f
}

func GetFont(n string) *sfnt.Font {
	return fonts[n]
}

func ReadImage(p string) (image.Image, error) {
	data, err := assets.ReadFile(path.Join("img", p))
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// Images is the map of all loaded images.
var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

// GetImage returns the image matching the given file name. IT ALSO LOADS IT.
func GetImage(p string) *ebiten.Image {
	if v, ok := Images[p]; ok {
		return v
	}
	img, err := ReadImage(p + ".png")
	if err != nil {
		log.Println("error reading image " + p)
		return nil
	}
	eimg := ebiten.NewImageFromImage(img)
	//eimg = ScaledImage(eimg, ebiten.DeviceScaleFactor())
	Images[p] = eimg
	return eimg
}

func ScaledImage(image *ebiten.Image, scaleAmount float64) *ebiten.Image {
	w, h := image.Bounds().Dx(), image.Bounds().Dy()
	scaledImage := ebiten.NewImage(int(scaleAmount*float64(w)), int(float64(h)*scaleAmount))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleAmount, scaleAmount)
	scaledImage.DrawImage(image, op)
	return scaledImage
}
