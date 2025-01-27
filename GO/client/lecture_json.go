package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type jsonItem struct {
	TCP_address  string `json:"TCP_address"`
	TCP_port     string `json:"TCP_port"`
	Fichiert     string `json:"image_t"`
	Orientationt string `json:"orientation_t"`
	Fichierc     string `json:"image_c"`
	Orientationc string `json:"orientation_c"`
	Fichierf     string `json:"forward"`
	Orientationf string `json:"orientation_forward"`
	Fichierblank string `json:"blank"`
	Fichiercross string `json:"cross"`
}

var corresp = map[string]int{
	"up":    1,
	"right": 2,
	"down":  3,
	"left":  4,
}

func lecture_json(input_file string) ([]int, []string) {
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
	var json_data jsonItem
	err = json.Unmarshal(bxteValue, &json_data) // Unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return nil, nil
	}

	if json_data.Orientationt == "" || json_data.Orientationc == "" || json_data.Orientationf == "" {
		fmt.Println("Orientation invalide :", json_data.Orientationc, json_data.Orientationt, json_data.Orientationf)
		return nil, nil
	}

	return []int{corresp[json_data.Orientationt], corresp[json_data.Orientationc], corresp[json_data.Orientationf]}, []string{json_data.TCP_address + ":" + json_data.TCP_port, json_data.Fichiert, json_data.Fichierc, json_data.Fichierf, json_data.Fichierblank, json_data.Fichiercross}
}
