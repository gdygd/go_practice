package general

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

type SvgGrp struct {
	XMLName xml.Name `xml:"g"`

	Id             string  `xml:"id,attr,omitempty"`
	Stroke         string  `xml:"stroke,attr,omitempty"`
	StrokeWidth    float64 `xml:"stroke-width,attr,omitempty"`
	Fill           string  `xml:"fill,attr,omitempty"`
	Strokelinejoin string  `xml:"stroke-linejoin,attr,omitempty"`
	Strokelinecap  string  `xml:"stroke-linecap,attr,omitempty"`
	FontFamily     string  `xml:"font-family,attr,omitempty"`

	Paths    []Path    `xml:"path"`
	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`
	Images   []SImage  `xml:"image"`
	Texts    []Text    `xml:"text"`
}

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	// Width   int      `xml:"width,attr,omitempty"`
	// Height  int      `xml:"height,attr,omitempty"`
	Width  string `xml:"width,attr,omitempty"`
	Height string `xml:"height,attr,omitempty"`
	// Width   int    `xml:"width"`
	// Height  int    `xml:"height"`
	ViewBox string `xml:"viewBox,attr,omitempty"`
	Style   string `xml:"style,attr,omitempty"`
	//Style string `xml:"-"`

	Images []SImage `xml:"image"`
	Paths  []Path   `xml:"path"`

	Rects    []Rect    `xml:"rect"`
	Polygons []Polygon `xml:"polygon"`

	Texts []Text `xml:"text"`

	Grps []SvgGrp `xml:"g,omitempty"`
}

type Path struct {
	// M, L Z만 지원
	XMLName xml.Name `xml:"path"`
	Id      string   `xml:"id,attr"`
	Desc    string   `xml:"d,attr"`

	StrokeWidth    float64 `xml:"stroke-width,attr,omitempty"`
	Stroke         string  `xml:"stroke,attr,omitempty"`
	Strokelinejoin string  `xml:"stroke-linejoin,attr,omitempty"`
	Strokelinecap  string  `xml:"stroke-linecap,attr,omitempty"`
	Fill           string  `xml:"fill,attr,omitempty"`

	Pos [][2]float64 `xml:"-"`
}

type Text struct {
	XMLName    xml.Name `xml:"text"`
	Id         string   `xml:"id,attr"`
	FontSize   float64  `xml:"font-size,attr,omitempty"`
	FontFamily string   `xml:"font-family,attr,omitempty"`
	Textanchor string   `xml:"text-anchor,attr,omitempty"`
	FontWeight string   `xml:"font-weight,attr,omitempty"`
	Text       string   `xml:",innerxml"`
	Style      string   `xml:"style,attr,omitempty"`
	X          float64  `xml:"x,attr"`
	Y          float64  `xml:"y,attr"`
	Fill       string   `xml:"fill,attr"`

	TextLength   string `xml:"textLength,attr"`
	LengthAdjust string `xml:"lengthAdjust,attr"`

	Alignmentbaseline string `xml:"alignment-baseline,attr"`

	rgb [3]int `xml:"-"`
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
	rgb     [3]int   `xml:"-"`
}

type Polygon struct {
	XMLName   xml.Name     `xml:"polygon"`
	Id        string       `xml:"id,attr"`
	Style     string       `xml:"style,attr,omitempty,fill"`
	StrPoints string       `xml:"points,attr"`
	Points    [][2]float64 `xml:"-"`

	Fill   string `xml:"fill,attr"`
	Stroke string `xml:"stroke,attr"`

	Fillrgb [3]int `xml:"-"`
	Linergb [3]int `xml:"-"`
}

type SImage struct {
	XMLName xml.Name `xml:"image"`
	Id      string   `xml:"id,attr"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	// Width   float64  `xml:"width,attr"`
	// Height  float64  `xml:"height,attr"`
	Width               float64 `xml:"width,attr"`
	Height              float64 `xml:"height,attr"`
	Href                string  `xml:"href,attr"`
	PreserveAspectRatio string  `xml:"preserveAspectRatio,attr"`
}

//var ttfpath string = "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"
//var ttfpath string = "/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf"
//var ttfpath string = "/usr/share/fonts/opentype/noto/NotoSansCJK-Bold.ttc"

// var ttfpath string = "/home/theroad/notosans/NotoSansKR-Regular.ttf"
// var ttfpath string = "/home/theroad/notosans/laravel-notosans-korean/public/fonts/NotoSansKR-Bold.ttf"
var ttfpath string = "./font/noto/NotoSansKR-Bold.ttf"

