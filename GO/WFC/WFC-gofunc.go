package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const DIM_X = 11
const DIM_Y = 11

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

// Définir le tableau de règles
var rules = [][][]int{
	// BLANK
	{
		{BLANK, T_UP, C_UP, C_RIGHT, F_H},      // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN, F_V}, // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT, F_H},   // south
		{BLANK, T_LEFT, C_UP, C_LEFT, F_V},     // west
	},
	// T_UP
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS, F_V}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, CROSS, F_H},              // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT, F_H},                  // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS, F_H},  // west
	},
	// T_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS, F_V}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS, F_H},      // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS, F_V},    // south
		{BLANK, T_LEFT, C_UP, C_LEFT, F_V},                    // west
	},
	// T_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT, F_H},                    // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS, F_H},     // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS, F_V},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS, F_H}, // west
	},
	// T_LEFT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS, F_V}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN, F_V},                // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS, F_V},    // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS, F_H},  // west
	},
	// C_UP
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS, F_V}, // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN, F_V},                // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT, F_H},                  // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS, F_H},  // west
	},
	// C_RIGHT
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS, F_V}, // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS, F_H},      // east
		{BLANK, T_DOWN, C_DOWN, C_LEFT, F_H},                  // south
		{BLANK, T_LEFT, C_UP, C_LEFT, F_V},                    // west
	},
	// C_DOWN
	{
		{BLANK, T_UP, C_UP, C_RIGHT, F_H},                  // north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS, F_H},   // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS, F_V}, // south
		{BLANK, T_LEFT, C_UP, C_LEFT, F_V},                 // west
	},
	// C_LEFT
	{
		{BLANK, T_UP, C_UP, C_RIGHT, F_H},                    // north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN, F_V},               // east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS, F_V},   // south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS, F_H}, // west
	},
	//CROSS
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, F_V}, //north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, F_H},      //east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, F_V},    //south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, F_H},  //west
	},
	//F_H
	{
		{BLANK, T_UP, C_UP, C_RIGHT},                    //north
		{T_UP, T_DOWN, T_LEFT, C_UP, C_LEFT, CROSS},     //east
		{BLANK, T_DOWN, C_DOWN, C_LEFT},                 //south
		{T_UP, T_RIGHT, T_DOWN, C_RIGHT, C_DOWN, CROSS}, //west
	},
	//F_V
	{
		{T_RIGHT, T_DOWN, T_LEFT, C_DOWN, C_LEFT, CROSS}, //north
		{BLANK, T_RIGHT, C_RIGHT, C_DOWN},                //east
		{T_UP, T_RIGHT, T_LEFT, C_UP, C_RIGHT, CROSS},    //south
		{BLANK, T_LEFT, C_UP, C_LEFT},                    //west
	},
}

type gridItem struct {
	collapsed bool  // Champ pour "collapsed"
	options   []int // Un tableau de tableaux d'entiers
}

