package main

import (
	"fmt"
)

type Point struct {
	Lat float64
	Lng float64
}

func PointInPolygon(point Point, polygon []Point) bool {
	n := len(polygon)
	inside := false

	j := n - 1
	for i := 0; i < n; i++ {
		pi := polygon[i]
		pj := polygon[j]

		intersect := ((pi.Lng > point.Lng) != (pj.Lng > point.Lng)) &&
			(point.Lat < (pj.Lat-pi.Lat)*(point.Lng-pi.Lng)/(pj.Lng-pi.Lng)+pi.Lat)

		if intersect {
			inside = !inside
		}
		j = i
	}

	return inside
}

func main() {

	polygon := []Point{
		{37.367726, 126.723248},
		{37.367965, 126.723467},
		{37.365224, 126.728137},
		{37.365028, 126.727740},
	}

	// 테스트할 점
	// testPoint := Point{37.3665, 126.7250}	// outside
	testPoint := Point{37.366678, 126.725130} // inside

	if PointInPolygon(testPoint, polygon) {
		fmt.Println("다각형 내부")
	} else {
		fmt.Println("다각형 외부")
	}
}
