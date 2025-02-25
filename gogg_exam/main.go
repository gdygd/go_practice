package main

import (
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var ttfpath string = "/home/theroad/notosans/laravel-notosans-korean/public/fonts/NotoSansKR-Bold.ttf"

func gofont() {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(1024, 1024)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello, world!", 512, 512, 0.5, 0.5)
	dc.SavePNG("gofont.png")
}

func gradiendtext() {
	var W float64 = 1024
	var H float64 = 512
	dc := gg.NewContext(1024, 512)

	// draw text
	dc.SetRGB(0, 0, 0)
	//dc.LoadFontFace("/Library/Fonts/Impact.ttf", 128)
	dc.LoadFontFace(ttfpath, 72)
	dc.DrawStringAnchored("Gradient Text", W/2, H/2, 0.5, 0.5)

	// get the context as an alpha mask
	mask := dc.AsMask()

	// clear the context
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// set a gradient
	g := gg.NewLinearGradient(0, 0, W, H)
	g.AddColorStop(0, color.RGBA{255, 0, 0, 255})
	g.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	dc.SetFillStyle(g)

	// using the mask, fill the context with the gradient
	dc.SetMask(mask)
	dc.DrawRectangle(0, 0, W, H)
	dc.Fill()

	dc.SavePNG("gradiendtext.png")

}

func memo() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	if err := dc.LoadFontFace(ttfpath, 72); err != nil {
		panic(err)
	}
	dc.SetRGB(0, 0, 0)
	s := "ONE DOES NOT SIMPLY"
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		//fmt.Printf("str:")
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := S/2 + float64(dx)
			y := S/2 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)

			//fmt.Printf("%s ", s)

		}
	}
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, S/2, S/2, 0.5, 0.5)
	//fmt.Printf("\n##%s \n", s)
	dc.SavePNG("memo.png")

}

func drawStringWrapped() {

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(1024, 1024)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	//dc.DrawStringAnchored("Hello, world!", 512, 512, 0.5, 0.5)
	dc.DrawStringWrapped("Hello, world!", 512, 512, 0.1, 0.1, 100, 0, gg.AlignLeft)
	dc.SavePNG("drawStringWrapped.png")

}

func drawtext() {
	const S = 300
	// im, err := gg.LoadImage("src.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(ttfpath, 72); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)

	dc.DrawRoundedRectangle(0, 0, 150, 150, 0)
	// dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	dc.Clip()
	dc.SavePNG("drawtext.png")

}

func main() {
	gofont()

	gradiendtext()

	memo()

	drawStringWrapped()

	drawtext()

}
