package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

var numer_port = 8000
var address = fmt.Sprintf("127.0.0.1:%d", numer_port)

func sendInt(conn net.Conn, prompt string) (int, error) {
	var value int
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)

		tempValue, err := strconv.Atoi(input)
		if err != nil || tempValue <= 0 {
			fmt.Println("Erreur : Veuillez entrer un nombre entier valide supérieur à 0.")
			continue
		}
		value = tempValue
		break
	}

	// Convertir l'entier en bytes
	buffer := make([]byte, 4) // Un int32 nécessite 4 octets
	binary.BigEndian.PutUint32(buffer, uint32(value))

	// Envoyer d'abord la taille des données
	sizeBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBuffer, uint32(len(buffer)))

	_, err := conn.Write(sizeBuffer)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'envoi de la taille : %v", err)
	}

	// Envoyer les données au serveur
	_, err = conn.Write(buffer)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'envoi des données : %v", err)
	}

	// Lire la confirmation du serveur
	confirmationBuffer := make([]byte, 1024)
	n, err := conn.Read(confirmationBuffer)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la lecture de la confirmation : %v", err)
	}
	fmt.Println("Confirmation du serveur :", string(confirmationBuffer[:n]))

	return value, nil
}

func main() {
	// Connexion au serveur
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecté au serveur.")

	Largeur, err := sendInt(conn, "Entrez la largeur de la grille que vous voulez générer (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

	Longueur, err := sendInt(conn, "Entrez la longueur de la grille que vous voulez générer (entier) : ")
	if err != nil {
		fmt.Println(err)
		return
	}

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

	display(matrix, Largeur, Longueur)
}
