package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

const BLANK = 0
const T_UP = 1
const T_RIGHT = 2
const T_DOWN = 3
const T_LEFT = 4
const C_UP = 5
const C_RIGHT = 6
const C_DOWN = 7
const C_LEFT = 8
const CROSS = 9
const F_H = 10
const F_V = 11

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
	// Initialisation des tuiles
	var T, C, F image.Image
	for i := 0; i < 12; i++ {
		Tiles = append(Tiles, image.Transparent)
	}
	orients, images := lecture_json("input.JSON")
	Tiles[BLANK], err = loadImage(images[4])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'blank.png': %w", err)
	}
	Tiles[CROSS], err = loadImage(images[5])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'cross.png': %w", err)
	}

	T, err = loadImage(images[1])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 't_.png': %w", err)
	}
	Tiles[T_UP], Tiles[T_RIGHT], Tiles[T_DOWN], Tiles[T_LEFT], err = flipImage(T, orients[0])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	C, err = loadImage(images[2])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'c_.png': %w", err)
	}
	Tiles[C_UP], Tiles[C_RIGHT], Tiles[C_DOWN], Tiles[C_LEFT], err = flipImage(C, orients[1])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	F, err = loadImage(images[3])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'f_.png': %w", err)
	}
	_, Tiles[F_H], Tiles[F_V], _, err = flipImage(F, orients[2])
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
	cellSize := Tiles[BLANK].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

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