func SetTTF(path string) {
	//ttfpath = path
}

func ViewboxToWH(viewbox string) (int, int) {

	//fmt.Printf("viewbox :%s\n", viewbox)

	viewinfo := strings.Split(viewbox, ",")
	if len(viewinfo) == 4 {
		//fmt.Printf("viewboxToWH:%s, %s\n", viewinfo[2], viewinfo[3])
		fval1, _ := strconv.ParseInt(viewinfo[2], 0, 64)
		fval2, _ := strconv.ParseInt(viewinfo[3], 0, 64)
		return int(fval1), int(fval2)
	}
	return 0, 0
}

func SvgpathtoPos(path string) [][2]float64 {

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

func DrawRect(r Rect, dc *gg.Context) {

	//dc.SetRGBA255(r[i].rgb[0], r[i].rgb[1], r[i].rgb[2], 255)
	dc.SetHexColor(r.Fill)

	//fmt.Printf("Rect X(%f) Y(%f) \n", r.X, r.Y)
	//fmt.Printf("Rect W(%f) H(%f) \n", r.Width, r.Height)

	dc.DrawRectangle(r.X, r.Y, r.Width, r.Height)

	dc.Fill()
}

func DrawText(t Text, dc *gg.Context) {

	//fmt.Printf("Text fill : [%s](%s)(%f)\n", t.Fill, t.Text, t.FontSize)
	dc.SetHexColor(t.Fill)

	// text
	//dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", 30)
	//dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoMono-Regular.ttf", t.FontSize)
	//dc.LoadFontFace("/home/theroad/notosans/NotoSansKR-Regular.ttf", t.FontSize)
	dc.LoadFontFace(ttfpath, t.FontSize)
	dc.DrawString(t.Text, t.X, t.Y)
	dc.Fill()
}

func DrawPath(p Path, width float64, stroke, fill string, dc *gg.Context) {

	//fmt.Printf("draw Path..\n")

	dc.Push()
	for j := 0; j < len(p.Pos); j++ {
		//fmt.Printf("drawPath point ###(%f, %f) \n", p.Pos[j][0], p.Pos[j][1])
		if j == 0 {
			//dc.MoveTo(p[i].Points[j][0], p[i].Points[j][1])
			dc.LineTo(p.Pos[j][0], p.Pos[j][1])
		} else {
			dc.LineTo(p.Pos[j][0], p.Pos[j][1])
		}
	}

	// fmt.Printf("polygonstroke : [%s]\n", p.Stroke)
	// fmt.Printf("polygonfill : [%s]\n", p.Fill)

	//dc.SetHexColor(p.Stroke)
	dc.SetHexColor(stroke)
	dc.SetLineWidth(5)

	dc.StrokePreserve()
	//dc.SetHexColor(p.Fill)
	dc.SetHexColor(fill)

	dc.Fill()
	dc.Pop()
}

func DrawPolygon(p Polygon, width float64, stroke, fill string, dc *gg.Context) {

	//fmt.Printf("POlygon fill, stroke : [%s](%s)\n", p.Fill, p.Stroke)

	dc.Push()
	for j := 0; j < len(p.Points); j++ {
		//fmt.Printf("drawPolygon point ###(%f, %f) \n", p.Points[j][0], p.Points[j][1])
		if j == 0 {
			//dc.MoveTo(p[i].Points[j][0], p[i].Points[j][1])
			dc.LineTo(p.Points[j][0], p.Points[j][1])
		} else {
			dc.LineTo(p.Points[j][0], p.Points[j][1])
		}
	}

	//fmt.Printf("polygonstroke : [%s]\n", p.Stroke)
	//fmt.Printf("polygonfill : [%s]\n", p.Fill)

	dc.SetHexColor(stroke)
	dc.SetLineWidth(width)

	dc.StrokePreserve()
	dc.SetHexColor(fill)

	dc.Fill()
	dc.Pop()
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

		if len(svginfo.Grps[i].Rects) > 0 {
			//drawRect
			for j := 0; j < len(svginfo.Grps[i].Rects); j++ {
				DrawRect(svginfo.Grps[i].Rects[j], dc)
			}
		}

		if len(svginfo.Grps[i].Paths) > 0 {
			//drawPath
			for j := 0; j < len(svginfo.Grps[i].Paths); j++ {
				w := svginfo.Grps[i].StrokeWidth
				stroke := svginfo.Grps[i].Stroke
				fill := svginfo.Grps[i].Paths[j].Fill
				DrawPath(svginfo.Grps[i].Paths[j], w, stroke, fill, dc)
			}
		}

		// Polygon
		if len(svginfo.Grps[i].Polygons) > 0 {
			for j := 0; j < len(svginfo.Grps[i].Polygons); j++ {
				w := svginfo.Grps[i].StrokeWidth
				stroke := svginfo.Grps[i].Stroke
				fill := svginfo.Grps[i].Polygons[j].Fill

				DrawPolygon(svginfo.Grps[i].Polygons[j], w, stroke, fill, dc)
			}
		}

		// text
		if len(svginfo.Grps[i].Texts) > 0 {
			for j := 0; j < len(svginfo.Grps[i].Texts); j++ {
				dc.Push()
				DrawText(svginfo.Grps[i].Texts[j], dc)
				dc.Pop()
			}
		}
	}
}

