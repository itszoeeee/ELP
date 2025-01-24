package main

import (
	"fmt"
	"sync"
)

type gridItem struct {
	val int
}

// Fonction qui ajoute 2 à chaque élément d'une sous-matrice
func addToSubMatrix(matrix []*gridItem, rowStart, rowEnd, colStart, colEnd, n int, wg *sync.WaitGroup) {
	defer wg.Done() // Indiquer que la tâche est terminée (cette ligne sera exécutée juste avant que la fonction retourne et donc garantie de tuer le process)

	for i := rowStart; i < rowEnd; i++ {
		for j := colStart; j < colEnd; j++ {
			index := i*n + j
			matrix[index].val += 2
		}
	}
}

// Worker qui récupère les sous-matrices à traiter depuis un canal
func worker(tasks <-chan func(), wg *sync.WaitGroup) {
	for task := range tasks {
		task() // Exécute la tâche
	}
}

// Fonction principale qui crée le pool de workers
func addToMatrix(matrix []*gridItem, n int) {
	var wg sync.WaitGroup
	tasks := make(chan func(), 4) // Canal pour envoyer des tâches aux workers

	// Créer un pool de 4 workers (on peut ajuster ce nombre en fonction du nombre de cœurs)
	for i := 0; i < 4; i++ {
		go worker(tasks, &wg)
	}

	// Diviser la matrice en 4 sous-matrices et envoyer les tâches aux workers
	wg.Add(4) // Nous avons 4 tâches à faire

	// Haut-gauche
	go func() {
		addToSubMatrix(matrix, 0, n/2, 0, n/2, n, &wg)
		tasks <- func() {} // Envoie une tâche pour être traitée par un worker
	}()

	// Haut-droit
	go func() {
		addToSubMatrix(matrix, 0, n/2, n/2, n, n, &wg)
		tasks <- func() {} // Envoie une tâche pour être traitée par un worker
	}()

	// Bas-gauche
	go func() {
		addToSubMatrix(matrix, n/2, n, 0, n/2, n, &wg)
		tasks <- func() {} // Envoie une tâche pour être traitée par un worker
	}()

	// Bas-droit
	go func() {
		addToSubMatrix(matrix, n/2, n, n/2, n, n, &wg)
		tasks <- func() {} // Envoie une tâche pour être traitée par un worker
	}()

	// Attendre que toutes les tâches soient terminées
	wg.Wait()

	// Fermer le canal des tâches
	close(tasks)
}

func main() {
	// Dimensions de la matrice (ex: 4x4)
	n := 4
	matrix := make([]*gridItem, n*n)

	// Initialisation de la matrice avec des éléments
	for i := 0; i < n*n; i++ {
		matrix[i] = &gridItem{val: 1}
	}

	// Affichage de la matrice avant ajout
	fmt.Println("Avant ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			index := i*n + j
			fmt.Print(matrix[index].val, " ")
		}
		fmt.Println()
	}

	// Appel de la fonction pour ajouter 2 à chaque sous-matrice en parallèle avec un pool de workers
	addToMatrix(matrix, n)

	// Affichage de la matrice après ajout
	fmt.Println("\nAprès ajout:")
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			index := i*n + j
			fmt.Print(matrix[index].val, " ")
		}
		fmt.Println()
	}
}
