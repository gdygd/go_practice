package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"

	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type SvgGrp struct {
	XMLName xml.Name `xml:"g"`

	Id             string  `xml:"id,attr"`
	Stroke         string  `xml:"stroke,attr"`
	StrokeWidth    float64 `xml:"stroke-width,attr"`
	Fill           string  `xml:"fill,attr"`
	Strokelinejoin string  `xml:"stroke-linejoin,attr"`
	Strokelinecap  string  `xml:"stroke-linecap,attr"`
	FontFamily     string  `xml:"font-family,attr,omitempty"`

	Paths    []Path    `xml:"path"`
	Texts    []Text    `xml:"text"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
	Images   []SImage  `xml:"image"`
}

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   int      `xml:"width,attr,omitempty"`
	Height  int      `xml:"height,attr,omitempty"`
	ViewBox string   `xml:"viewBox,attr,omitempty"`

	Grps     []SvgGrp  `xml:"g"`
	Paths    []Path    `xml:"path"`
	Texts    []Text    `xml:"text"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
	Images   []SImage  `xml:"image"`
}

type Path struct {
	// M, L Z만 지원
	XMLName xml.Name `xml:"path"`
	Id      string   `xml:"id,attr"`
	Desc    string   `xml:"d,attr"`
	Pos     [][2]float64
}

type Text struct {
	XMLName    xml.Name `xml:"text"`
	Id         string   `xml:"id,attr"`
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
	Id      string   `xml:"id,attr"`
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
	Id        string   `xml:"id,attr"`
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
	Id      string   `xml:"id,attr"`
	Href    string   `xml:"href,attr"`
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

func drawRect(r Rect, dc *gg.Context) {

	//dc.SetRGBA255(r[i].rgb[0], r[i].rgb[1], r[i].rgb[2], 255)
	dc.SetHexColor(r.Fill)

	fmt.Printf("Rect X(%f) Y(%f) \n", r.X, r.Y)
	fmt.Printf("Rect W(%f) H(%f) \n", r.Width, r.Height)

	dc.DrawRectangle(r.X, r.Y, r.Width, r.Height)

	dc.Fill()
}

func drawText(t Text, dc *gg.Context) {

	fmt.Printf("Text fill : [%s](%s)(%f)\n", t.Fill, t.Text, t.FontSize)
	dc.SetHexColor(t.Fill)

	// text
	//dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", 30)
	//dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", t.FontSize)
	//dc.LoadFontFace("/home/theroad/notosans/NotoSansKR-Regular.ttf", t.FontSize)
	dc.LoadFontFace("/home/theroad/notosans/NotoSansKR-Regular.ttf", t.FontSize)
	dc.DrawString(t.Text, t.X, t.Y)
	//dc.DrawString("abcd", t.X, t.Y)
	dc.Fill()
}

func drawPath(p Path, width float64, stroke, fill string, dc *gg.Context) {

	fmt.Printf("draw Path..\n")

	dc.Push()
	for j := 0; j < len(p.Pos); j++ {
		fmt.Printf("drawPath point ###(%f, %f) \n", p.Pos[j][0], p.Pos[j][1])
		if j == 0 {
			//dc.MoveTo(p[i].Points[j][0], p[i].Points[j][1])
			dc.LineTo(p.Pos[j][0], p.Pos[j][1])
		} else {
			dc.LineTo(p.Pos[j][0], p.Pos[j][1])
		}
	}

	// fmt.Printf("polygonstroke : [%s]\n", p.Stroke)
	// fmt.Printf("polygonfill : [%s]\n", p.Fill)

	dc.SetHexColor(stroke)
	//dc.SetHexColor("#00FF00")
	//dc.SetLineWidth(width)
	dc.SetLineWidth(10)

	dc.StrokePreserve()
	dc.SetHexColor(fill)
	//dc.SetHexColor("F0F0F0")

	dc.Fill()
	dc.Pop()
}

//func drawPolygon(p Polygon, dc *gg.Context) {
func drawPolygon(p Polygon, width float64, stroke, fill string, dc *gg.Context) {

	fmt.Printf("POlygon fill, stroke : [%s](%s)\n", p.Fill, p.Stroke)

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

	fmt.Printf("polygonstroke : [%s]\n", p.Stroke)
	fmt.Printf("polygonfill : [%s]\n", p.Fill)

	dc.SetHexColor(stroke)
	dc.SetLineWidth(10)

	dc.StrokePreserve()
	dc.SetHexColor(fill)

	dc.Fill()
	dc.Pop()
}

func addLabel(img *image.RGBA, x, y int, label string) {
	// userface := basicfont.Face7x13
	// userface.Width = 6
	// userface.Ascent = 30

	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		//Face: userface,
		Dot: point,
	}
	d.DrawString(label)
}

