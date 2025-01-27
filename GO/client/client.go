package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type prompt_dataItem struct {
	dim_x    int
	dim_y    int
	rand     int
	div_x    int
	div_y    int
	nbWorker int
}

func promptInt(prompt string) (int, error) {
	var value int
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)

		tempValue, err := strconv.Atoi(input)
		if err != nil || tempValue <= 0 {
			fmt.Println("Erreur : Veuillez entrer un nombre entier positif valide.")
			continue
		}
		value = tempValue
		break
	}
	return value, nil
}

func prompt(prompt_data *prompt_dataItem) {
	var err error
	prompt_data.dim_x, err = promptInt("Entrez la largeur de la grille que vous voulez générer (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	prompt_data.dim_y, err = promptInt("Entrez la hauteur de la grille que vous voulez générer (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	prompt_data.rand, err = promptInt("Entrez le pourcentage de cases vides ou avec des routes que vous souhaitez (entier entre 0 et 100) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Vérifier si la valeur est inférieure à 100

	prompt_data.div_x, err = promptInt("Entrez la division sur la largeur de la grille que vous voulez paralléliser (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Vérifier que la division sur x est plus petite que la dimension sur x

	prompt_data.div_y, err = promptInt("Entrez la division sur la hauteur de la grille que vous voulez paralléliser (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Vérifier que la division sur y est plus petite que la dimension sur y

	prompt_data.nbWorker, err = promptInt("Entrez le nombre de threads maximal que vous voulez utiliser (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendInt(conn net.Conn, value int) error {
	// Convertir l'entier en bytes
	buffer := make([]byte, 4) // Un int32 nécessite 4 octets
	binary.BigEndian.PutUint32(buffer, uint32(value))

	// Envoyer d'abord la taille des données
	sizeBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuffer, uint32(len(buffer)))

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la taille : %v", err)
	}

	// Envoyer les données au serveur
	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi des données : %v", err)
	}

	// Lire la confirmation du serveur
	confirmationBuffer := make([]byte, 1024)
	n, err := conn.Read(confirmationBuffer)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de la confirmation : %v", err)
	}
	fmt.Println("Confirmation du serveur :", string(confirmationBuffer[:n]))

	return nil
}

func main() {
	// --- Connexion au serveur ---
	_, address := lecture_json("input.JSON")
	conn, err := net.Dial("tcp", address[0])
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecté au serveur.")

	// --- Récupération des données via la prompt ---
	var prompt_data prompt_dataItem // Structure pour stocker les données à envoyer
	prompt(&prompt_data)            // Enregistrement des données via le prompt

	// --- Envoie des données pour générer la grille ---
	// err = sendInt(conn, Largeur)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// --- Affichage de la progression ---

	// --- Récupération de la grille ---
	matrixBuffer := make([]byte, 8192) // Taille plus grande pour contenir une matrice
	n, err := conn.Read(matrixBuffer)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la matrice :", err)
		return
	}
	// Désérialiser la matrice
	var matrix [][]int
	err = json.Unmarshal(matrixBuffer[:n], &matrix)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation JSON :", err)
		return
	}

	// --- Exportation de l'image de sortie ---
	display(matrix, prompt_data.dim_x, prompt_data.dim_y)
}
