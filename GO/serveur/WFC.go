package main

import (
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
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

// Structure représentant les cellules de la grille
type gridItem struct {
	collapsed bool  // Champ pour "collapsed"
	options   []int // Un tableau de tableaux d'entiers
}

// Structure représentant une option avec une valeur et un poids
type weightedItem struct {
	value  int
	weight int
}

// Fonction pour effectuer un tirage pondéré
func weighted_random(weightedOptions []weightedItem) int {
	// Génère un nombre aléatoire entre 0 et 100 (poids total)
	randomWeight := rand.Intn(100)

	// Parcours les éléments pour trouver celui correspondant au poids généré
	currentWeight := 0
	for _, option := range weightedOptions {
		currentWeight += option.weight
		if randomWeight <= currentWeight {
			return option.value
		}
	}
	return BLANK // Ne devrait jamais arriver si les poids sont bien définis
}

// Fonction qui prend une liste d'options et leur associe un poids
func set_weight(options []int, proba int) []weightedItem {
	var weightedOptions []weightedItem
	containsBlank := false // Vérifie si blank est bien disponible parmi les options possibles
	for _, option := range options {
		if option == BLANK {
			containsBlank = true
			break
		}
	}

	// Parcours des options
	for i := 0; i < len(options); i++ {
		if containsBlank {
			if options[i] == BLANK {
				weightedOptions = append(weightedOptions, weightedItem{BLANK, proba}) // Ajoute un poids de proba
			} else {
				weightedOptions = append(weightedOptions, weightedItem{options[i], (100 - proba) / (len(options) - 1)}) // Ajoute un poids de 1-proba (ne pas tirer un BLANK) selon toutes les autres options disponibles sauf BLANK
			}
		} else {
			weightedOptions = append(weightedOptions, weightedItem{options[i], 100 / len(options)})
		}
	}
	return weightedOptions
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

// Fonction pour initialiser la grille
func grid_init(grid *[][]*gridItem, dim_x, dim_y int) {
	for j := 0; j < dim_y; j++ { // Initialiser chaque ligne de la grille
		var row []*gridItem          // Créer un slice vide pour chaque ligne
		for i := 0; i < dim_x; i++ { // Initialiser chaque cellule dans la ligne
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

// Algorithme de Wave Function Collapsed
func WFC(grid *[][]*gridItem, step, proba int, nb_cell_collapsed *int64) {
	gridCopy := *grid

	for _, row := range *grid { // Vérifie les cellules déjà collapsed
		for _, cell := range row {
			if !cell.collapsed {
				step++
			}
		}
	}

	for k := 0; k < step; k++ { // Boucle principale

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
		var randomItem *gridItem = smallestItems[rand.Intn(len(smallestItems))] // Sélectionner une clé aléatoire parmi celles disponibles
		randomItem.collapsed = true                                             // Collapsed l'élément
		if len(randomItem.options) != 0 {                                       // Vérifie qu'il existe une option disponible, sinon affiche un carré vide
			var weightedOptions []weightedItem = set_weight(randomItem.options, proba)
			var pick = weighted_random(weightedOptions) // Choisir un option disponible (aléatoirement pondérées)
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
		*grid = nextGrid                      // Cette affectation oblige de passer grid en tant que pointeur car le passage par référence par défaut d'un slice ne permet de créer de nouveau élément
		atomic.AddInt64(nb_cell_collapsed, 1) // Incrémente le nombre d'itération pour le calcul du pourcentage
	}
}

// Worker qui récupère les sous-grilles à traiter depuis un canal
func worker(tasks <-chan func(), wg *sync.WaitGroup) {
	for task := range tasks {
		task()    // Exécute la tâche
		wg.Done() // Ferme le process
	}
}

// Fonction pour envoyer le calcul de WFC dans une Goroutines
func multi_process(grid *[][]*gridItem, dim_x, dim_y, proba, div_x, div_y, numWorkers int, nb_cell_collapsed *int64) {
	var wg sync.WaitGroup
	wg.Add(div_x * div_y)                   // Diviser la matrice en div_x * div_y sous-grilles et envoyer les tâches aux workers
	tasks := make(chan func(), div_x*div_y) // Canal pour envoyer des tâches aux workers

	// Créer un pool de workers
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, &wg) // création de go function
	}

	// Calcul des tailles des sous-grilles
	colsPerSubGrid := dim_x / div_x
	rowsPerSubGrid := dim_y / div_y

	for j := 0; j < div_y; j++ {
		for i := 0; i < div_x; i++ {
			// Calcul des indices pour chaque sous-matrice
			rowStart := (j * rowsPerSubGrid) + 3 // permet de créer la séparation entre les sous-grilles à compléter
			rowEnd := ((j + 1) * rowsPerSubGrid) - 3
			colStart := (i * colsPerSubGrid) + 3
			colEnd := ((i + 1) * colsPerSubGrid) - 3

			// Ajuster la première et la dernière sous-matrice pour qu'elle couvre tout l'espace (en cas de division non parfaitement égale) elle permet également de calculer les sous grille jusqu'au bordure de la grande
			if j == 0 {
				rowStart = 0
			}
			if i == 0 {
				colStart = 0
			}
			if j == div_y-1 {
				rowEnd = dim_y
			}
			if i == div_x-1 {
				colEnd = dim_x
			}

			// Créer un slice pour la sous-matrice spécifique à cette tâche
			if rowEnd > rowStart && colEnd > colStart { // On vérifie qu'on ne crée pas un slice vide
				subGrid := make([][]*gridItem, rowEnd-rowStart)
				for r := rowStart; r < rowEnd; r++ {
					subGrid[r-rowStart] = (*grid)[r][colStart:colEnd]
				}

				// Envoyer la tâche au canal, qui sera récupéré par un worker
				tasks <- func() {
					WFC(&subGrid, 0, proba, nb_cell_collapsed) // Remarque : c'est inutile de passer subGrid par référence car Go le fait indirectement mais l'objectif est de se rapprocher de la fonction WFC qui elle necessite un passage par référence à cause de nextGrid
				}
			}
		}
	}
	wg.Wait()    // Attendre que toutes les tâches soient terminées
	close(tasks) // Fermer le canal des tâches une fois qu'elles sont toutes envoyées
}

func progress(stopChan <-chan struct{}, conn net.Conn, dim_x, dim_y int, nb_cell_collapsed *int64) { // Calcul de la progession
	var lastPercentage = -1

	for {
		select {
		case <-time.After(100 * time.Millisecond): // Toutes les 100 millisecondes
			percentage := int((atomic.LoadInt64(nb_cell_collapsed) * 100)) / (dim_x * dim_y) // Converstion en int de la structure atomic pour protéger l'accès
			if percentage != lastPercentage {
				send_int(conn, percentage) // Envoie du pourcentage au client
				lastPercentage = percentage
			}
		case <-stopChan:
			send_int(conn, 100) // Signal que la grille est complète
			return
		}
	}
}

func grid_process(grid_TCP *[][]int, dim_x, dim_y, proba, div_x, div_y, numWorkers int, nb_cell_collapsed *int64) {

	// ----- Initialisation -----
	var grid [][]*gridItem         // Création de la grille
	grid_init(&grid, dim_x, dim_y) // Initialisation des éléments de la grille

	// ----- Boucle principale -----
	if div_x == 1 && div_y == 1 { // Si on utilise pas de Goroutines
		WFC(&grid, 0, proba, nb_cell_collapsed)
	} else {
		multi_process(&grid, dim_x, dim_y, proba, div_x, div_y, numWorkers, nb_cell_collapsed)
		WFC(&grid, 1, proba, nb_cell_collapsed)
	}

	// Grille à retourner par le serveur TCP
	for j := 0; j < dim_y; j++ {
		var row []int                // Créer un slice vide pour chaque ligne
		for i := 0; i < dim_x; i++ { // Initialiser chaque cellule dans la ligne
			cell := grid[j][i]
			if cell.collapsed && len(cell.options) != 0 { // Vérifie si la cellule est bien collapsed
				row = append(row, cell.options[0])
			} else {
				row = append(row, -1) // Renvoie -1 sur la cellule n'est pas collapsed
			}
		}
		*grid_TCP = append(*grid_TCP, row) // Ajouter la ligne à la grille
	}
}