func drawImage2(simg SImage, dc *gg.Context) {
	fmt.Printf("drawImage2....\n")

	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	//img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	addLabel(img, 20, 30, "Hello Go")

	dc.DrawImage(img, 0, 0)
}

func drawImage(simg SImage, dc *gg.Context) {

	// img, errimg := gg.LoadImage(simg.Href)
	// if errimg != nil {
	// 	fmt.Printf("errimg : %v", errimg)
	// }

	// dc.DrawImage(img, 0, 0)

	strhrefs := strings.Split(simg.Href, "base64,")

	var strbase64 string = ""
	strbase64 = strhrefs[1]
	fmt.Printf("base64str : %s\n", strbase64)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(strbase64))
	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("drawImage decoding err : %v", err)
		return
	}
	dc.DrawImage(m, 0, 0)
}

func Drawsvg(svginfo Svg, dc *gg.Context) {
	//func drawText(t Text, dc *gg.Context) {
	//func drawPolygon(p Polygon, dc *gg.Context) {

	for i := 0; i < len(svginfo.Grps); i++ {
		//fmt.Printf("rect : %v\n", svginfo.Grps[i].Rects)
		if len(svginfo.Grps[i].Rects) > 0 {
			//drawRect
			for j := 0; j < len(svginfo.Grps[i].Rects); j++ {
				drawRect(svginfo.Grps[i].Rects[j], dc)
			}

		}

		if len(svginfo.Grps[i].Paths) > 0 {
			//drawPath
			for j := 0; j < len(svginfo.Grps[i].Paths); j++ {
				w := svginfo.Grps[i].StrokeWidth
				stroke := svginfo.Grps[i].Stroke
				fill := svginfo.Grps[i].Fill
				drawPath(svginfo.Grps[i].Paths[j], w, stroke, fill, dc)
			}

		}

		// Polygon
		if len(svginfo.Grps[i].Polygons) > 0 {
			for j := 0; j < len(svginfo.Grps[i].Polygons); j++ {
				w := svginfo.Grps[i].StrokeWidth
				stroke := svginfo.Grps[i].Stroke
				fill := svginfo.Grps[i].Fill

				drawPolygon(svginfo.Grps[i].Polygons[j], w, stroke, fill, dc)
			}
		}

		// text
		if len(svginfo.Grps[i].Texts) > 0 {
			for j := 0; j < len(svginfo.Grps[i].Texts); j++ {
				dc.Push()
				drawText(svginfo.Grps[i].Texts[j], dc)
				dc.Pop()
			}
		}

		// image
		if len(svginfo.Grps[i].Images) > 0 {
			for j := 0; j < len(svginfo.Grps[i].Images); j++ {
				dc.Push()
				//drawText(svginfo.Grps[i].Texts[j], dc)
				//drawImage(svginfo.Grps[i].Images[j], dc)
				drawImage2(svginfo.Grps[i].Images[j], dc)
				dc.Pop()
			}
		}
	}
}

func viewboxToWH(viewbox string) (int, int) {

	fmt.Printf("viewbox :%s\n", viewbox)

	viewinfo := strings.Split(viewbox, ",")
	if len(viewinfo) == 4 {
		fmt.Printf("viewboxToWH:%s, %s\n", viewinfo[2], viewinfo[3])
		fval1, _ := strconv.ParseInt(viewinfo[2], 0, 64)
		fval2, _ := strconv.ParseInt(viewinfo[3], 0, 64)
		return int(fval1), int(fval2)
	}
	return 0, 0
}

