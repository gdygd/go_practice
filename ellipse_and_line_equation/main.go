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

func main() {
	// 타원의 중심과 반경
	//h, k := 126.967067, 37.352215 // 타원의 중심 16m
	//h, k := 126.966882, 37.352232 // 타원의 중심 2m
	h, k := 126.966960, 37.352220 // 타원의 중심 5m
	// a, b := 1.13e-4, 8.98e-5      // 타원의 반경 (경도, 위도 방향)
	a, b := getLonRadius(8, k), 8/111320.0 // 타원의 반경 (경도, 위도 방향)

	// 직선의 방정식 y = mx + c
	x1, y1 := 126.966894, 37.352308
	x2, y2 := 126.966867, 37.352133

	// 기울기 계산
	m := (y2 - y1) / (x2 - x1)
	// 절편 c 계산
	c := y1 - m*x1

	// 이차방정식 계수 및 판별식 계산
	A, B, C, D := quadraticEquation(a, b, h, k, m, c)

	// 결과 출력
	fmt.Printf("이차방정식 계수:\nA = %.6f\nB = %.6f\nC = %.6f\n", A, B, C)
	fmt.Printf("판별식 D = %.6f\n", D)

	// 위치 관계 판별
	if D > 0 {
		fmt.Println("직선과 타원은 서로 다른 두 점에서 만난다.")
	} else if D == 0 {
		fmt.Println("직선과 타원은 한 점에서 만난다 (접선).")
	} else {
		fmt.Println("직선과 타원은 만나지 않는다.")
	}

	A1, B1, C1, D1 := quadraticEquation2(a, b, h, k, m, c)

	// 결과 출력
	fmt.Printf("이차방정식 계수:\nA = %.6f\nB = %.6f\nC = %.6f\n", A1, B1, C1)
	fmt.Printf("판별식 D = %.6f\n", D1)

	// 위치 관계 판별
	if D1 > 0 {
		fmt.Println("직선과 타원은 서로 다른 두 점에서 만난다.")
	} else if D1 == 0 {
		fmt.Println("직선과 타원은 한 점에서 만난다 (접선).")
	} else {
		fmt.Println("직선과 타원은 만나지 않는다.")
	}
}
