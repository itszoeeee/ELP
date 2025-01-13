package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
)

const DIM = 20

const BLANK = 0
const UP = 1
const RIGHT = 2
const DOWN = 3
const LEFT = 4

// Définir le tableau de règles
var rules = [][][]int{
	// BLANK
	{
		{BLANK, UP},
		{BLANK, RIGHT},
		{BLANK, DOWN},
		{BLANK, LEFT},
	},
	// UP
	{
		{RIGHT, DOWN, LEFT},
		{UP, DOWN, LEFT},
		{BLANK, DOWN},
		{UP, RIGHT, DOWN},
	},
	// RIGHT
	{
		{RIGHT, DOWN, LEFT},
		{UP, DOWN, LEFT},
		{UP, RIGHT, LEFT},
		{BLANK, LEFT},
	},
	// DOWN
	{
		{BLANK, UP},
		{UP, DOWN, LEFT},
		{UP, RIGHT, LEFT},
		{UP, RIGHT, DOWN},
	},
	// LEFT
	{
		{RIGHT, DOWN, LEFT},
		{BLANK, RIGHT},
		{UP, RIGHT, LEFT},
		{UP, RIGHT, DOWN},
	},
}

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

// Fonction pour vérifier si l'option est valide
func checkValid(option *[]int, valid []int) {
	var newOption []int
	for _, elem := range *option {
		element := elem
		found := false
		// Vérifie si l'élément est dans le tableau valid
		for _, v := range valid {
			if element == v {
				found = true
				break
			}
		}
		// Si l'élément est valide, on l'ajoute à la nouvelle slice
		if found {
			newOption = append(newOption, elem)
		}
	}
	// Met à jour la slice d'origine
	*option = newOption
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
			collapsed: false,                               // Initialisé à false
			options:   []int{BLANK, UP, RIGHT, DOWN, LEFT}, // Options fixes
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

	var gridCopy []*gridItem
	gridCopy = grid
	var step int = len(gridCopy)
	var compteur int = 0
	var old_percent int = -1

	for len(gridCopy) > 1 { // Boucle principale
		// --- Faire une fonction ---
		// ----- Condition d'affichage (si la tuile est collapsed) -----
		for i := 0; i < DIM; i++ {
			for j := 0; j < DIM; j++ {
				var cell = grid[i+j*DIM]
				if cell.collapsed && len(cell.options) != 0 {
					var index = cell.options[0]
					placeImageInMatrix(outputImage, Tiles[index], i, j, cellSize)
				}
			}
		}
		// -------------------------
		// fmt.Println("Grille:")
		// for _, v := range grid {
		// 	fmt.Print(v.collapsed)
		// 	fmt.Println("	", v.options)
		// }

		gridCopy = nil // Réinitialiser gridCopy
		for _, cell := range grid {
			if !cell.collapsed {
				gridCopy = append(gridCopy, cell)
			}
		}

		// Trouver la longueur minimale de "options" dans la grille
		minLength := len(gridCopy[0].options)
		for _, v := range gridCopy {
			if len(v.options) < minLength {
				minLength = len(v.options)
			}
		}

		// Créer une nouvelle tranche pour les éléments ayant la longueur minimale
		var smallestItems []*gridItem
		for _, v := range gridCopy {
			if len(v.options) == minLength {
				smallestItems = append(smallestItems, v)
			}
		}

		// Choisir aléatoirement les éléments de plus petites longueurs
		// rand.Seed(time.Now().UnixNano())                                        // Initialiser le générateur de nombres aléatoires avec l'heure actuelle
		var randomItem *gridItem = smallestItems[rand.Intn(len(smallestItems))] // Sélectionner une clé aléatoire parmi celles disponibles
		randomItem.collapsed = true                                             // collapsed l'élément
		if len(randomItem.options) != 0 {                                       // vérifie que l'élèment c'est pas vide (erreur)
			var pick int = randomItem.options[rand.Intn(len(randomItem.options))] // choisir un option disponible (aléatoirement)
			randomItem.options = []int{pick}
		}

		// Création de la tuile suivante
		var nextGrid []*gridItem
		for i := 0; i < DIM; i++ {
			for j := 0; j < DIM; j++ {
				var index = j + i*DIM
				if grid[index].collapsed {
					nextGrid = append(nextGrid, grid[index])
				} else {
					var cell_option = []int{BLANK, UP, RIGHT, DOWN, LEFT}

					// Look up
					if i > 0 {
						var up = grid[j+(i-1)*DIM]
						var validOption []int
						for _, option := range up.options {
							var valid []int = rules[option][2]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look right
					if j < DIM-1 {
						var right = grid[j+1+i*DIM]
						var validOption []int
						for _, option := range right.options {
							var valid []int = rules[option][3]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look down
					if i < DIM-1 {
						var down = grid[j+(i+1)*DIM]
						var validOption []int
						for _, option := range down.options {
							var valid []int = rules[option][0]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look left
					if j > 0 {
						var left = grid[j-1+i*DIM]
						var validOption []int
						for _, option := range left.options {
							var valid []int = rules[option][1]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Crée une nouvelle instance de gridItem
					cell := &gridItem{
						collapsed: false,       // Initialisé à false
						options:   cell_option, // Options fixes
					}
					nextGrid = append(nextGrid, cell) // Ajouter l'élément à la grille
				}
			}
		}

		grid = nextGrid

		// Calcul du pourcentage
		var percent int = (compteur * 101) / step
		compteur++

		// Affichage de la barre de chargement
		if old_percent != percent {
			fmt.Printf("\r[") // permet de revenir au début de la ligne sans en ajouter une nouvelle, afin de mettre à jour la même ligne de la console.
			for j := 0; j < 50; j++ {
				if j < percent/2 {
					fmt.Print("=")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Printf("] %d%%", percent)
			old_percent = percent
		}
	}

	// --- Faire une fonction ---
	// ----- Condition d'affichage (si la tuile est collapsed) -----
	for i := 0; i < DIM; i++ {
		for j := 0; j < DIM; j++ {
			var cell = grid[i+j*DIM]
			if cell.collapsed && len(cell.options) != 0 {
				var index = cell.options[0]
				placeImageInMatrix(outputImage, Tiles[index], i, j, cellSize)
			}
		}
	}
	// -------------------------

	// ----- Exportation de l'image -----

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
