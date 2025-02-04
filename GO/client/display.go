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
func load_image(path string) (image.Image, error) {
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

func create_tile() (Tiles []image.Image, err error) {
	// Initialisation des tuiles
	var T, C, F image.Image
	for i := 0; i < 12; i++ {
		Tiles = append(Tiles, image.Transparent)
	}

	orients, images := lecture_json("input.JSON")

	Tiles[BLANK], err = load_image(images[3])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'blank.png': %w", err)
	}

	Tiles[CROSS], err = load_image(images[4])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'cross.png': %w", err)
	}

	T, err = load_image(images[0])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 't_.png': %w", err)
	}
	Tiles[T_UP], Tiles[T_RIGHT], Tiles[T_DOWN], Tiles[T_LEFT], err = flipImage(T, orients[0])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	C, err = load_image(images[1])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'c_.png': %w", err)
	}
	Tiles[C_UP], Tiles[C_RIGHT], Tiles[C_DOWN], Tiles[C_LEFT], err = flipImage(C, orients[1])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	F, err = load_image(images[2])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors du chargement de l'image 'f_.png': %w", err)
	}
	_, Tiles[F_H], Tiles[F_V], _, err = flipImage(F, orients[2])
	if err != nil {
		return Tiles, fmt.Errorf("erreur lors de la premiere rotation: %w", err)
	}

	return Tiles, err
}

// Fonction pour dessiner l'image dans une case spécifique (x, y) de la matrice
func place_image(dst *image.RGBA, src image.Image, gridX, gridY, cellSize int) {
	// Position de départ de la cellule
	startX := gridX * cellSize
	startY := gridY * cellSize

	// Dessine l'image src dans la cellule
	for y := 0; y < cellSize; y++ {
		for x := 0; x < cellSize; x++ {
			color := src.At(x, y)              // Récupère la couleur de chaque pixel de l'image source
			dst.Set(startX+x, startY+y, color) // Place cette couleur dans l'image de destination
		}
	}
}

func display(grid [][]int, dim_x, dim_y int, erreur *bool) {
	// Création des tuile
	Tiles, err := create_tile()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Dimensions de l'image de sortie
	cellSize := Tiles[BLANK].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées
	// Créer de l'image de sortie
	outputImage := image.NewRGBA(image.Rect(0, 0, cellSize*dim_x, cellSize*dim_y))

	// Affichage des tuiles (si la tuile est collapsed)
	for j := 0; j < dim_y; j++ {
		for i := 0; i < dim_x; i++ {
			var cell = grid[j][i]
			if cell != -1 {
				place_image(outputImage, Tiles[cell], i, j, cellSize)
			} else {
				*erreur = true
				place_image(outputImage, image.Transparent, i, j, cellSize)
			}
		}
	}

	// Exporter l'image résultante dans un fichier PNG
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Println("Erreur lors de la création de l'image de sortie :", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, outputImage)
	if err != nil {
		fmt.Println("Erreur lors de l'exportation de l'image :", err)
	}

	fmt.Println("Image exportée avec succès dans output.png")
}
