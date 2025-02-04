package main

import (
	"fmt"
	"image"
)

var Orientations = map[int][]int{ // Équivalent d'un dico pour faire correspondre une orientation a une structure Orientation
	1: {1, 2, 3, 4},
	2: {2, 3, 4, 1},
	3: {3, 4, 1, 2},
	4: {4, 1, 2, 3},
}

func flipImage(img image.Image, orient int) (image.Image, image.Image, image.Image, image.Image, error) {
	// Rotation des images en fonction de param (1 pour une rotation, 2 pour deux rotations, ...)
	if Orientations[orient] == nil {
		return nil, nil, nil, nil, fmt.Errorf("orientation invalide")
	}

	// Trouver la liste des orientations à partir de l'orientation actuelle
	orientationList := Orientations[orient]
	result := make([]image.Image, 4)
	var rotated image.Image

	result[orientationList[0]-1] = img
	rotated = rotate90(img)
	result[orientationList[1]-1] = rotated
	rotated = rotate90(rotated)
	result[orientationList[2]-1] = rotated
	rotated = rotate90(rotated)
	result[orientationList[3]-1] = rotated

	return result[0], result[1], result[2], result[3], nil
}

// Fonction pour effectuer une rotation de 90 degrés vers la droite
func rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newX := bounds.Max.Y - y - 1
			newY := x
			rotated.Set(newX, newY, img.At(x, y))
		}
	}

	return rotated
}
