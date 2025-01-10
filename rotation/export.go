package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

// Structure pour nos 5 images avec le nom du fichier et l'orientation (blank,up,down,right,left)
type Image struct {
	Fichier     string `json:"image"`
	Orientation string `json:"orientation"`
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
	var image1 Image
	err = json.Unmarshal(bxteValue, &image1) //unmarshal parse le json et stocke le resultat dans &image
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	//copie des deux fichiers en input dans le bon dossier
	copie(image1.Fichier, "../pattern_test/"+image1.Fichier) //si on sait pas ce qu'on nous donne on peut remplacer par image1.Fichier
	copie("blank.png", "../pattern_test/blank.png")

	// Afficher les données pour tester
	//fmt.Printf("Image: %s\n", image0.Fichier)
	//fmt.Printf("Orientation: %s\n", image1.Orientation)

	//Creation des 3 autres images en fonction de image1.Orientation
	var par1, par2, par3 string
	if image1.Orientation == "up" {
		par1 = "down"
		par2 = "right"
		par3 = "left"
	}
	if image1.Orientation == "down" {
		par1 = "up"
		par2 = "left"
		par3 = "right"
	}
	if image1.Orientation == "right" {
		par1 = "left"
		par2 = "down"
		par3 = "up"
	}
	if image1.Orientation == "left" {
		par1 = "right"
		par2 = "up"
		par3 = "down"
	}
	erreur1 := flipImage(image1.Fichier, "../pattern_test/"+par1+".png", 2)
	if erreur1 != nil {
		fmt.Println("Erreur down :", erreur1)
		return
	}

	erreur2 := flipImage(image1.Fichier, "../pattern_test/"+par2+".png", 1)
	if erreur2 != nil {
		fmt.Println("Erreur left :", erreur2)
		return
	}

	erreur3 := flipImage(image1.Fichier, "../pattern_test/"+par3+".png", 3)
	if erreur3 != nil {
		fmt.Println("Erreur right:", erreur3)
		return
	}
}

func copie(Source, Dest string) {

	// Ouvrir le fichier source
	src, err := os.Open(Source)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier source:", err)
		return
	}
	defer src.Close()

	// Créer ou ouvrir le fichier de destination
	dst, err := os.Create(Dest)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier destination:", err)
		return
	}
	defer dst.Close()

	// Copier le contenu de l'image source vers le fichier de destination
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("Erreur lors de la copie du fichier:", err)
		return
	}
}

func flipImage(inputFile, outputFile string, param int) error {
	// Ouvrir le fichier d'entrée
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de l'image : %v", err)
	}
	defer file.Close()

	// Décoder l'image
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("erreur lors du décodage de l'image : %v", err)
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
		return fmt.Errorf("valeur invalide pour param : %d. Les valeurs valides sont 1, 2 ou 3", param)
	}

	// Créer le fichier de sortie
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier de sortie : %v", err)
	}
	defer output.Close()

	// Encoder et sauvegarder l'image retournée
	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(output, rotated)
	case "jpeg":
		err = jpeg.Encode(output, rotated, nil)
	default:
		return fmt.Errorf("format d'image non supporté : %s", format)
	}

	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage de l'image : %v", err)
	}

	return nil
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
