package main

import (
	"fmt"

	"GoTicTacToe/utils"
)

func main() {
	var fov uint8 = 66
	screenWidth := 800
	playerDir := utils.Vector{X: 0.0, Y: 1.0}

	k := utils.CalculateK(fov)
	fmt.Printf("FOV: %d degrees, K: %.4f\n", fov, k)
	for i := 0; i <= screenWidth; i += screenWidth / 8 {
		rayDirection := utils.GetRayDirection(playerDir, k, screenWidth, i)
		fmt.Printf("ScreenX: %d, RayDir: (%.2f, %.2f)\n", i, rayDirection.X, rayDirection.Y)
	}
}
