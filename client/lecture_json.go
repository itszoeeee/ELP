package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Json struct {
	Fichiert     string `json:"imaget"`
	Orientationt string `json:"orientationt"`
	Fichierc     string `json:"imagec"`
	Orientationc string `json:"orientationc"`
}

var corresp = map[string]int{
	"up":    1,
	"right": 2,
	"down":  3,
	"left":  4,
}

func ouverture_json(input_file string) ([]int, []string) {
	// Ouvrir le fichier JSON
	file, err := os.Open(input_file)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return nil, nil
	}
	defer file.Close()

	// Lire le contenu du fichier
	bxteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return nil, nil
	}

	// Décoder le JSON dans une structure
	var images Json
	err = json.Unmarshal(bxteValue, &images) //unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return nil, nil
	}

	if images.Orientationt == "" || images.Orientationc == "" {
		fmt.Println("Orientation invalide :", images.Orientationc, images.Orientationt)
		return nil, nil
	}

	return []int{corresp[images.Orientationt], corresp[images.Orientationc]}, []string{images.Fichiert, images.Fichierc}
}
