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
	Fichierc     string `json:"imaget"`
	Orientationc string `json:"orientationt"`
}

var corresp = map[string]int{
	"up":    1,
	"right": 2,
	"down":  3,
	"left":  4,
}

func ouverture_json(input_file string) []int {
	// Ouvrir le fichier JSON
	file, err := os.Open(input_file)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return nil
	}
	defer file.Close()

	// Lire le contenu du fichier
	bxteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return nil
	}

	// Décoder le JSON dans une structure
	var images Json
	err = json.Unmarshal(bxteValue, &images) //unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return nil
	}

	if images.Orientationt == "" || images.Orientationc == "" {
		fmt.Println("Orientation invalide :", image1.Orientation)
		return nil
	}

	return []int{corresp[images.Orientationt], corresp[images.Orientationc]}
}
