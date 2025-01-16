package main

import (
	"fmt"
	"sync"
)

type gridItem struct {
	value int
}

// Fonction qui ajoute 2 à chaque élément d'une sous-matrice
func addToSubMatrix(matrix []*gridItem, rowStart, rowEnd, colStart, colEnd, n int, wg *sync.WaitGroup) {
	defer wg.Done() // Indiquer que la tâche est terminée

	for i := rowStart; i < rowEnd; i++ {
		for j := colStart; j < colEnd; j++ {
			index := i*n + j
			matrix[index].value += 2
		}
	}
}

// Worker qui récupère les sous-matrices à traiter depuis un canal
func worker(tasks <-chan func(), wg *sync.WaitGroup) {
	for task := range tasks {
		task() // Exécute la tâche
	}
}

// Fonction principale qui crée le pool de workers et divise la matrice
func addToMatrix(matrix []*gridItem, n, x, y, numWorkers int) {
	var wg sync.WaitGroup
	tasks := make(chan func(), x*y) // Canal pour envoyer des tâches aux workers

	// Créer un pool de `numWorkers` workers (4 workers dans votre cas)
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

			// Envoyer la tâche au canal, qui sera récupéré par un worker
			tasks <- func() {
				addToSubMatrix(matrix, rowStart, rowEnd, colStart, colEnd, n, &wg)
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
	numWorkers := 4 // Nombre de workers (goroutines) à utiliser en parallèle

	matrix := make([]*gridItem, n*n)

	// Initialisation de la matrice avec des éléments
	for i := 0; i < n*n; i++ {
		matrix[i] = &gridItem{value: 1}
	}

	// Affichage de la matrice avant ajout
	fmt.Println("Avant ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			index := i*n + j
			fmt.Print(matrix[index].value, " ")
		}
		fmt.Println()
	}

	// Appel de la fonction pour ajouter 2 à chaque sous-matrice en parallèle avec un pool de workers
	addToMatrix(matrix, n, x, y, numWorkers)

	// Affichage de la matrice après ajout
	fmt.Println("\nAprès ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			index := i*n + j
			fmt.Print(matrix[index].value, " ")
		}
		fmt.Println()
	}
}
