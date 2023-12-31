package detect

import (
	"image"
	"image/color"
	"math"
	"sync"
)

var (
	horizontal = [3][3]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	vertical = [3][3]int{
		{1, 2, 1},
		{0, 0, 0},
		{-1, -2, -1},
	}
)

func Luminance(aColor color.Color) float64 {

	red, green, blue, alpha := aColor.RGBA()
	red /= 257
	green /= 257
	blue /= 257
	alpha /= 257

	// need to convert uint32 to float64
	return float64(float64(0.299)*float64(red)+float64(0.587)*float64(green)+float64(0.114)*float64(blue)) - float64(alpha)
}

const BUFSIZE = 0x7fffffff

// uses sobel gradient instead of simple difference approach. :P
func FindEdgesV2(img image.Image) []*Point {
	res, idx := make([]*Point,BUFSIZE), 0
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	size := width*height - width*2 - height*2 + 4
	px := 0
	bar := CreateBar(size, "reading pixels...")
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			gradient := [3][3]int{}
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					gradient[i][j] = int(Luminance(img.At(x-1+i, y-1+i)))
				}
			}

			gx, gy := 0, 0
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					gx += gradient[i][j] * horizontal[i][j]
					gy += gradient[i][j] * vertical[i][j]
				}
			}
			colorCode := int(math.Sqrt(float64(gx*gx + gy*gy)))
			if colorCode > 80 {
				res[idx] = &Point{x, y}
				idx ++
			}

			px++
			bar.Set(px)
		}
	}
	return res[:idx]
}

func FindEdgesV2Async(img image.Image) []*Point {
	res := make([]*Point, BUFSIZE)

	c := make(chan *Point)
	wg := sync.WaitGroup{}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	size := width*height - width*2 - height*2 + 4
	px := 0
	bar := CreateBar(size, "reading pixels...")
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			wg.Add(1)
			go func(x, y int, c chan *Point) {
				defer wg.Done()
				gradient := [3][3]int{}
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						gradient[i][j] = int(Luminance(img.At(x-1+i, y-1+i)))
					}
				}

				gx, gy := 0, 0
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						gx += gradient[i][j] * horizontal[i][j]
						gy += gradient[i][j] * vertical[i][j]
					}
				}
				colorCode := int(math.Sqrt(float64(gx*gx + gy*gy)))
				if colorCode > 80 {
					c <- &Point{x, y}
				}
			}(x, y, c)

			px++
			bar.Set(px)
		}
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	idx := 0
	for p := range c {
		res[idx] = p
	}
	return res[:idx]
}
