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
	var sens string
	if image1.Orientation == "up" {
		par1 = "down"
		par2 = "right"
		par3 = "left"
		sens = "vertical"
	}
	if image1.Orientation == "down" {
		par1 = "up"
		par2 = "left"
		par3 = "right"
		sens = "vertical"
	}
	if image1.Orientation == "right" {
		par1 = "left"
		par2 = "down"
		par3 = "up"
		sens = "horizontal"
	}
	if image1.Orientation == "left" {
		par1 = "right"
		par2 = "up"
		par3 = "down"
		sens = "horizontal"
	}
	erreur1 := flipImage(image1.Fichier, "../pattern_test/"+par1+".png", 1, sens)
	if erreur1 != nil {
		fmt.Println("Erreur down :", erreur1)
		return
	}

	erreur2 := flipImage(image1.Fichier, "../pattern_test/"+par2+".png", 2, sens)
	if erreur2 != nil {
		fmt.Println("Erreur left :", erreur2)
		return
	}

	erreur3 := flipImage(image1.Fichier, "../pattern_test/"+par3+".png", 3, sens)
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

func flipImage(inputFile, outputFile string, param int, sens string) error {
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
	if param == 1 && sens == "vertical" { //rotation 180 degres, retournement vertical
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				newY := bounds.Max.Y - (y - bounds.Min.Y) - 1
				flipped.Set(x, newY, img.At(x, y))
			}
		}
	}
	if param == 1 && sens == "horizontal" { //rotation 180 degres, retournement horizontal
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				newX := bounds.Max.X - (x - bounds.Min.X) - 1
				flipped.Set(newX, y, img.At(x, y))
			}
		}
	}
	if param == 2 { // rotation de 90 degrés anti-horaire
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				newX := bounds.Max.Y - (y - bounds.Min.Y) - 1
				newY := x - bounds.Min.X
				flipped.Set(newX, newY, img.At(x, y))
			}
		}
	}
	if param == 3 { // rotation de 90 degrés anti-horaire
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				newX := y - bounds.Min.Y
				newY := bounds.Max.X - (x - bounds.Min.X) - 1
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
