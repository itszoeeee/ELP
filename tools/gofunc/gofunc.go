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
	defer wg.Done()

	// Index de la matrice à une position (i, j)
	for i := rowStart; i < rowEnd; i++ {
		for j := colStart; j < colEnd; j++ {
			// Calcul de l'indice dans la slice (matrice)
			index := i*n + j
			matrix[index].val += 2
		}
	}
}

// Fonction principale pour gérer l'ajout parallèle sur la matrice
func addToMatrix(matrix []*gridItem, n int) {
	var wg sync.WaitGroup

	// Diviser la matrice en 4 sous-matrices
	// Haut-gauche
	wg.Add(1)
	go addToSubMatrix(matrix, 0, n/2, 0, n/2, n, &wg)

	// Haut-droit
	wg.Add(1)
	go addToSubMatrix(matrix, 0, n/2, n/2, n, n, &wg)

	// Bas-gauche
	wg.Add(1)
	go addToSubMatrix(matrix, n/2, n, 0, n/2, n, &wg)

	// Bas-droit
	wg.Add(1)
	go addToSubMatrix(matrix, n/2, n, n/2, n, n, &wg)

	// Attendre que toutes les goroutines aient terminé
	wg.Wait()
}

func main() {
	// Dimensions de la matrice (ex: 4x4)
	n := 4
	var matrix []*gridItem

	// Initialisation de la matrice avec des éléments
	for i := 0; i < n*n; i++ {
		matrix = append(matrix, &gridItem{val: 1})
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

	// Appel de la fonction pour ajouter 2 à chaque sous-matrice en parallèle
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
