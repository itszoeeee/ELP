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

	var image0 Image
	image0.Fichier = "blank.png"
	image0.Orientation = "blank"

	//copie des deux fichiers en input dans le bon dossier
	copie("up.png", "../pattern_test/up.png") //si on sait pas ce qu'on nous donne on peut remplacer par image1.Fichier
	copie("blank.png", "../pattern_test/blank.png")

	// Afficher les données pour tester
	//fmt.Printf("Image: %s\n", image0.Fichier)
	//fmt.Printf("Orientation: %s\n", image1.Orientation)

	//Creation des 3 autres images
	par1 := "down"
	inputFile := "up.png"
	outputFile1 := "../pattern_test/down.png"

	erreur1 := flipImage(inputFile, outputFile1, par1)
	if erreur1 != nil {
		fmt.Println("Erreur down :", err)
		return
	}
	par2 := "left"
	outputFile2 := "../pattern_test/left.png"

	erreur2 := flipImage(inputFile, outputFile2, par2)
	if erreur2 != nil {
		fmt.Println("Erreur left :", err)
		return
	}
	par3 := "right"
	outputFile3 := "../pattern_test/right.png"

	erreur3 := flipImage(inputFile, outputFile3, par3)
	if erreur3 != nil {
		fmt.Println("Erreur right:", err)
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

func flipImage(inputFile, outputFile, param string) error {
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

	// Créer une nouvelle image retournée
	bounds := img.Bounds()
	flipped := image.NewRGBA(bounds)

	//operations de retournement
	if param == "down" {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				newY := bounds.Max.Y - (y - bounds.Min.Y) - 1
				flipped.Set(x, newY, img.At(x, y))
			}
		}
	}
	if param == "right" {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				newX := bounds.Max.Y - (y - bounds.Min.Y) - 1
				newY := x - bounds.Min.X
				flipped.Set(newX, newY, img.At(x, y))
			}
		}
	}
	if param == "left" {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				// Calculer la nouvelle position après une rotation de 90 degrés anti-horaire
				newX := y - bounds.Min.Y
				newY := bounds.Max.X - (x - bounds.Min.X) - 1

				// Définir le pixel à la nouvelle position
				flipped.Set(newX, newY, img.At(x, y))
			}
		}
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
		err = png.Encode(output, flipped) //on utilise des png pour l'instant mais utile si fichier d'entree diff
	case "jpeg":
		err = jpeg.Encode(output, flipped, nil)
	default:
		return fmt.Errorf("format d'image non supporté : %s", format)
	}

	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage de l'image : %v", err)
	}

	return nil
}
