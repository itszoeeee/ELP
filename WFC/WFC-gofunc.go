package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
)

const DIM_X = 4
const DIM_Y = 3

const BLANK = 0
const T_UP = 1
const T_RIGHT = 2
const T_DOWN = 3
const T_LEFT = 4
const C_UP = 5
const C_RIGHT = 6
const C_DOWN = 7
const C_LEFT = 8

// Définir le tableau de règles
var rules = [][][]int{
	// BLANK
	{
		{BLANK, T_UP, C_UP, C_RIGHT},      // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN}, // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},   // south
		{BLANK, T_LEFT, C_UP, C_LEFT},     // west
	},
	// T_UP
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP},              // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},           // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN},  // west
	},
	// T_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT},      // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT},    // south
		{BLANK, T_LEFT, C_UP, C_LEFT},             // west
	},
	// T_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT},             // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT},     // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN}, // west
	},
	// T_LEFT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},         // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT},    // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN},  // west
	},
	// C_UP
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},         // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},           // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN},  // west
	},
	// C_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT},      // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},           // south
		{BLANK, T_LEFT, C_UP, C_LEFT},             // west
	},
	// C_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT},           // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT},   // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT}, // south
		{BLANK, T_LEFT, C_UP, C_LEFT},          // west
	},
	// C_LEFT
	{
		{BLANK, T_UP, C_UP, C_RIGHT},             // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},        // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN}, // west
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

func createTile() (Tiles []image.Image, err error) {
	// Charger l'image pattern.png
	for k := 0; k < len(rules); k++ {
		Tiles = append(Tiles, image.Transparent)
	}

	Tiles[0], err = loadImage("pattern/blank.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'blank.png': %w", err)
	}
	Tiles[1], err = loadImage("pattern/t_up.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_up.png': %w", err)
	}
	Tiles[2], err = loadImage("pattern/t_right.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_right.png': %w", err)
	}
	Tiles[3], err = loadImage("pattern/t_down.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_down.png': %w", err)
	}
	Tiles[4], err = loadImage("pattern/t_left.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_left.png': %w", err)
	}
	Tiles[5], err = loadImage("pattern/c_up.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_up.png': %w", err)
	}
	Tiles[6], err = loadImage("pattern/c_right.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_right.png': %w", err)
	}
	Tiles[7], err = loadImage("pattern/c_down.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_down.png': %w", err)
	}
	Tiles[8], err = loadImage("pattern/c_left.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_left.png': %w", err)
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

