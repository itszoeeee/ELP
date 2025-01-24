package main

import (
	"fmt"
	"time"
)

func main() {
	// Variable x définissant la taille de la liste
	x := 50

	// Affiche une barre de chargement pour parcourir la liste entre 0 et x
	for i := 0; i <= x; i++ {
		// Calcul du pourcentage d'avancement
		percent := (i * 100) / x

		// Affichage de la barre de chargement
		fmt.Printf("\r[") // permet de revenir au début de la ligne sans en ajouter une nouvelle, afin de mettre à jour la même ligne de la console.
		for j := 0; j < 50; j++ {
			if j < percent/2 {
				fmt.Print("=")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Printf("] %d%%", percent)

		// Pause de 1 seconde
		time.Sleep(1 * time.Second)
	}

	// Une fois la boucle terminée, affiche "terminé"
	fmt.Printf("\nChargement terminé !\n")
}