func SerializeSvg(s Svg) string {

	out, _ := xml.MarshalIndent(s, " ", "  ")
	strxml1 := string(out)
	strslice := strings.Split(strxml1, "\n")
	var svgstr string = ""
	for _, str := range strslice {
		svgstr += str
	}

	return svgstr
}

type PathInfo struct {
	Id             string
	StrokeWidth    float64
	Stroke         string
	Strokelinejoin string
	Strokelinecap  string
	Fill           string
	Zyn            bool

	Pos [][2]float64
}

type ImageInfo struct {
	Id                  string
	X                   float64
	Y                   float64
	Width               float64
	Height              float64
	Href                string // rawdata
	PreserveAspectRatio string
}

type TextInfo struct {
	Id                string
	FontFamily        string
	Textanchor        string
	FontSize          float64
	FontWeight        string
	X                 float64
	Y                 float64
	TextLength        string
	LengthAdjust      string
	Fill              string
	Text              string
	Alignmentbaseline string
}

func (s *Svg) Init(w, h int) {
	// <svg> 초기화
	// init width, height, viewbox
	// init Xmlns
	s.ViewBox = fmt.Sprintf("0,0,%d,%d", w, h)
	s.Width = fmt.Sprintf("%dpx", w)
	s.Height = fmt.Sprintf("%dpx", h)
	s.Xmlns = "http://www.w3.org/2000/svg"

	s.Paths = []Path{}
	s.Images = []SImage{}
	s.Texts = []Text{}

}

func (s *Svg) SetSvgStyle(style string) {
	s.Style = style
}

func (s *Svg) AddPath(info PathInfo) {
	var p Path
	p.Id = info.Id
	p.StrokeWidth = info.StrokeWidth
	p.Stroke = info.Stroke
	p.Strokelinejoin = info.Strokelinejoin
	p.Strokelinecap = info.Strokelinecap
	p.Fill = info.Fill

	var strPos string = ""
	for idx, pos := range info.Pos {
		if idx == 0 {
			strPos += fmt.Sprintf("M%.1f,%.1f", pos[0], pos[1])
		} else {
			if idx+1 == len(info.Pos) {
				if info.Zyn {
					strPos += fmt.Sprintf("L%.1f,%.1fZ", pos[0], pos[1])
				} else {
					strPos += fmt.Sprintf("L%.1f,%.1f", pos[0], pos[1])
				}

			} else {
				strPos += fmt.Sprintf("L%.1f,%.1f", pos[0], pos[1])
			}
		}
	}
	p.Desc = strPos

	s.Paths = append(s.Paths, p)

}

func (s *Svg) AddImage(info ImageInfo) {

	var img SImage
	img.Id = info.Id
	img.X = info.X
	img.Y = info.Y
	img.Width = info.Width
	img.Height = info.Height
	img.Href = info.Href
	img.PreserveAspectRatio = info.PreserveAspectRatio

	s.Images = append(s.Images, img)

}

func (s *Svg) AddText(info TextInfo) {
	var t Text
	t.Id = info.Id
	t.FontFamily = info.FontFamily
	t.Textanchor = info.Textanchor
	t.FontSize = info.FontSize
	t.FontWeight = info.FontWeight
	t.X = info.X
	t.Y = info.Y
	t.TextLength = info.TextLength
	t.LengthAdjust = info.LengthAdjust
	t.Fill = info.Fill
	t.Text = info.Text
	t.Alignmentbaseline = info.Alignmentbaseline

	s.Texts = append(s.Texts, t)
}
