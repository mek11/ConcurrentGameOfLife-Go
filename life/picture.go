package life

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

// Pixel represent a cell to be drawn
type pixel struct {
	x, y, v int
}

// ImgCreator allow to create an image
type ImgCreator interface {
	create(chPixels <-chan pixel, chDone chan<- struct{})
}

// GifCreator allows to create a gif picture
type GifCreator struct {
	width, height, dimPx int
}

// NewGifCreator return a new gif creator
func NewGifCreator(width int, height int) *GifCreator {
	return &GifCreator{
		width:  width,
		height: height,
		dimPx:  5,
	}
}

func (g *GifCreator) create(chPixels <-chan pixel, chDone chan<- struct{}) {
	var images []*image.Paletted
	var delays []int
	palette := []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, //black
		color.RGBA{0x00, 0xff, 0xff, 0xff}, //cyan
		color.RGBA{0x00, 0x00, 0xff, 0xff}, //blue
		color.RGBA{0x00, 0xff, 0x00, 0xff}, //green
		color.RGBA{0xff, 0x00, 0x00, 0xff}, //red
		color.RGBA{0xff, 0x00, 0xff, 0xff}, //magenta
		color.RGBA{0xff, 0xff, 0x00, 0xff}, //yellow
		color.RGBA{0xff, 0xff, 0xff, 0xff}, //white
	}

	counter := 0
	img := image.NewPaletted(image.Rect(0, 0, g.width*g.dimPx, g.height*g.dimPx), palette)

	for p := range chPixels {
		for i := 0; i < g.dimPx; i++ {
			for j := 0; j < g.dimPx; j++ {
				img.SetColorIndex(g.dimPx*p.x+i, g.dimPx*p.y+j, uint8(p.v))
			}
		}
		counter++
		if counter == g.width*g.height {
			images = append(images, img)
			delays = append(delays, 0)
			img = image.NewPaletted(image.Rect(0, 0, g.width*g.dimPx, g.height*g.dimPx), palette)
			counter = 0
		}
	}

	f, err := os.Create("game.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	chDone <- struct{}{}
}

// TestCreator allow to test the creation of a picture
type TestCreator struct{}

// NewTestCreator return a new test img creator
func NewTestCreator() *TestCreator {
	return &TestCreator{}
}

func (t *TestCreator) create(chPixels <-chan pixel, chDone chan<- struct{}) {
	for range chPixels {
		// draw a cell in the current image
		// if an image is complete, set up the next one
	}
	// create the final image
	chDone <- struct{}{}
}
