package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
)

func giftest1() {
	// 각 프레임을 저장할 배열 생성
	var images []*image.Paletted
	var delays []int

	// 각 프레임을 생성하고 images 배열에 추가
	// 예시로 빨간색과 파란색 프레임을 추가합니다.
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	// 첫 번째 프레임
	frame1 := image.NewPaletted(image.Rect(0, 0, 100, 100), color.Palette{red})
	images = append(images, frame1)
	delays = append(delays, 10) // 10ms 지연

	// 두 번째 프레임
	frame2 := image.NewPaletted(image.Rect(0, 0, 100, 100), color.Palette{blue})
	images = append(images, frame2)
	delays = append(delays, 10) // 10ms 지연

	// 파일 생성
	file, err := os.Create("animation.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// gif 이미지 인코딩 및 파일 저장
	gif.EncodeAll(file, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func giftest2() {
	fileNames := []string{"img1.png", "img2.png", "img3.png"} // png 파일의 경로 및 이름을 적절히 수정하세요.
	//fileNames := []string{"out.png", "out2.png"} // png 파일의 경로 및 이름을 적절히 수정하세요.
	images := []*image.Paletted{}
	delays := []int{50, 100, 500}

	for _, fileName := range fileNames {
		reader, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer reader.Close()

		img, err := png.Decode(reader)
		if err != nil {
			log.Fatalf("Error decoding PNG: %v", err)
		}

		bounds := img.Bounds()

		drawer := draw.FloydSteinberg
		palettedImg := image.NewPaletted(bounds, palette.Plan9)
		drawer.Draw(palettedImg, img.Bounds(), img, image.ZP)

		images = append(images, palettedImg)
		//delays = append(delays, 100) // 각 프레임 간의 딜레이를 조정하세요.
	}

	outputFile, err := os.Create("animationfrompng.gif") // 생성될 gif 파일의 경로 및 이름을 적절히 수정하세요.
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = gif.EncodeAll(outputFile, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		log.Fatalf("Error encoding GIF: %v", err)
	}

	fmt.Println("GIF 생성이 완료되었습니다.")

}

func main() {
	fmt.Printf("gif test..")

	giftest1()

	giftest2()
}