func client(grid [][]*gridItem) {
	// Création des tuiles
	Tiles, err := createTile()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Dimensions de l'image de sortie
	gridWidth, gridHeight := DIM_X, DIM_Y
	cellSize := Tiles[0].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

	// Créer de l'image de sortie
	outputImage := createEmptyImage(cellSize*gridWidth, cellSize*gridHeight)

	// Affichage des tuiles (si la tuile est collapsed)
	for j := 0; j < DIM_Y; j++ {
		for i := 0; i < DIM_X; i++ {
			var cell = grid[j][i]
			if cell.collapsed && len(cell.options) != 0 {
				placeImageInMatrix(outputImage, Tiles[cell.options[0]], i, j, cellSize)
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

func grid_init(grid *[][]*gridItem) {
	// Initialiser chaque ligne de la grille
	for j := 0; j < DIM_Y; j++ {
		var row []*gridItem          // Créer un slice vide pour chaque ligne
		for i := 0; i < DIM_X; i++ { // Initialiser chaque cellule dans la ligne
			cell := &gridItem{ // Créer une nouvelle instance de gridItem
				collapsed: false, // Initialisé à false
				// Options par défaut
				options: []int{BLANK, T_UP, T_RIGHT, T_DOWN, T_LEFT, C_UP, C_RIGHT, C_DOWN, C_LEFT},
			}
			row = append(row, cell) // Ajouter la cellule à la ligne
		}
		*grid = append(*grid, row) // Ajouter la ligne à la grille
	}
}

func WFC(grid *[][]*gridItem) {
	//defer wg.Done() // Indiquer que la tâche est terminée (cette ligne sera exécutée juste avant que la fonction retourne et donc garantie de tuer le process)

	println(len((*grid)[0]))
	println(len(*grid))

	gridCopy := *grid
	// var step int = len(gridCopy)
	// var compteur int = 0
	// var old_percent int = -1

	for k := 0; k < len((*grid)[0])*len(*grid); k++ { // Boucle principale
		//for i := 0; i < 2; i++ {
		// fmt.Println("Grille:")
		// for _, row := range *grid {
		// 	for _, v := range row {
		// 		fmt.Print(v.collapsed)
		// 		fmt.Println("	", v.options)
		// 	}
		// }

		// Trouver la longueur minimale de "options" dans la grille copiée
		minLength := len(rules) // Il s'agit du maximum d'options possibles
		for _, row := range gridCopy {
			for _, v := range row {
				if len(v.options) < minLength {
					minLength = len(v.options)
				}
			}
		}

		// Créer une nouvelle tranche pour les éléments ayant la longueur minimale
		var smallestItems []*gridItem
		for _, row := range gridCopy {
			for _, v := range row {
				if len(v.options) == minLength {
					smallestItems = append(smallestItems, v)
				}
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

		gridCopy = nil             // Réinitialiser gridCopy
		var nextGrid [][]*gridItem // Création de la tuile suivante
		for j := 0; j < len(*grid); j++ {
			var row []*gridItem
			var copyRow []*gridItem
			for i := 0; i < len((*grid)[0]); i++ {
				cell := (*grid)[j][i]
				if cell.collapsed {
					row = append(row, cell) // Si la cellule est "collapsed", on garde la même cellule
				} else {
					copyRow = append(copyRow, cell)

					var cell_option = []int{BLANK, T_UP, T_RIGHT, T_DOWN, T_LEFT, C_UP, C_RIGHT, C_DOWN, C_LEFT} // Ajouter avec une  boucle (avec len(rules)) ?

					// Look north
					if j > 0 {
						var north = (*grid)[j-1][i]
						var validOption []int
						for _, option := range north.options {
							var valid []int = rules[option][2]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look east
					if i < len((*grid)[0])-1 {
						var east = (*grid)[j][i+1]
						var validOption []int
						for _, option := range east.options {
							var valid []int = rules[option][3]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look south
					if j < len(*grid)-1 {
						var south = (*grid)[j+1][i]
						var validOption []int
						for _, option := range south.options {
							var valid []int = rules[option][0]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look west
					if i > 0 {
						var west = (*grid)[j][i-1]
						var validOption []int
						for _, option := range west.options {
							var valid []int = rules[option][1]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					cell.options = cell_option // On modifie les options de la case
					row = append(row, cell)    // Ajouter la cellule à la ligne
				}
			}
			nextGrid = append(nextGrid, row) // Ajouter la ligne à nextGrid
			gridCopy = append(gridCopy, copyRow)
		}

		*grid = nextGrid // Cette affectation oblige de passer grid en tant que pointeur car le passage par référence par défaut d'un slice ne permet de créer de nouveau élément

		// // Calcul du pourcentage
		// var percent int = (compteur * 101) / step
		// compteur++

		// // Affichage de la barre de chargement
		// if old_percent != percent {
		// 	fmt.Printf("\r[") // permet de revenir au début de la ligne sans en ajouter une nouvelle, afin de mettre à jour la même ligne de la console.
		// 	for j := 0; j < 50; j++ {
		// 		if j < percent/2 {
		// 			fmt.Print("=")
		// 		} else {
		// 			fmt.Print(" ")
		// 		}
		// 	}
		// 	fmt.Printf("] %d%%", percent)
		// 	old_percent = percent
		// }
	}
}

func main() {
	// ----- Initialisation -----
	// Création de la grille
	var grid [][]*gridItem
	// Initialisation des éléments de la grille
	grid_init(&grid)
	// ----- Fin de l'initialisation -----

	// ----- Boucle principale -----
	subgrid := grid[:1][0:1]
	WFC(&subgrid)
	// ----- Fin de la boucle principale -----

	// TEST
	// grid[0][3].collapsed = false

	// Affichage de la grille à retourner par le serveur TCP
	fmt.Println("\n\nGrille renvoyée par le serveur TCP:")
	var grid_TCP [][]int
	for j := 0; j < DIM_Y; j++ {
		var row []int                // Créer un slice vide pour chaque ligne
		for i := 0; i < DIM_X; i++ { // Initialiser chaque cellule dans la ligne
			cell := grid[j][i]
			if cell.collapsed {
				row = append(row, cell.options[0])
			} else {
				row = append(row, -1)
			}
		}
		grid_TCP = append(grid_TCP, row) // Ajouter la ligne à la grille
	}

	fmt.Print(grid_TCP)
	client(grid)
}
