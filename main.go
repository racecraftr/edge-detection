package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"racecraftr/edge-detection/detect"

	"github.com/urfave/cli/v2"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

var app = &cli.App{
	Action: func(ctx *cli.Context) error {
		path := ctx.Args().Get(0)
		// we now have the path to the file.
		f, err := os.Open(path)
		Check(err)
		defer f.Close()

		img, _, err := image.Decode(f)
		Check(err)
		points := detect.FindEdgesV2(img)
		fmt.Println()

		resName := filepath.Base(path) + "-edges.png"
		width, height := img.Bounds().Dx(), img.Bounds().Dy()

		resImg := image.NewRGBA(image.Rect(
			0, 0,
			width, height,
		))

		wrote := 0
		bar := detect.CreateBar(len(points), "creating image...")
		for _, p := range points {
			if p == nil {
				break
			}
			resImg.Set(p.X, p.Y, color.RGBA{0xff, 0x0, 0x0, 0xff})
			wrote++
			bar.Set(wrote)
		}

		resFile, err := os.Create(resName)
		Check(err)
		defer resFile.Close()

		fmt.Println()
		fmt.Println("encoding image...")
		png.Encode(resFile, resImg)

		return nil
	},
}

func main() {
	app.Run(os.Args)
}
