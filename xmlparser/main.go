package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/mskrha/svg2png"
)

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   int      `xml:"width,attr,omitempty"`
	Height  int      `xml:"height,attr,omitempty"`
	ViewBox string   `xml:"viewBox,attr,omitempty"`

	Paths    []Path    `xml:"path"`
	Texts    []Text    `xml:"text"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
	Images   []SImage  `xml:"image"`
}

type SvgGrp struct {
	Id          int    `xml:"id,attr"`
	Stroke      string `xml:"stroke,attr"`
	StrokeWidth string `xml:"stroke-width,attr"`
	Fill        string `xml:"fill,attr"`
	FontFamily  string `xml:"font-family,attr,omitempty"`

	Paths    []Path    `xml:"path"`
	Texts    []Text    `xml:"text"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
	Images   []SImage  `xml:"image"`
}

type Path struct {
	// M, L Z만 지원
	XMLName xml.Name `xml:"path"`
	Desc    string   `xml:"d,attr"`
	Pos     [][2]float64
}

type Text struct {
	XMLName    xml.Name `xml:"text"`
	Id         int      `xml:"id,attr"`
	FontSize   float64  `xml:"font-size,attr,omitempty"`
	FontFamily string   `xml:"font-family,attr,omitempty"`
	Text       string   `xml:",innerxml"`
	Style      string   `xml:"style,attr,omitempty"`
	X          float64  `xml:"x,attr"`
	Y          float64  `xml:"y,attr"`
	Fill       string   `xml:"fill,attr"`

	rgb [3]int
}

type Rect struct {
	XMLName xml.Name `xml:"rect"`
	Id      int      `xml:"id,attr"`
	Style   string   `xml:"style,attr,omitempty,fill"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
	Fill    string   `xml:"fill,attr"`
	rgb     [3]int
}

type Polygon struct {
	XMLName   xml.Name `xml:"polygon"`
	Id        int      `xml:"id,attr"`
	Style     string   `xml:"style,attr,omitempty,fill"`
	StrPoints string   `xml:"points,attr"`
	Points    [][2]float64

	Fill   string `xml:"fill,attr"`
	Stroke string `xml:"stroke,attr"`

	Fillrgb [3]int
	Linergb [3]int
}

type SImage struct {
	XMLName xml.Name `xml:"image"`
	Id      int      `xml:"id,attr"`
	Href    string   `xml:"href,attr"`
}

type Svg2 struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   int      `xml:"width,attr,omitempty"`
	Height  int      `xml:"height,attr,omitempty"`
	ViewBox string   `xml:"viewBox,attr,omitempty"`

	Paths    []Path    `xml:"path"`
	Texts    []Text    `xml:"text"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
}

type VoidTag struct {
	XMLName1 xml.Name `xml:"text"`
	XMLName2 xml.Name `xml:"rect"`
	XMLName3 xml.Name `xml:"polygon"`
}

func svgpathtoPos(path string) {
	//
	var fpos [][2]float64 = [][2]float64{}

	posinfo := strings.Split(path, " ")

	var i int = 0
	var posxy [2]float64
	for _, pos := range posinfo {
		fmt.Println(pos)
		var strNum string = ""
		fmt.Printf("ch:%c\n", pos[0])
		if (pos[0] >= 0x41 && pos[0] <= 0x5A) ||
			(pos[0] >= 0x61 && pos[0] <= 0x7A) {

			strNum = pos[1:]
		} else {
			strNum = pos[:]
		}

		if len(strNum) > 0 {
			if i%2 == 0 {
				fval, _ := strconv.ParseFloat(strNum, 0)
				posxy[0] = fval
			} else {
				fval, _ := strconv.ParseFloat(strNum, 0)
				posxy[1] = fval

				fpos = append(fpos, posxy)
				posxy[0] = 0
				posxy[1] = 0
			}

			fmt.Printf("strnum(%d)(idx:%d):%s \n", i, i/2, strNum)
			i++
		}

	}

	fmt.Printf("poslist %v\n", fpos)
}