func svgpathtoPos(path string) [][2]float64 {
	//var fpos [][2]float64 = [][2]float64{}

	//path := `M40,10L50,10L80,10L80,7L90,15L80,23L80,20L55,20L40,20L40,10`

	posinfo := strings.Split(path, "L")
	var fpos [][2]float64 = [][2]float64{}

	for _, pos := range posinfo {
		var strXY string = ""

		if (pos[0] >= 0x41 && pos[0] <= 0x5A) ||
			(pos[0] >= 0x61 && pos[0] <= 0x7A) {

			strXY = pos[1:]
		} else {
			strXY = pos[:]
		}
		// fmt.Printf("strXY : %s \n", strXY)
		xy := strings.Split(strXY, ",")
		fx, _ := strconv.ParseFloat(xy[0], 0)
		fy, _ := strconv.ParseFloat(xy[1], 0)
		// fmt.Printf("XY : %f, %f \n", fx, fy)

		fpos = append(fpos, [2]float64{fx, fy})
	}
	//fmt.Printf("fpos : %v\n", fpos)

	return fpos
}

func svgtoimg() {
	// Open the SVG file
	// svgFile, err := os.Open("rect.svg")
	// if err != nil {
	// 	fmt.Println("Error opening SVG file:", err)
	// 	return
	// }
	// defer svgFile.Close()

	// // Create a new image context
	// dc := gg.NewContext(0, 0)
	// err = dc.LoadSVG(svgFile)
	// if err != nil {
	// 	fmt.Println("Error loading SVG:", err)
	// 	return
	// }

	// // Save the image as BMP
	// err = dc.SavePNG("svgrect.png")
	// if err != nil {
	// 	fmt.Println("Error saving BMP:", err)
	// 	return
	// }

	// fmt.Println("SVG converted to BMP successfully!")

}

