package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// Fonction pour ouvrir une image et la décoder
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Fonction pour créer une image vide avec une taille spécifique (largeur, hauteur)
func createEmptyImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

// Fonction pour dessiner l'image dans une case spécifique (x, y) de la matrice
func placeImageInMatrix(dst *image.RGBA, src image.Image, gridX, gridY, cellSize int) {
	// Position de départ de la cellule
	startX := gridX * cellSize
	startY := gridY * cellSize

	// Dessine l'image src dans la cellule
	for y := 0; y < cellSize; y++ {
		for x := 0; x < cellSize; x++ {
			// Récupère la couleur de chaque pixel de l'image source
			color := src.At(x, y)
			// Place cette couleur dans l'image de destination
			dst.Set(startX+x, startY+y, color)
		}
	}
}

func main() {
	// Charger l'image pattern.png
	// blank, err := loadImage("pattern/blank.png")
	up, err := loadImage("pattern/up.png")
	// right, err := loadImage("pattern/right.png")
	// down, err := loadImage("pattern/down.png")
	// left, err := loadImage("pattern/left.png")
	if err != nil {
		fmt.Println("Erreur lors du chargement des images:", err)
		return
	}

	// Dimensions de la "grande image" (matrice 3x3)
	gridWidth, gridHeight := 5, 5
	cellSize := up.Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

	// Créer une image vide qui servira de "grande image"
	outputImage := createEmptyImage(cellSize*gridWidth, cellSize*gridHeight)

	// Choisir l'emplacement où on veut placer l'image (par exemple, case (1,1))
	// Case (0,0) est en haut à gauche, case (2,2) est en bas à droite.
	// Par exemple, on place l'image dans la case (1,1)
	// placeImageInMatrix(outputImage, blank, 1, 1, cellSize)
	placeImageInMatrix(outputImage, up, 1, 0, cellSize)
	// placeImageInMatrix(outputImage, right, 2, 1, cellSize)
	// placeImageInMatrix(outputImage, down, 1, 2, cellSize)
	// placeImageInMatrix(outputImage, left, 0, 1, cellSize)

	// Exporter l'image résultante dans un fichier PNG
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Erreur lors de la création de l'image de sortie:", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, outputImage)
	if err != nil {
		fmt.Println("Erreur lors de l'exportation de l'image:", err)
	}

	fmt.Println("Image exportée avec succès dans output.png")
}