var completedCells int64 = 0

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

	Tiles[0], err = loadImage("GO/pattern/blank.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'blank.png': %w", err)
	}
	Tiles[1], err = loadImage("GO/pattern/t_up.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_up.png': %w", err)
	}
	Tiles[2], err = loadImage("GO/pattern/t_right.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_right.png': %w", err)
	}
	Tiles[3], err = loadImage("GO/pattern/t_down.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_down.png': %w", err)
	}
	Tiles[4], err = loadImage("GO/pattern/t_left.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 't_left.png': %w", err)
	}
	Tiles[5], err = loadImage("GO/pattern/c_up.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_up.png': %w", err)
	}
	Tiles[6], err = loadImage("GO/pattern/c_right.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_right.png': %w", err)
	}
	Tiles[7], err = loadImage("GO/pattern/c_down.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_down.png': %w", err)
	}
	Tiles[8], err = loadImage("GO/pattern/c_left.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'c_left.png': %w", err)
	}
	Tiles[9], err = loadImage("GO/pattern/cross.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'cross.png': %w", err)
	}
	Tiles[10], err = loadImage("GO/pattern/f_h.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'f_h.png': %w", err)
	}
	Tiles[11], err = loadImage("GO/pattern/f_v.png")
	if err != nil {
		return Tiles, fmt.Errorf("Erreur lors du chargement de l'image 'f_v.png': %w", err)
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

func affichage(grid [][]*gridItem) {
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
	outFile, err := os.Create("GO/output.png")
	if err != nil {
		fmt.Println("\n\nErreur lors de la création de l'image de sortie:\n", err)
		return
	}
	defer outFile.Close()

	err = png.Encode(outFile, outputImage)
	if err != nil {
		fmt.Println("\n\nErreur lors de l'exportation de l'image:\n", err)
	}

	fmt.Println("\n\nImage exportée avec succès dans output.png\n")
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
	for j := 0; j < DIM_Y; j++ { // Initialiser chaque ligne de la grille
		var row []*gridItem          // Créer un slice vide pour chaque ligne
		for i := 0; i < DIM_X; i++ { // Initialiser chaque cellule dans la ligne
			// Options par défaut
			var default_option []int
			for k := 0; k < len(rules); k++ {
				default_option = append(default_option, k)
			}
			cell := &gridItem{ // Créer une nouvelle instance de gridItem
				collapsed: false, // Initialisé à false
				options:   default_option,
			}
			row = append(row, cell) // Ajouter la cellule à la ligne
		}
		*grid = append(*grid, row) // Ajouter la ligne à la grille
	}
}

func WFC(grid *[][]*gridItem, step int) {
	gridCopy := *grid

	for _, row := range *grid {
		for _, cell := range row {
			if !cell.collapsed {
				step++
			}
		}
	}

	for k := 0; k < step; k++ { // Boucle principale
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
		// rand.Seed(time.Now().UnixNano()) // Initialiser le générateur de nombres aléatoires avec l'heure actuelle
		var randomItem *gridItem = smallestItems[rand.Intn(len(smallestItems))] // Sélectionner une clé aléatoire parmi celles disponibles
		randomItem.collapsed = true                                             // collapsed l'élément
		if len(randomItem.options) != 0 {                                       // vérifie qu'il existe une option disponible, sinon affiche un carré vide
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
					var cell_option []int
					for k := 0; k < len(rules)-1; k++ {
						cell_option = append(cell_option, k)
					}
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
		*grid = nextGrid                    // Cette affectation oblige de passer grid en tant que pointeur car le passage par référence par défaut d'un slice ne permet de créer de nouveau élément
		atomic.AddInt64(&completedCells, 1) // Incrémente le nombre d'itération pour le calcul du pourcentage
	}
}

// Worker qui récupère les sous-grilles à traiter depuis un canal
func worker(tasks <-chan func(), wg *sync.WaitGroup) {
	for task := range tasks {
		task() // Exécute la tâche
		wg.Done()
	}
}

func multi_process(grid *[][]*gridItem, div_x, div_y, numWorkers int) {
	var wg sync.WaitGroup
	wg.Add(div_x * div_y)                   // Diviser la matrice en div_x * div_y sous-grilles et envoyer les tâches aux workers
	tasks := make(chan func(), div_x*div_y) // Canal pour envoyer des tâches aux workers

	// Créer un pool de workers
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, &wg) // création de go function
	}

	// Calcul des tailles des sous-grilles
	colsPerSubGrid := DIM_X / div_x
	rowsPerSubGrid := DIM_Y / div_y

	for j := 0; j < div_y; j++ {
		for i := 0; i < div_x; i++ {
			// Calcul des indices pour chaque sous-matrice
			rowStart := (j * rowsPerSubGrid) + 2 // permet de créer la séparation entre les sous-grilles à compléter
			rowEnd := ((j + 1) * rowsPerSubGrid) - 2
			colStart := (i * colsPerSubGrid) + 2
			colEnd := ((i + 1) * colsPerSubGrid) - 2

			// Ajuster la première et la dernière sous-matrice pour qu'elle couvre tout l'espace (en cas de division non parfaitement égale) elle permet également de calculer les sous grille jusqu'au bordure de la grande
			if j == 0 {
				rowStart = 0
			}
			if i == 0 {
				colStart = 0
			}
			if j == div_y-1 {
				rowEnd = DIM_Y
			}
			if i == div_x-1 {
				colEnd = DIM_X
			}

			// Créer un slice pour la sous-matrice spécifique à cette tâche
			if rowEnd > rowStart && colEnd > colStart { // On vérifie qu'on ne crée pas un slice vide
				subGrid := make([][]*gridItem, rowEnd-rowStart)
				for r := rowStart; r < rowEnd; r++ {
					subGrid[r-rowStart] = (*grid)[r][colStart:colEnd]
				}

				// Envoyer la tâche au canal, qui sera récupéré par un worker
				tasks <- func() {
					WFC(&subGrid, 0) // Remarque : c'est inutile de passer subGrid par référence car Go le fait indirectement mais l'objectif est de se rapprocher de la fonction WFC qui elle necessite un passage par référence à cause de nextGrid
				}
			}
		}
	}
	wg.Wait()    // Attendre que toutes les tâches soient terminées
	close(tasks) // Fermer le canal des tâches une fois qu'elles sont toutes envoyées
}