func main() {
	svgtoimg()
	// svgpathtoPos()

	// Open our xmlFile
	xmlFile, err := os.Open("./pathsvg__.svg")
	//xmlFile, err := os.Open("./rect_text_copy.svg")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var svginfo Svg
	xml.Unmarshal(byteValue, &svginfo)
	// fmt.Printf("\n\n%s\n\n", string(byteValue))
	// fmt.Printf("svr info : [%v]\n\n", svginfo)
	//fmt.Printf("svr grp info : [%v]\n\n", svginfo.Grps)
	//fmt.Printf("svr Polygons info : [%v]\n\n", svginfo.Polygons)

	w, h := viewboxToWH(svginfo.ViewBox)
	fmt.Printf("wh : %d, %d\n", w, h)
	svginfo.Width = w
	svginfo.Height = h

	fmt.Printf("grp))================================================>>\n")
	for i := 0; i < len(svginfo.Grps); i++ {
		// path
		fmt.Printf("PATH))>>\n")
		if len(svginfo.Grps[i].Paths) > 0 {
			fmt.Printf("Path))================================================>>\n")
			for j := 0; j < len(svginfo.Grps[i].Paths); j++ {
				fmt.Printf("xmlname: %v\n", svginfo.Grps[i].Paths[j].XMLName)
				fmt.Println("Path Id: " + svginfo.Grps[i].Paths[j].Id)
				fmt.Println("Path Desc: " + svginfo.Grps[i].Paths[j].Desc)
				//svgpathtoPos(svginfo.Paths[i].Desc)
				fpos := svgpathtoPos(svginfo.Grps[i].Paths[j].Desc)
				svginfo.Grps[i].Paths[j].Pos = [][2]float64{}
				for _, pos := range fpos {
					svginfo.Grps[i].Paths[j].Pos = append(svginfo.Grps[i].Paths[j].Pos, [2]float64{pos[0], pos[1]})
				}
				fmt.Printf("path pos : [%v]\n", svginfo.Grps[i].Paths[j].Pos)
			}
			fmt.Printf("======================================================>>\n")
		}

		// text
		if len(svginfo.Grps[i].Texts) > 0 {
			fmt.Printf("TEXT))>>\n")
			for j := 0; j < len(svginfo.Grps[i].Texts); j++ {
				fmt.Printf("xmlname: %v\n", svginfo.Grps[i].Texts[j].XMLName)
				fmt.Printf("Text Id:%s \n", svginfo.Grps[i].Texts[j].Id)
				//fmt.Println("Text Type: " + svginfo.Texts[i]. FontSize)
				fmt.Println("Text font: " + svginfo.Grps[i].Texts[j].FontFamily)
				fmt.Println("Text text: " + svginfo.Grps[i].Texts[j].Text)
				fmt.Println("Text style: " + svginfo.Grps[i].Texts[j].Style)
				fmt.Println("Text fill: " + svginfo.Grps[i].Texts[j].Fill)
				fmt.Printf("Text X(%f) Y(%f) \n", svginfo.Grps[i].Texts[j].X, svginfo.Grps[i].Texts[j].Y)

				//svginfo.Grps[i].Texts[i].rgb = ColorHex2RGB(svginfo.Grps[i].Texts[j].Fill)
			}
			fmt.Printf("======================================================>>\n")
		}

		// Polygon
		if len(svginfo.Grps[i].Polygons) > 0 {
			fmt.Printf("POLYGON))>>\n")

			for j := 0; j < len(svginfo.Grps[i].Polygons); j++ {
				//fmt.Printf("xmlname: %v\n", svginfo.Grps[i].Polygons[j].XMLName)
				fmt.Printf("Polygon Id:%s \n", svginfo.Grps[i].Polygons[j].Id)
				fmt.Println("Polygon style: " + svginfo.Grps[i].Polygons[j].Style)
				// fmt.Printf("Polygon Fillrgb:%s ", svginfo.Grps[i].Polygons[j].Fillrgb)
				// fmt.Printf("Polygon Linergb:%s ", svginfo.Grps[i].Polygons[j].Linergb)
				fmt.Printf("Polygon points(%s) \n", svginfo.Grps[i].Polygons[j].StrPoints)

				//svginfo.Grps[i].Polygons[i].Fillrgb = ColorHex2RGB(svginfo.Polygons[i].Fill)
				//svginfo.Grps[i].Polygons[i].Linergb = ColorHex2RGB(svginfo.Polygons[i].Stroke)

				svginfo.Grps[i].Polygons[j].Points = [][2]float64{}

				// points parsing
				//strtmp := "alfa bravo charlie delta echo foxtrot golf"
				points := strings.Split(svginfo.Grps[i].Polygons[j].StrPoints, " ")
				for _, strXY := range points {
					// strXY : x,y
					xy := strings.Split(strXY, ",")
					fmt.Printf("X[%s], Y[%s]\n", xy[0], xy[1])
					nx, _ := strconv.ParseFloat(xy[0], 2)
					ny, _ := strconv.ParseFloat(xy[1], 2)

					svginfo.Grps[i].Polygons[j].Points = append(svginfo.Grps[i].Polygons[j].Points, [2]float64{nx, ny})

				}

				// // polygon close를 위해 첫번째 점 마지막 인덱스에 추가
				// if len(points) > 0 {
				// 	svginfo.Grps[i].Polygons[j].Points = append(svginfo.Grps[i].Polygons[j].Points, svginfo.Grps[i].Polygons[j].Points[0])
				// }

				//fmt.Printf("Polygon point arr(%v) \n", svginfo.Grps[i].Polygons[j].Points)
			}
			fmt.Printf("======================================================>>\n")

		}

		// Rect
		if len(svginfo.Grps[i].Rects) > 0 {
			fmt.Printf("RECT))>>\n")
			for j := 0; j < len(svginfo.Grps[i].Rects); j++ {
				//fmt.Printf("xmlname: %v\n", svginfo.Grps[i].Texts[j].XMLName)
				fmt.Printf("Text Id:%s \n", svginfo.Grps[i].Rects[j].Id)
				fmt.Printf("Rect X(%f) Y(%f) \n", svginfo.Grps[i].Rects[j].X, svginfo.Grps[i].Rects[j].Y)
				fmt.Printf("Rect W(%f) H(%f) \n", svginfo.Grps[i].Rects[j].Width, svginfo.Grps[i].Rects[j].Height)
				fmt.Println("Rect fill: " + svginfo.Grps[i].Rects[j].Fill)

				//svginfo.Grps[i].Texts[i].rgb = ColorHex2RGB(svginfo.Grps[i].Texts[j].Fill)
			}
			fmt.Printf("======================================================>>\n")
		}

	}

	fmt.Printf("h:%d w:%d\n", svginfo.Width, svginfo.Height)
	// svg size
	dc := gg.NewContext(svginfo.Width, svginfo.Height)
	Drawsvg(svginfo, dc)
	dc.SavePNG("svgbmp.bmp")
}
