package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const DIM = 20

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
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, CROSS},              // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},                  // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS},  // west
	},
	// T_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS},      // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS},    // south
		{BLANK, T_LEFT, C_UP, C_LEFT},                    // west
	},
	// T_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT},                    // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS},     // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS}, // west
	},
	// T_LEFT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},                // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS},    // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS},  // west
	},
	// C_UP
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},                // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},                  // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS},  // west
	},
	// C_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS},      // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},                  // south
		{BLANK, T_LEFT, C_UP, C_LEFT},                    // west
	},
	// C_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT},                  // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS},   // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS}, // south
		{BLANK, T_LEFT, C_UP, C_LEFT},                 // west
	},
	// C_LEFT
	{
		{BLANK, T_UP, C_UP, C_RIGHT},                    // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},               // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS}, // west
	},
	//CROSS
	{
		{T_DOWN, T_RIGHT, T_LEFT, C_LEFT, C_DOWN}, //north
		{T_LEFT, T_UP, T_DOWN, C_LEFT, C_UP},      //east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT},    //south
		{T_RIGHT, T_UP, T_DOWN, C_RIGHT, C_DOWN},  //west
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
	for i := 0; i < len(rules); i++ {
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
	Tiles[9], err = loadImage("pattern/cross.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'crosss.png': %w", err)
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

func client(grid []*gridItem) {
	// Création des tuiles
	Tiles, err := createTile()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Dimensions de l'image de sortie
	gridWidth, gridHeight := DIM, DIM
	cellSize := Tiles[0].Bounds().Dx() // On suppose que la cellule fait la même taille que les images chargées

	// Créer de l'image de sortie
	outputImage := createEmptyImage(cellSize*gridWidth, cellSize*gridHeight)

	// Affichage des tuiles (si la tuile est collapsed)
	for i := 0; i < DIM; i++ {
		for j := 0; j < DIM; j++ {
			var cell = grid[i+j*DIM]
			if cell.collapsed && len(cell.options) != 0 {
				var index = cell.options[0]
				placeImageInMatrix(outputImage, Tiles[index], i, j, cellSize)
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

func init_grid() []*gridItem {
	// Création de la grille
	var grid []*gridItem

	// Initialisation des éléments de la grille
	for i := 0; i < DIM*DIM; i++ {
		// Crée une nouvelle instance de gridItem
		cell := &gridItem{
			collapsed: false, // Initialisé à false
			// Ajouter avec une  boucle (avec len(rules)) ?
			options: []int{BLANK, T_UP, T_RIGHT, T_DOWN, T_LEFT, C_UP, C_RIGHT, C_DOWN, C_LEFT, CROSS}, // Options fixes
		}
		grid = append(grid, cell) // Ajouter l'élément à la grille
	}
	return grid
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

// Structure représentant un élément avec une valeur et un poids.
type WeightedItem struct {
	Value  int
	Weight int
}

// Fonction pour effectuer un tirage pondéré.
func weightedRandom(items []WeightedItem) int {
	// Calculer la somme totale des poids.
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}

	// Générer un nombre aléatoire entre 0 et la somme des poids.
	rand.Seed(time.Now().UnixNano())
	randomWeight := rand.Intn(totalWeight)

	// Parcourir les éléments pour trouver celui correspondant au poids généré.
	currentWeight := 0
	for _, item := range items {
		currentWeight += item.Weight
		if randomWeight < currentWeight {
			return item.Value
		}
	}

	return 0 // Ne devrait jamais arriver si les poids sont bien définis.
}

func proba(liste []int, p int) []WeightedItem { // fonction qui prend une liste d'options et leur associe un poids
	var items []WeightedItem
	containsBlank := false //pour verifier si blank ets bien disponible parmi les options possibles
	for _, elem := range liste {
		if elem == 0 {
			containsBlank = true
			break
		}
	}

	// Parcours de la liste
	for i := 0; i < len(liste); i++ {
		if containsBlank {
			if liste[i] == 0 {
				items = append(items, WeightedItem{0, p}) // l'option 0 correspond à "blank"
			} else {
				items = append(items, WeightedItem{liste[i], (100 - p) / (len(liste) - 1)})
			}
		} else {
			items = append(items, WeightedItem{liste[i], 100 / len(liste)})
		}
	}
	return items
}

func main() {
	// ----- Initialisation -----
	// Création de la grille
	var grid []*gridItem

	grid = init_grid()
	// ----- Fin de l'initialisation -----

	// ----- Boucle principale -----

	var gridCopy []*gridItem
	gridCopy = grid
	var step int = len(gridCopy)
	var compteur int = 0
	var old_percent int = -1

	for len(gridCopy) > 1 { // Boucle principale
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
			// var pick int = randomItem.options[rand.Intn(len(randomItem.options))] // choisir un option disponible (aléatoirement)
			// randomItem.options = []int{pick}
			var p = 50 // probabilité qu'on veut de tirer "blank" (en %)
			var items []WeightedItem = proba(randomItem.options, p)
			var result = weightedRandom(items)
			randomItem.options = []int{result}
		}

		// Création de la tuile suivante
		var nextGrid []*gridItem
		for i := 0; i < DIM; i++ {
			for j := 0; j < DIM; j++ {
				var index = j + i*DIM
				if grid[index].collapsed {
					nextGrid = append(nextGrid, grid[index])
				} else {
					var cell_option = []int{BLANK, T_UP, T_RIGHT, T_DOWN, T_LEFT, C_UP, C_RIGHT, C_DOWN, C_LEFT, CROSS} // Ajouter avec une  boucle (avec len(rules)) ?

					// Look north
					if i > 0 {
						var north = grid[j+(i-1)*DIM]
						var validOption []int
						for _, option := range north.options {
							var valid []int = rules[option][2]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look east
					if j < DIM-1 {
						var east = grid[j+1+i*DIM]
						var validOption []int
						for _, option := range east.options {
							var valid []int = rules[option][3]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look south
					if i < DIM-1 {
						var south = grid[j+(i+1)*DIM]
						var validOption []int
						for _, option := range south.options {
							var valid []int = rules[option][0]
							validOption = append(validOption, valid...)
						}
						checkValid(&cell_option, validOption)
					}
					// Look west
					if j > 0 {
						var west = grid[j-1+i*DIM]
						var validOption []int
						for _, option := range west.options {
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

	// TEST
	//grid[3].collapsed = false

	// Affichage de la grille à retourner par le serveur TCP
	var grid_TCP []int
	fmt.Println("\n\n Grille renvoyée par le serveur TCP:")
	for _, v := range grid {
		if v.collapsed {
			grid_TCP = append(grid_TCP, v.options[0])
		} else {
			grid_TCP = append(grid_TCP, -1)
		}
	}

	fmt.Print(grid_TCP)
	client(grid)
}
