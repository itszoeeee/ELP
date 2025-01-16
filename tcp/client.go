package main

import (
	"encoding/binary"
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
		// Lire un message depuis l'utilisateur
		fmt.Print("Entrez la dimension de la grille que vous voulez generer: ")
		var message string
		fmt.Scanln(&message)

		// Convertir la chaîne en entier
		number, err := strconv.Atoi(message)
		if err != nil {
			fmt.Println("Erreur : Veuillez entrer un nombre valide.")
			return
		}

		buffer := make([]byte, 4) // Un int32 nécessite 4 octets
		binary.BigEndian.PutUint32(buffer, uint32(number))

		// Calculer la taille des données (ici, 4 octets pour l'entier)
		dataSize := int32(len(buffer))
		sizeBuffer := make([]byte, 4) // 4 octets pour la taille
		binary.BigEndian.PutUint32(sizeBuffer, uint32(dataSize))

		// Envoyer d'abord la taille des données
		_, err = conn.Write(sizeBuffer)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi de la taille :", err)
			break
		}

		// Envoyer les données au serveur
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi des données :", err)
			break
		}

		// Si le message est "quit", fermer la connexion
		if message == "quit" {
			fmt.Println("Déconnexion...")
			break
		}

		// Lire la réponse du serveur
		buffer = make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erreur lors de la lecture de la réponse :", err)
			break
		}

		fmt.Println("Réponse du serveur :", string(buffer[:n]))
	}
}
