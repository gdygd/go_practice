package main

import (
	"fmt"
	"math"
)

// 타원과 직선의 관계 테스트 완료.

// 이차방정식의 계수를 계산하고 판별식을 구하는 함수
func quadraticEquation(a, b, h, k, m, c float64) (float64, float64, float64, float64) {
	// 이차방정식의 계수
	A := a*a*m*m + b*b
	B := 2*a*a*m*(c-k) - 2*b*b*h
	C := a*a*(c-k)*(c-k) + b*b*h*h - a*a*b*b

	// 판별식 계산
	D := B*B - 4*A*C

	return A, B, C, D
}

// 이차방정식의 계수를 계산하고 판별식을 구하는 함수
func quadraticEquation2(a, b, h, k, m, c float64) (float64, float64, float64, float64) {
	// 이차방정식의 계수
	A := a*a*m*m + b*b
	B := 2*a*a*m*(c-k) - 2*b*b*h
	C := b*b*h*h + a*a*c*c - 2*a*a*k*c + a*a*k*k - a*a*b*b

	// 판별식 계산
	D := B*B - 4*A*C

	return A, B, C, D
}

func getLonRadius(meter int, lat float64) float64 {
	// 경도반지름
	angle := lat * math.Pi / 180.0 // To degree
	lon1dist := 111320 * math.Cos(angle)
	dist := (1 / lon1dist) * float64(meter)
	return dist
}

func isInsideEllipse(x, y, h, k, a, b float64) bool {
	// 타원 내부 판별
	value := ((x-h)*(x-h))/(a*a) + ((y-k)*(y-k))/(b*b)
	return value < 1
}

func main() {
	// 타원의 중심과 반경
	h, k := 126.966960, 37.352220          // 타원의 중심 5m
	a, b := getLonRadius(8, k), 8/111320.0 // 타원의 반경 (경도, 위도 방향)

	// 직선의 방정식 y = mx + c
	x, y := 126.966894, 37.352308

	if isInsideEllipse(x, y, h, k, a, b) {
		fmt.Printf("inside Ellipse \n")
	} else {
		fmt.Printf("is not inside Ellipse \n")
	}
}
