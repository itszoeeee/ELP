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

func main() {
	// Connexion au serveur
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecté au serveur.")
	for {
		// Lire une entrée utilisateur valide
		fmt.Print("Entrez la dimension de la grille que vous voulez générer (entier valide) : ")
		var message string
		fmt.Scanln(&message)

		// Convertir la chaîne en entier
		number, err := strconv.Atoi(message)
		if err != nil || number <= 0 {
			fmt.Println("Erreur : Veuillez entrer un nombre entier valide supérieur à 0.")
			continue
		}

		// Convertir l'entier en bytes
		buffer := make([]byte, 4) // Un int32 nécessite 4 octets
		binary.BigEndian.PutUint32(buffer, uint32(number))

		// Envoyer d'abord la taille des données
		dataSize := int32(len(buffer))
		sizeBuffer := make([]byte, 4)
		binary.BigEndian.PutUint32(sizeBuffer, uint32(dataSize))

		_, err = conn.Write(sizeBuffer)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi de la taille :", err)
			return
		}

		// Envoyer les données au serveur
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi des données :", err)
			return
		}

		// Lire la confirmation du serveur
		confirmationBuffer := make([]byte, 1024)
		n, err := conn.Read(confirmationBuffer)
		if err != nil {
			fmt.Println("Erreur lors de la lecture de la confirmation :", err)
			return
		}
		fmt.Println("Confirmation du serveur :", string(confirmationBuffer[:n]))
		break
	}

	matrixBuffer := make([]byte, 8192) // Taille plus grande pour contenir une matrice
	n, err := conn.Read(matrixBuffer)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la matrice :", err)
		return
	}
	// Désérialiser la matrice

	var matrix []int
	err = json.Unmarshal(matrixBuffer[:n], &matrix)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation JSON :", err)
		return
	}

	// Transformer en matrice 20x20
	const DIM = 20
	grid := make([][]int, DIM)
	for j := 0; j < DIM; j++ {
		grid[j] = make([]int, DIM)
		for i := 0; i < DIM; i++ {
			grid[j][i] = matrix[i*DIM+j]
		}
	}

	client(grid)
}