func convert1() {
	//var input []byte

	// Fill the input with SVG image

	//dat, err := ioutil.ReadFile("./svg2.svg")
	//dat, err := ioutil.ReadFile("./rect_text.svg")
	dat, err := ioutil.ReadFile("./pathsvg.svg")

	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}

	converter := svg2png.New()

	output, err := converter.Convert(dat)
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := os.Create("rect_text2.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	werr := ioutil.WriteFile("rect_text2.png", output, 0644)
	if werr != nil {
		fmt.Printf("werr %v", werr)
	}

	// Read the output PNG image

}

type Point struct {
	X, Y float64
}

func TmpPolygon(n int) []Point {
	result := make([]Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = Point{math.Cos(a), math.Sin(a)}
	}
	return result
}

func TempDrawing() {
	const W = 1200
	const H = 120
	const S = 100
	dc := gg.NewContext(W, H)
	dc.SetHexColor("#FFFFFF")
	dc.Clear()
	n := 5
	points := TmpPolygon(n)
	for x := S / 2; x < W; x += S {
		dc.Push()
		s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), H/2)
		dc.Rotate(rand.Float64() * 2 * math.Pi)
		dc.Scale(s, s)
		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor("#FF0000")
		dc.StrokePreserve()
		dc.SetHexColor("#00FF00")
		dc.Fill()
		dc.Pop()
	}
	dc.SavePNG("stars.png")
}

func drawRect(r Rect, dc *gg.Context) {
	dc.Clear()

	//dc.SetRGBA255(r[i].rgb[0], r[i].rgb[1], r[i].rgb[2], 255)
	dc.SetHexColor(r.Fill)

	fmt.Printf("Rect X(%f) Y(%f) \n", r.X, r.Y)
	fmt.Printf("Rect W(%f) H(%f) \n", r.Width, r.Height)

	dc.DrawRectangle(r.X, r.Y, r.Width, r.Height)

	dc.Fill()

}

func drawText(t Text, dc *gg.Context) {
	dc.SetHexColor(t.Fill)

	// text
	//dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", 30)
	dc.LoadFontFace("/home/theroad/notosans/NotoSansKR-Regular.ttf", 30)
	dc.DrawString(t.Text, t.X, t.Y)
	dc.Fill()
}

func drawImage(simg SImage, dc *gg.Context) {

	img, errimg := gg.LoadImage(simg.Href)
	if errimg != nil {
		fmt.Printf("errimg : %v", errimg)
	}

	dc.DrawImage(img, 0, 0)
}

func drawPolygon(p Polygon, dc *gg.Context) {
	dc.Clear()

	dc.Push()
	for j := 0; j < len(p.Points); j++ {
		fmt.Printf("drawPolygon point ###(%f, %f) \n", p.Points[j][0], p.Points[j][1])
		if j == 0 {
			//dc.MoveTo(p[i].Points[j][0], p[i].Points[j][1])
			dc.LineTo(p.Points[j][0], p.Points[j][1])
		} else {
			dc.LineTo(p.Points[j][0], p.Points[j][1])
		}
	}

	dc.SetHexColor(p.Stroke)
	dc.SetLineWidth(5)

	dc.StrokePreserve()
	dc.SetHexColor(p.Fill)

	dc.Fill()
	dc.Pop()
}

func drawRects(r []Rect, dc *gg.Context) {

	dc.Clear()
	//dc.SetRGBA(0, 1, 0, 1)
	// dc.SetRGBA255(0, 255, 0, 255)

	for i := 0; i < len(r); i++ {
		//dc.SetRGBA255(r[i].rgb[0], r[i].rgb[1], r[i].rgb[2], 255)
		dc.SetHexColor("#00FF00")

		fmt.Printf("Rect X(%f) Y(%f) \n", r[i].X, r[i].Y)
		fmt.Printf("Rect W(%f) H(%f) \n", r[i].Width, r[i].Height)

		dc.DrawRectangle(r[i].X, r[i].Y, r[i].Width, r[i].Height)
	}
	dc.Fill()
}

func drawTexts(t []Text, dc *gg.Context) {

	//dc.SetRGB(0, 0, 1)
	//dc.SetRGB255(0, 0, 255)
	dc.SetHexColor("#0000FF")

	// text
	dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", 30)
	for i := 0; i < len(t); i++ {
		dc.SetRGB255(t[i].rgb[0], t[i].rgb[1], t[i].rgb[2])
		dc.DrawString(t[i].Text, t[i].X, t[i].Y)
		dc.Fill()
	}

}

