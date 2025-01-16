package main

import (
	"fmt"
	"sync"
)

type gridItem struct {
	val int
}

// Fonction qui ajoute 2 à chaque élément d'une sous-matrice
func addToSubMatrix(subMatrix [][]*gridItem, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < len(subMatrix); i++ {
		for j := 0; j < len(subMatrix[i]); j++ {
			subMatrix[i][j].val += 2
		}
	}

	// Affichage de la sous-matrice pendant ajout
	fmt.Println("\nPendant ajout:")
	for i := 0; i < len(subMatrix); i++ {
		for j := 0; j < len(subMatrix[i]); j++ {
			fmt.Print(subMatrix[i][j].val, " ")
		}
		fmt.Println()
	}
}

// Worker qui récupère les sous-matrices à traiter depuis un canal
func worker(tasks <-chan func(), wg *sync.WaitGroup) {
	for task := range tasks {
		task() // Exécute la tâche
	}
}

// Fonction principale qui crée le pool de workers et divise la matrice
func addToMatrix(matrix [][]*gridItem, n, x, y, numWorkers int) {
	var wg sync.WaitGroup
	tasks := make(chan func(), x*y) // Canal pour envoyer des tâches aux workers

	// Créer un pool de `numWorkers` workers
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, &wg)
	}

	// Calcul des tailles des sous-matrices
	rowsPerSubMatrix := n / x
	colsPerSubMatrix := n / y

	// Diviser la matrice en x * y sous-matrices et envoyer les tâches aux workers
	wg.Add(x * y) // Nous avons x * y tâches à faire

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			// Calcul des indices pour chaque sous-matrice
			rowStart := i * rowsPerSubMatrix
			rowEnd := (i + 1) * rowsPerSubMatrix
			colStart := j * colsPerSubMatrix
			colEnd := (j + 1) * colsPerSubMatrix

			// Ajuster la dernière sous-matrice pour qu'elle couvre tout l'espace (en cas de division non parfaitement égale)
			if i == x-1 {
				rowEnd = n
			}
			if j == y-1 {
				colEnd = n
			}

			// Créer un slice pour la sous-matrice spécifique à cette tâche
			subMatrix := make([][]*gridItem, rowEnd-rowStart)
			for r := rowStart; r < rowEnd; r++ {
				subMatrix[r-rowStart] = matrix[r][colStart:colEnd]
			}

			// Envoyer la tâche au canal, qui sera récupéré par un worker
			tasks <- func() {
				addToSubMatrix(subMatrix, &wg)
			}
		}
	}

	// Attendre que toutes les tâches soient terminées
	wg.Wait()

	// Fermer le canal des tâches une fois qu'elles sont toutes envoyées
	close(tasks)
}

func main() {
	// Dimensions de la matrice (par exemple 6x6)
	n := 6          // Taille de la matrice
	x := 2          // Nombre de divisions sur les lignes
	y := 3          // Nombre de divisions sur les colonnes
	numWorkers := 1 // Nombre de workers (goroutines) à utiliser en parallèle

	var matrix [][]*gridItem

	// Initialisation de la matrice avec des éléments
	for i := 0; i < n; i++ {
		var row []*gridItem
		for j := 0; j < n; j++ {
			row = append(row, &gridItem{val: 1})
		}
		matrix = append(matrix, row)
	}

	// Affichage de la matrice avant ajout
	fmt.Println("Avant ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Print(matrix[i][j].val, " ")
		}
		fmt.Println()
	}

	// Appel de la fonction pour ajouter 2 à chaque sous-matrice en parallèle avec un pool de workers
	addToMatrix(matrix, n, x, y, numWorkers)

	// Affichage de la matrice après ajout
	fmt.Println("\nAprès ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			fmt.Print(matrix[i][j].val, " ")
		}
		fmt.Println()
	}
}
