package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	//"github.com/mskrha/svg2png"
	svg "github.com/ajstarks/svgo"
	//"github.com/ajstarks/svgo/float"
	"github.com/fogleman/gg"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	//dimg "github.com/llgcode/draw2d/draw2dimg"
	"github.com/JoshVarga/svgparser"
	dsvg "github.com/llgcode/draw2d/draw2dsvg"
)

func convert(strfile string) {
	w, h := 512, 512

	//in, err := os.Open("file_example_SVG_20kB.svg")
	in, err := os.Open(strfile)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	icon, _ := oksvg.ReadIconStream(in)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	out, err := os.Create("out5.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", out)
}

func convert1() {
	w, h := 512, 512

	//in, err := os.Open("file_example_SVG_20kB.svg")
	in, err := os.Open("output.svg")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	icon, _ := oksvg.ReadIconStream(in)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	out, err := os.Create("svg4.bmp")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rgba)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", out)

}

// func convert2() {
// 	// Input and output file paths
// 	svgFilePath := "svg2.svg"
// 	bmpFilePath := "output.bmp"

// 	// Create a new image context
// 	dc := gg.NewContext(0, 0)
// 	err := dc.LoadSVG(svgFilePath)
// 	if err != nil {
// 		fmt.Println("Failed to load SVG file:", err)
// 		return
// 	}

// 	// Set the size of the context to match the SVG image
// 	width, height := dc.Width(), dc.Height()
// 	dc = gg.NewContext(int(width), int(height))
// 	dc.LoadSVG(svgFilePath)

// 	// Save the image as BMP
// 	dc.SaveBMP(bmpFilePath)

// 	// Open the resulting BMP file
// 	cmd := exec.Command("open", bmpFilePath) // Modify for Windows or Linux
// 	err = cmd.Run()
// 	if err != nil {
// 		fmt.Println("Failed to open the BMP file:", err)
// 		return
// 	}

// 	fmt.Println("SVG to BMP conversion completed successfully.")

// }

func convert3() {
	//var input []byte
	// input, err := os.Open("svg2.svg")

	// if err != nil {
	// 	panic(err)
	// }
	// defer input.Close()

	// // Fill the input with SVG image

	// converter := svg2png.New()

	// output, err := converter.Convert(input)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// out, err := os.Create("out.bmp")
	// if err != nil {
	// 	panic(err)
	// }
	// defer out.Close()

	// err = png.Encode(output, output)
	// if err != nil {
	// 	panic(err)
	// }
}

func convert4() {

	// Create a new SVG file
	file, err := os.Create("output2.svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new SVGo object
	canvas := svg.New(file)

	// Set the canvas size
	canvas.Start(500, 500)

	// Draw a rectangle
	canvas.Rect(100, 100, 300, 200, "fill:none;stroke:black")

	// Draw a circle
	canvas.Circle(250, 250, 50, "fill:red")

	//canvas.Image(0, 0, 512, 512, "1.jpg", "0.50")
	canvas.Text(300, 300, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")

	// End the SVG file
	canvas.End()

	convert("output2.svg")
}

// func convert5() {

// 	svgData, err := ioutil.ReadFile("output2.svg")
// 	if err != nil {
// 		// Handle the error
// 	}

// 	outputFile, err := os.Create("output2.png")
// 	if err != nil {
// 		// Handle the error
// 	}
// 	defer outputFile.Close()

// 	canvas := svg.New(outputFile)
// 	canvas.Start(500, 500) // Specify the width and height of the canvas

// 	svgReader := strings.NewReader(string(svgData))
// 	//err = canvas.Parse(svgReader, float.DefaultErrFmt)
// 	err = canvas.Parse(svgReader, svg.DefaultErrFmt)
// 	if err != nil {
// 		// Handle the error
// 	}
// 	canvas.End()

// 	canvas.Writer.Flush()

// }

func convert6() {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	// if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 96); err != nil {
	// 	panic(err)
	// }
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	//dc.SavePNG("out6.png")
	dc.SavePNG("out6.bmp")

}

func convert7() {
	//w, h := 1024, 1024
	const S = 1024

	// // load svg and convert rgba
	// in, err := os.Open("output.svg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer in.Close()

	// icon, _ := oksvg.ReadIconStream(in)
	// icon.SetTarget(0, 0, float64(w), float64(h))
	// rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	// load image
	img, err := gg.LoadImage("svg4.bmp")
	if err != nil {
		fmt.Printf("err : %v", err)
	}

	// image modify
	dc := gg.NewContextForImage(img)
	// dc.SetRGB(1, 1, 1)
	// dc.Clear()
	dc.SetRGB(0, 0, 0)
	// if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 96); err != nil {
	// 	panic(err)
	// }
	dc.DrawStringAnchored("Hello, world!", 200, 200, 0.5, 0.5)

	dc.SavePNG("out2.bmp")

}

type SvgStruct struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr,omitempty"`
	Height  string   `xml:"height,attr,omitempty"`
	ViewBox string   `xml:"viewBox,attr,omitempty"`

	Image    []dsvg.Image    `xml:"image"`
	Path     []dsvg.Path     `xml:"path"`
	Position []dsvg.Position `xml:"position"`
	Rect     []dsvg.Rect     `xml:"rect"`
	Rune     []dsvg.Rune     `xml:"rune"`
	Text     []dsvg.Text     `xml:"text"`
}

func convert8() {
	// decoding svg
	// draw2d library

	var svg *dsvg.Svg = dsvg.NewSvg()
	//var svg SvgStruct

	fmt.Printf("(1)\n%v\n", svg)

	dat, err := ioutil.ReadFile("./svg4.svg")
	check(err)
	fmt.Printf("%s\n", string(dat))

	//var members Members
	xmlerr := xml.Unmarshal(dat, svg)
	if xmlerr != nil {
		panic(xmlerr)
	}

	fmt.Printf("\n-------------------------)\nUnmarshal.....\n")
	fmt.Println(svg)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func convert9() {
	dat, err := ioutil.ReadFile("./svg4.svg")
	check(err)
	strSvg := string(dat)

	fmt.Printf("(1)\n%s\n", string(dat))

	reader := strings.NewReader(strSvg)
	element, _ := svgparser.Parse(reader, false)

	fmt.Printf("(2)\nelement : %v\n", element)

	rectangles := element.FindAll("text")
	fmt.Printf("(3)\r rectangles : %v\n", rectangles[0])

	// rectangles := element.FindAll("rect")

	// fmt.Printf("First child fill: %s\n", rectangles[0].Attributes["fill"])
	// fmt.Printf("Second child height: %s", rectangles[1].Attributes["height"])

}

func main() {

	//convert1()
	//convert4()
	//convert6()
	//convert7()
	//convert8()
	convert9()
}
