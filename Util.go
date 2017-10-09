package main

import "math"

func distanceBetween(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func lengthOf(x float64, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func normalize(x float64, y float64) (float64, float64) {
	length := lengthOf(x, y)
	if length != 0 {
		return x / length, y / length
	}

	return x, y
}

func limitToLength(x float64, y float64, len float64) (float64, float64) {
	oldLen := lengthOf(x, y)

	if (len > oldLen) {
		return x, y
	}

	newX, newY := normalize(x, y)

	return newX * len, newY * len
}
