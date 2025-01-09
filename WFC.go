package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
)

const DIM = 2

const BLANK = 0
const UP = 1
const RIGHT = 2
const DOWN = 3
const LEFT = 4

type gridItem struct {
	collapsed bool  // Champ pour "collapsed"
	options   []int // Un tableau de tableaux d'entiers
}

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

func createTile() (Tiles [5]image.Image, err error) {
	// Charger l'image pattern.png
	for i := 0; i < 5; i++ {
		Tiles[i] = image.Transparent
	}

	Tiles[0], err = loadImage("pattern/blank.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'blank.png': %w", err)
	}
	Tiles[1], err = loadImage("pattern/up.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'up.png': %w", err)
	}
	Tiles[2], err = loadImage("pattern/right.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'right.png': %w", err)
	}
	Tiles[3], err = loadImage("pattern/down.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'down.png': %w", err)
	}
	Tiles[4], err = loadImage("pattern/left.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'left.png': %w", err)
	}
	return Tiles, nil
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

	// ----- Initialisation -----

	// Création des tuiles
	Tiles, err := createTile()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Création de la grille
	var grid []*gridItem

	// Initialisation des éléments de la grille
	for i := 0; i < DIM*DIM; i++ {
		// Crée une nouvelle instance de gridItem
		cell := &gridItem{
			collapsed: false,                // Initialisé à false
			options:   []int{1, 2, 3, 4, 5}, // Options fixes
		}
		grid = append(grid, cell) // Ajouter l'élément à la grille
	}

	// Dimensions de l'image de sortie
	gridWidth, gridHeight := DIM, DIM
	cellSize := Tiles[0].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

	// Créer de l'image de sortie
	outputImage := createEmptyImage(cellSize*gridWidth, cellSize*gridHeight)

	// ----- Fin de l'initialisation -----

	// ----- Boucle principale -----

	// TEST
	// grid[0].options = []int{BLANK, UP}
	// grid[1].options = []int{DOWN, LEFT}

	// // Trier la grille par la longueur de "options" (inutile)
	// sort.Slice(grid, func(i, j int) bool {
	// 	return len(grid[i].options) < len(grid[j].options)
	// })

	// Trouver la longueur minimale de "options" dans la grille
	minLength := len(grid[0].options)
	for _, v := range grid {
		if len(v.options) < minLength {
			minLength = len(v.options)
		}
	}

	// Créer une nouvelle tranche pour les éléments ayant la longueur minimale
	var smallestItems []*gridItem
	for _, v := range grid {
		if len(v.options) == minLength {
			smallestItems = append(smallestItems, v)
		}
	}

	// Choisir aléatoirement les éléments de plus petites longueurs
	var randomItem *gridItem = smallestItems[rand.Intn(len(smallestItems))] // Sélectionner une clé aléatoire parmi celles disponibles
	fmt.Printf("\nÉlément sélectionné : %d\n", randomItem.options)          // Afficher l'élément sélectionné
	randomItem.collapsed = true                                             // collapsed l'élément
	var pick int = rand.Intn(len(randomItem.options))                       // choisir un option disponible (aléatoirement)
	randomItem.options = []int{pick}

	// ----- Condition d'affichage (si la tuile est collapsed) -----
	for i := 0; i < DIM; i++ {
		for j := 0; j < DIM; j++ {
			var cell = grid[i+j*DIM]
			if cell.collapsed {
				var index = cell.options[0]
				placeImageInMatrix(outputImage, Tiles[index], i, j, cellSize)
			}
		}
	}

	// Afficher la nouvelle grille
	fmt.Println("Nouvelle grille:")
	for _, v := range grid {
		fmt.Println(v.options)
	}

	// Création de la tuile suivante
	// var nextTile []*gridItem
	for i := 0; i < DIM; i++ {
		for j := 0; j < DIM; j++ {
			// var index = grid[i+j*DIM]
		}
	}

	// ----- Exportation de l'image -----

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
