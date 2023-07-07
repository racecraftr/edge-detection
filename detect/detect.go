package detect

import (
	"image"
	"image/color"
	"math"

	"github.com/schollz/progressbar/v3"
)

type Point struct {
	X, Y int
}

type PointSet map[Point]bool

func getLuma(c color.Color) uint32 {
	l := color.GrayModel.Convert(c)
	y, _, _, _ := l.RGBA()
	return y / 257
}

func isEdge(c1, c2 color.Color) bool {
	rawl1, rawl2 := getLuma(c1), getLuma(c2)

	l1, l2 := float64(rawl1)+0.05, float64(rawl2)+0.05
	l1, l2 = math.Max(l1, l2), math.Min(l1, l2)
	diff := (l1-l2)/l1*100
	return diff >= 17
}

func CreateBar(maxsize int, desc string) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		maxsize,
		progressbar.OptionShowBytes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "[red]=[reset]",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func FindEdges(img image.Image) PointSet {
	bounds := img.Bounds().Max
	width, height := bounds.X, bounds.Y
	size, px := width*height, 0
	res := make(PointSet)

	bar := CreateBar(size, "reading pixels...")
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p := Point{X: x, Y: y}
			c1 := img.At(x, y)
			if x < width-1 {
				c2 := img.At(x+1, y)
				if isEdge(c1, c2) {
					res[p] = true
				}
			}
			if y < height-1 {
				c2 := img.At(x, y+1)
				if isEdge(c1, c2) {
					res[p] = true
				}
			}
			px++
			bar.Set(px)

		}
	}
	return res
}
