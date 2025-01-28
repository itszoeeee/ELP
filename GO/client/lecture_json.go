package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type jsonItem struct {
	TCP_address   string `json:"TCP_address"`
	TCP_port      string `json:"TCP_port"`
	Image_t       string `json:"image_t"`
	Orientation_t string `json:"orientation_t"`
	Image_c       string `json:"image_c"`
	Orientation_c string `json:"orientation_c"`
	Image_f       string `json:"image_f"`
	Orientation_f string `json:"orientation_f"`
	Image_blank   string `json:"image_blank"`
	Fichiercross  string `json:"image_cross"`
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

	if json_data.Orientation_t == "" || json_data.Orientation_c == "" || json_data.Orientation_f == "" {
		fmt.Println("Orientation invalide :", json_data.Orientation_c, json_data.Orientation_t, json_data.Orientation_f)
		return nil, nil
	}

	return []int{corresp[json_data.Orientation_t], corresp[json_data.Orientation_c], corresp[json_data.Orientation_f]}, []string{json_data.Image_t, json_data.Image_c, json_data.Image_f, json_data.Image_blank, json_data.Fichiercross, json_data.TCP_address + ":" + json_data.TCP_port}
}
