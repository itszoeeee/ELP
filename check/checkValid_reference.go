package main

import "fmt"

// Fonction pour vérifier si l'option est valide
func checkValid_ref(option []int, valid []int) {
	for i := len(option) - 1; i >= 0; i-- {
		element := option[i]
		found := false
		// Vérifie si l'élément est dans le tableau valid
		for _, v := range valid {
			if element == v {
				found = true
				break
			}
		}
		// Si l'élément n'est pas valide, on le supprime
		if !found {
			option = append(option[:i], option[i+1:]...)
		}
	}
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
	checkValid_ref(option, valid)

	fmt.Println("option après vérification:", option)
}
