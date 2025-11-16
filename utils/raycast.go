package utils

import "math"

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func CalculateK(fovDegrees uint8) float32 {
	return float32(math.Tan(Radians(float64(fovDegrees) / 2)))
}

func calculateCameraX(screenX int, screenWidth int) float32 {
	return 2.0*float32(screenX)/float32(screenWidth) - 1.0
}

func GetRayDirection(playerDir Vector, k float32, screenWidth int, screenX int) Vector {
	perp := Perp(playerDir)
	plane := perp.Scale(k)
	cameraX := calculateCameraX(screenX, screenWidth)
	return playerDir.Add(plane.Scale(cameraX))
}