// Fonction de rapporteur de progression
func progressReporter(stopChan <-chan struct{}) {
	var lastPercentage = -1

	for {
		select {
		case <-time.After(time.Millisecond):
			percentage := int((atomic.LoadInt64(&completedCells) * 101) / (DIM_X * DIM_Y)) // converstion en int
			if percentage != lastPercentage {
				fmt.Printf("\r[") // permet de revenir au début de la ligne sans en ajouter une nouvelle, afin de mettre à jour la même ligne de la console.
				for j := 0; j < 50; j++ {
					if j < percentage/2 {
						fmt.Print("=")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Printf("] %d%%", percentage)
				lastPercentage = percentage
			}
		case <-stopChan:
			fmt.Printf("\nGénération de l'image terminée")
			return
		}
	}
}

func app(numWorkers, div_x, div_y int) {
	if div_x <= 0 || div_y <= 0 {
		fmt.Println("\n\nErreur : div_x et div_y doivent être strictement positifs\n")
	}
	if numWorkers == 0 {
		fmt.Println("\n\nErreur : numWorkers doit être strictement positif\n")
	} else {
		stopChan := make(chan struct{}) // Créer un canal pour arrêter le rapporteur de progression
		go progressReporter(stopChan)   // Lancer le calcul de progression

		// ----- Initialisation -----
		var grid [][]*gridItem // Création de la grille
		grid_init(&grid)       // Initialisation des éléments de la grille
		// ----- Fin de l'initialisation -----

		// ----- Boucle principale -----
		if div_x == 1 && div_y == 1 {
			WFC(&grid, 0)
		} else {
			multi_process(&grid, div_x, div_y, numWorkers)
			WFC(&grid, 1)
		}
		// ----- Fin de la boucle principale -----

		// Grille à retourner par le serveur TCP
		var grid_TCP [][]int
		for j := 0; j < DIM_Y; j++ {
			var row []int                // Créer un slice vide pour chaque ligne
			for i := 0; i < DIM_X; i++ { // Initialiser chaque cellule dans la ligne
				cell := grid[j][i]
				if cell.collapsed {
					row = append(row, cell.options...)
				} else {
					row = append(row, -1)
				}
			}
			grid_TCP = append(grid_TCP, row) // Ajouter la ligne à la grille
		}

		close(stopChan)                    // Ferme le canal de l'avancement
		time.Sleep(100 * time.Millisecond) // Attendre que le canal de l'avancement soit bien fermé
		// fmt.Print(grid_TCP)
		affichage(grid)
	}
}

func main() {
	nbWorkers := 9
	div_x := 3
	div_y := 3

	app(nbWorkers, div_x, div_y)
}
