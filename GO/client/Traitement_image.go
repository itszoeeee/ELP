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

func createTile_Json() (Tiles []image.Image, err error) {
	// Charger l'image pattern.png
	for i := 0; i < 12; i++ {
		Tiles = append(Tiles, image.Transparent)
	}
	Tiles[0], err = loadImage("../pattern/blank.png")
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'blank.png': %w", err)
	}
	Tiles[9], err = loadImage("../pattern/cross.png")
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'cross.png': %w", err)
	}
	orients, fichiers := ouverture_json("input.JSON")
	Tiles[1], Tiles[2], Tiles[3], Tiles[4], err = flipImage(fichiers[0], orients[0])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	Tiles[5], Tiles[6], Tiles[7], Tiles[8], err = flipImage(fichiers[1], orients[1])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	_, Tiles[10], Tiles[11], _, err = flipImage(fichiers[2], orients[2])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	return Tiles, err
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

func display(grid [][]int, width, height int) {
	// Création des tuile
	Tiles, err := createTile_Json()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Dimensions de l'image de sortie
	gridWidth, gridHeight := width, height
	cellSize := Tiles[0].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

	// Créer de l'image de sortie
	outputImage := createEmptyImage(cellSize*gridWidth, cellSize*gridHeight)

	// Affichage des tuiles (si la tuile est collapsed)
	for j := 0; j < width; j++ {
		for i := 0; i < height; i++ {
			var cell = grid[j][i]
			if cell != -1 {
				index := cell
				placeImageInMatrix(outputImage, Tiles[index], i, j, cellSize)
			} else {
				placeImageInMatrix(outputImage, image.Transparent, i, j, cellSize)
			}
		}
	}

	// Exporter l'image résultante dans un fichier PNG
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("\n\n Erreur lors de la création de l'image de sortie:\n", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, outputImage)
	if err != nil {
		fmt.Println("\n\n Erreur lors de l'exportation de l'image:\n", err)
	}

	fmt.Println("\n\n Image exportée avec succès dans output.png\n")
}
