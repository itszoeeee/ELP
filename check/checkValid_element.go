package main

import "fmt"

// Fonction pour vérifier si l'option est valide
func checkValid_elem(option []int, valid []int) []int {
	var result []int
	for _, element := range option {
		found := false
		// Vérifie si l'élément est dans le tableau valid
		for _, v := range valid {
			if element == v {
				found = true
				break
			}
		}
		// Si l'élément est valide, on l'ajoute à la tranche result
		if found {
			result = append(result, element)
		}
	}
	return result
}

func main() {
	// Déclaration des tableaux
	array1 := []int{0, 1}
	array2 := []int{1, 2}

	// Concatenation des deux tableaux
	valid := append(array1, array2...)

	fmt.Println("valid:", valid)

	// Option initiale
	option := []int{0, 1, 2, 3, 4}

	// Appel de la fonction checkValid
	option = checkValid_elem(option, valid)

	fmt.Println("option après vérification:", option)
}
