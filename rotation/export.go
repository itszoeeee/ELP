package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
)

// Structure pour nos 5 images avec le nom du fichier et l'orientation (blank,up,down,right,left)
type Image struct {
	Fichier     string `json:"image"`
	Orientation *Orientation
}

type Temp struct {
	Fichiert     string `json:"image"`
	Orientationt string `json:"orientation"`
}

// Structure pour une orientation
type Orientation struct {
	Name       string
	Transform1 string
	Transform2 string
	Transform3 string
}

var orientations = map[string]*Orientation{ //equivalent d'un dico pour faire correspondre une orientation a une structure Orientation
	"up":    {"up", "right", "down", "left"},
	"down":  {"down", "left", "up", "right"},
	"right": {"right", "down", "left", "up"},
	"left":  {"left", "up", "right", "down"},
}

func main() {
	// Ouvrir le fichier JSON
	file, err := os.Open("input.json")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer file.Close()

	// Lire le contenu du fichier
	bxteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return
	}

	// Décoder le JSON dans une structure
	var temp1 Temp //creation d'un Temp pour recuperer l'orientation dans le json puis la transformer en struct Orientation
	var image1 Image
	err = json.Unmarshal(bxteValue, &temp1) //unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	image1.Fichier = temp1.Fichiert
	image1.Orientation = orientations[temp1.Orientationt] //remplace "down"(ou autre) par une structure Orientation grace a la var orientations
	if image1.Orientation == nil {
		fmt.Println("Orientation invalide :", image1.Orientation)
		return
	}

	// Créer les 3 autres images
	var imagea, imageb, imagec image.Image
	imagea, err = flipImage(image1.Fichier, 1)
	if err != nil {
		fmt.Println("Erreur Transform1 :", err)
		return
	}

	imageb, err = flipImage(image1.Fichier, 2)
	if err != nil {
		fmt.Println("Erreur Transform2 :", err)
		return
	}

	imagec, err = flipImage(image1.Fichier, 3)
	if err != nil {
		fmt.Println("Erreur Transform3 :", err)
		return
	}
}

func flipImage(inputFile string, param int) (image.Image, error) {
	// Ouvrir le fichier d'entrée
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de l'image : %v", err)
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de l'image : %v", err)
	}

	// Rotation des images en fonction de param (1 pour une rotation, 2 pour deux rotations, ...)
	var rotated image.Image
	switch param {
	case 1: // Rotation de 90 degrés vers la droite
		rotated = rotate90(img)
	case 2: // Rotation de 180 degrés vers la droite
		rotated = rotate90(rotate90(img))
	case 3: // Rotation de 270 degrés (trois rotations de 90 degres)
		rotated = rotate90(rotate90(rotate90(img)))
	default:
		return nil, fmt.Errorf("valeur invalide pour param : %d. Les valeurs valides sont 1, 2 ou 3", param)
	}

	return rotated, nil
}

// Fonction pour effectuer une rotation de 90 degrés vers la droite
func rotate90(img image.Image) image.Image {
	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			newX := bounds.Max.Y - y - 1
			newY := x
			rotated.Set(newX, newY, img.At(x, y))
		}
	}

	return rotated
}