func drawPolygons(p []Polygon, dc *gg.Context) {
	// Draw Polygon

	dc.Clear()
	//dc.SetRGB(1, 1, 0)

	dc.Clear()

	for i := 0; i < len(p); i++ {
		//dc.SetRGB255(p[i].rgb[0], p[i].rgb[1], p[i].rgb[2])
		dc.SetHexColor("#FFFF00")

		dc.Push()

		for j := 0; j < len(p[i].Points); j++ {
			fmt.Printf("drawPolygon point ###(%f, %f) \n", p[i].Points[j][0], p[i].Points[j][1])
			if j == 0 {
				//dc.MoveTo(p[i].Points[j][0], p[i].Points[j][1])
				dc.LineTo(p[i].Points[j][0], p[i].Points[j][1])
			} else {
				dc.LineTo(p[i].Points[j][0], p[i].Points[j][1])
			}
		}

		//dc.SetHexColor("#FF00FF")
		dc.SetHexColor(p[i].Stroke)
		dc.SetLineWidth(5)

		dc.StrokePreserve()
		//dc.SetHexColor("#FFFF00")
		dc.SetHexColor(p[i].Fill)

		dc.Fill()
		dc.Pop()
	}

	// dc.SetRGB(1, 1, 1)
	// dc.Fill()
	// dc.ClosePath()
	// dc.Stroke()

}

func ColorHex2RGB(s string) [3]int {
	// s : #ff0000 (RRGGBB)
	// s = "#ff0000"
	//check string

	var rgb [3]int
	// check first char
	if s[0] == '#' {
		fmt.Printf("first char = # \n")
	} else {
		fmt.Printf("Invalid first char = # \n")
		rgb[0] = 0
		rgb[1] = 0
		rgb[2] = 0
		return rgb
	}

	var strR string = s[1:3]
	fmt.Printf("strR : %s\n", strR)

	var strG string = s[3:5]
	fmt.Printf("strG : %s\n", strG)

	var strB string = s[5:7]
	fmt.Printf("strB : %s\n", strB)

	valR, _ := strconv.ParseInt(strR, 16, 16)
	valG, _ := strconv.ParseInt(strG, 16, 16)
	valB, _ := strconv.ParseInt(strB, 16, 16)

	fmt.Printf("rr : %d, %d, %d \n", valR, valG, valB)

	rgb[0] = int(valR)
	rgb[1] = int(valG)
	rgb[2] = int(valB)

	return rgb
}

func xmlDecode() {
	//xmlif := []interface{}{}
	xmlif := map[string]interface{}{}

	//xmlFile, err := os.Open("rect_text2.svg")
	xmlFile, err := os.Open("pathsvg.svg")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	//--------------------------------------
	// xml interface test
	decoder := xml.NewDecoder(xmlFile)
	err1 := decoder.Decode(&xmlif)
	if err1 != nil {
		fmt.Printf("xmlDecode err:%s\n", err1.Error())
		return
	} else {
		fmt.Printf("xmlDecode xmlif.. (%v)\n", xmlif)
	}

	//--------------------------------------

}

