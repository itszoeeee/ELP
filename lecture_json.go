package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Image struct {
	Fichier     string `json:"image"`
	Orientation string `json:"orientation"`
}

var corresp = map[string]int{
	"up":    1,
	"right": 2,
	"down":  3,
	"left":  4,
}

func ouverture_json(input_file string) int {
	// Ouvrir le fichier JSON
	file, err := os.Open(input_file)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return (-1)
	}
	defer file.Close()

	// Lire le contenu du fichier
	bxteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return (-1)
	}

	// Décoder le JSON dans une structure
	var image1 Image
	err = json.Unmarshal(bxteValue, &image1) // unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return (-1)
	}

	if image1.Orientation == "" {
		fmt.Println("Orientation invalide :", image1.Orientation)
		return (-1)
	}

	return (corresp[image1.Orientation])
}
