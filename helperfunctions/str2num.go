package helperfunctions

import (
	"fmt"
	"math"
)

const pi = math.Pi
const d2r = pi / 180
const r2d = 180 / pi

//Str2Num : convert string to float64
func Str2Num(val string) float64 {
	var v float64
	fmt.Sscanf(val, "%f", &v)
	return v
}

//Deg2Rad : convert deg to radian
func Deg2Rad(val float64) float64 {
	return val * d2r
}

//Rad2Deg : convert radian to deg
func Rad2Deg(val float64) float64 {
	return val * r2d
}