func Drawsvg(svginfo Svg, dc *gg.Context) {
	const maxid = 100

	//for targetid ~100
	//  for rects
	//	for texts
	//	for polygons

	for targetid := 1; targetid <= maxid; {
		var markId = targetid

		for i := 0; i < len(svginfo.Rects); i++ {
			if svginfo.Rects[i].Id == targetid {
				drawRect(svginfo.Rects[i], dc)
				fmt.Printf("drawRect %d\n", svginfo.Rects[i].Id)
				targetid += 1
			}
		}

		for i := 0; i < len(svginfo.Polygons); i++ {
			if svginfo.Polygons[i].Id == targetid {
				drawPolygon(svginfo.Polygons[i], dc)
				fmt.Printf("drawPolygon %d\n", svginfo.Polygons[i].Id)
				targetid += 1
			}
		}

		for i := 0; i < len(svginfo.Texts); i++ {
			if svginfo.Texts[i].Id == targetid {
				drawText(svginfo.Texts[i], dc)
				fmt.Printf("drawTexts %d\n", svginfo.Texts[i].Id)
				targetid += 1
			}
		}

		for i := 0; i < len(svginfo.Images); i++ {
			if svginfo.Images[i].Id == targetid {
				drawImage(svginfo.Images[i], dc)
				fmt.Printf("drawImage %d\n", svginfo.Images[i].Id)
				targetid += 1
			}
		}

		if markId == targetid {
			targetid += 1
		}
	}

	// // search rect
	// for idx := 1; idx <= renderIdx; {
	// 	//search rect
	// 	if svginfo.Rects[idx-1].Id == renderIdx {
	// 		//drawRect
	// 		fmt.Printf("drawRect %d\n", svginfo.Rects[idx-1].Id)
	// 		idx += 1
	// 	} else if svginfo.Polygons[idx-1].Id == renderIdx {
	// 		//drawPolygon
	// 		fmt.Printf("drawPolygon %d\n", svginfo.Polygons[idx-1].Id)
	// 		idx += 1
	// 	} else if svginfo.Texts[idx-1].Id == renderIdx {
	// 		//drawText
	// 		fmt.Printf("drawText %d\n", svginfo.Texts[idx-1].Id)
	// 		idx += 1
	// 	} else {
	// 		fmt.Printf("No search id.. %d\n", idx)
	// 		idx += 1
	// 	}
	// }

}

