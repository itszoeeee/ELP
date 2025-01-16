package main

import "fmt"

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

func main() {
	// Option initiale
	option := []int{0, 1, 2, 3, 4}

	// Déclaration des tableaux
	array1 := []int{0, 1}
	array2 := []int{1, 2}

	// Concatenation des deux tableaux
	valid := append(array1, array2...)
	valid2 := []int{2, 4, 1, 0, 4}

	fmt.Println("valid:", valid)
	fmt.Println("valid2:", valid2)

	// Appel de la fonction checkValid
	checkValid(&option, valid2)
	checkValid(&option, valid)

	fmt.Println("option après vérification:", option)
}
