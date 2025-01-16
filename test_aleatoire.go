package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Poids associés aux options (ex : option 1 a plus de chances de sortir)
	weights := []int{10, 5, 15, 20, 10, 5, 25, 5, 5} // Somme totale : 100

	// Calculer la somme totale des poids
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// Générer un nombre aléatoire entre 0 et totalWeight
	randomValue := rand.Intn(totalWeight)

	// Choisir l'option en fonction du poids
	cumulativeWeight := 0
	for i, weight := range weights {
		cumulativeWeight += weight
		if randomValue < cumulativeWeight {
			fmt.Printf("Option choisie : %d\n", i+1)
			break
		}
	}
}