func main() {
	xmlDecode()

	convert1()

	TempDrawing()

	ColorHex2RGB("#ff0000")

	fmt.Printf("xml parser...")

	// Open our xmlFile
	xmlFile, err := os.Open("rect_text.svg")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// //--------------------------------------
	// // xml interface test
	// decoder := xml.NewDecoder(xmlFile)
	// err1 := decoder.Decode(&xmlif)
	// if err1 != nil {
	// 	fmt.Println(err1.Error())
	// 	return
	// } else {
	// 	fmt.Printf("xmlif.. (%v)", xmlif...)
	// }

	//--------------------------------------

	fmt.Println("Successfully Opened users.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var svginfo Svg
	xml.Unmarshal(byteValue, &svginfo)

	fmt.Printf("%s\n", string(byteValue))

	fmt.Printf("svr info %v\n", svginfo)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	fmt.Printf("Path))================================================>>\n")
	for i := 0; i < len(svginfo.Paths); i++ {
		fmt.Printf("xmlname: %v\n", svginfo.Paths[i].XMLName)
		fmt.Println("Path Desc: " + svginfo.Paths[i].Desc)
		svgpathtoPos(svginfo.Paths[i].Desc)
	}
	fmt.Printf("Rect))================================================>>\n")

	for i := 0; i < len(svginfo.Rects); i++ {
		fmt.Printf("xmlname: %v\n", svginfo.Rects[i].XMLName)
		fmt.Printf("Rect Id:%d \n", svginfo.Rects[i].Id)
		fmt.Println("Rect style: " + svginfo.Rects[i].Style)
		fmt.Printf("Rect X(%f) Y(%f) \n", svginfo.Rects[i].X, svginfo.Rects[i].Y)
		fmt.Printf("Rect W(%f) H(%f) \n", svginfo.Rects[i].Width, svginfo.Rects[i].Height)
		fmt.Println("Rect fill: " + svginfo.Rects[i].Fill)

		svginfo.Rects[i].rgb = ColorHex2RGB(svginfo.Rects[i].Fill)
	}

	fmt.Printf("Polygon))============================================>>\n")

	for i := 0; i < len(svginfo.Polygons); i++ {
		fmt.Printf("xmlname: %v\n", svginfo.Polygons[i].XMLName)
		fmt.Printf("Polygon Id:%d \n", svginfo.Polygons[i].Id)
		fmt.Println("Polygon style: " + svginfo.Polygons[i].Style)
		fmt.Printf("Polygon Fillrgb:%s ", svginfo.Polygons[i].Fillrgb)
		fmt.Printf("Polygon Linergb:%s ", svginfo.Polygons[i].Linergb)
		fmt.Printf("Polygon points(%s) \n", svginfo.Polygons[i].StrPoints)

		svginfo.Polygons[i].Fillrgb = ColorHex2RGB(svginfo.Polygons[i].Fill)
		svginfo.Polygons[i].Linergb = ColorHex2RGB(svginfo.Polygons[i].Stroke)

		svginfo.Polygons[i].Points = [][2]float64{}

		// points parsing
		//strtmp := "alfa bravo charlie delta echo foxtrot golf"
		points := strings.Split(svginfo.Polygons[i].StrPoints, " ")
		for _, strXY := range points {
			// strXY : x,y
			xy := strings.Split(strXY, ",")
			fmt.Printf("X[%s], Y[%s]\n", xy[0], xy[1])
			nx, _ := strconv.ParseFloat(xy[0], 2)
			ny, _ := strconv.ParseFloat(xy[1], 2)

			svginfo.Polygons[i].Points = append(svginfo.Polygons[i].Points, [2]float64{nx, ny})

		}

		fmt.Printf("Polygon point arr(%v) \n", svginfo.Polygons[i].Points)
	}

	fmt.Printf("Text================================================>>\n")

	for i := 0; i < len(svginfo.Texts); i++ {
		fmt.Printf("xmlname: %v\n", svginfo.Texts[i].XMLName)
		fmt.Printf("Text Id:%d \n", svginfo.Texts[i].Id)
		//fmt.Println("Text Type: " + svginfo.Texts[i]. FontSize)
		fmt.Println("Text font: " + svginfo.Texts[i].FontFamily)
		fmt.Println("Text text: " + svginfo.Texts[i].Text)
		fmt.Println("Text style: " + svginfo.Texts[i].Style)
		fmt.Println("Text fill: " + svginfo.Texts[i].Fill)
		fmt.Printf("Text X(%f) Y(%f) \n", svginfo.Texts[i].X, svginfo.Texts[i].Y)

		svginfo.Texts[i].rgb = ColorHex2RGB(svginfo.Texts[i].Fill)
	}
	fmt.Printf("-----------------------------------------------------\n")
	fmt.Printf("=====================================================\n")

	// 여기까지 SVG parsing!

	// 여기서부터 drawing

	// svg size
	dc := gg.NewContext(svginfo.Width, svginfo.Height)

	// // Draw rectangulrar
	// fmt.Printf("Draw rect==>>\n")
	// drawRects(svginfo.Rects, dc)

	// // Draw Polygon
	// fmt.Printf("Draw polygon==>>\n")
	// drawPolygons(svginfo.Polygons, dc)

	// // Draw Polygon
	// fmt.Printf("Draw text==>>\n")
	// drawTexts(svginfo.Texts, dc)

	// load and draw image

	fmt.Printf("########drawsvg======================>\n")
	// draw rect, text, polygon, image
	Drawsvg(svginfo, dc)
	fmt.Printf("########=============================>\n")

	// img, errimg := gg.LoadImage("1.jpg")
	// if errimg != nil {
	// 	fmt.Printf("errimg : %v", errimg)
	// }

	// dc.DrawImage(img, 50, 50)

	//dc.SavePNG("rect_text1.png")
	dc.SavePNG("rect_text2.bmp")

	var aaa []int = []int{}
	aaa = append(aaa, 1)
	aaa = append(aaa, 2)
	aaa = append(aaa, 3)

	fmt.Printf("aaa1:(%v)\n", aaa)
	if len(aaa) > 0 {
		aaa = aaa[1:]
		fmt.Printf("aaa2:(%v)\n", aaa)
	}

	if len(aaa) > 0 {
		aaa = aaa[1:]
		fmt.Printf("aaa3:(%v) %d\n", aaa, len(aaa))
	}

	if len(aaa) > 0 {
		aaa = aaa[1:]
		fmt.Printf("aaa4:(%v) %d\n", aaa, len(aaa))
	}

	if len(aaa) > 0 {
		aaa = aaa[1:]
		fmt.Printf("aaa5:(%v) %d\n", aaa, len(aaa))
	}

	fmt.Printf("aaa###:(%v) %d\n", aaa, len(aaa))

}
