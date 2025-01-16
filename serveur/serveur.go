package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

var matrix = []int{3, 8, 0, -1, 8, 6, 8, 7, 5, 7, 3, 1, 5, 2, 5, 2, 5, 7, 3, 3, 2, 1, 3, 5, 6, 8, 2, 5, 7, 4,
	6, 8, 7, 4, 7, 4, 7, 4, 2, 1, 5, 0, 2, 8, 7, 4, 2, 3, 1, 1, 3, 1, 1, 1, 1, 1, 4, 6, 1, 3,
	8, 7, 1, 5, 2, 4, 6, 5, 7, 8, 6, 3, 8, 7, 3, 3, 5, 7, 3, 4, 6, 5, 7, 3, 4, 2, 3, 3, 5, 2,
	3, 5, 2, 1, 1, 1, 3, 1, 4, 6, 8, 0, 6, 4, 2, 4, 6, 5, 7, 1, 5, 7, 1, 3, 3, 8, 6, 3, 1, 3,
	4, 0, 0, 2, 1, 4, 7, 3, 1, 3, 8, 2, 3, 4, 6, 5, 0, 6, 3, 4, 5, 7, 3, 1, 8, 2, 1, 5, 7, 4,
	6, 5, 6, 4, 0, 7, 3, 3, 1, 4, 8, 6, 4, 0, 2, 1, 8, 0, 2, 5, 0, 0, 0, 2, 3, 4, 2, 5, 0, 2,
	5, 7, 5, 0, 2, 3, 5, 7, 4, 0, 7, 8, 0, 2, 5, 6, 1, 3, 3, 5, 8, 6, 8, 0, 2, 4, 7, 1, 5, 0,
	6, 1, 3, 5, 7, 3, 8, 2, 1, 3, 5, 0, 6, 3, 5, 2, 1, 3, 3, 8, 7, 8, 6, 8, 2, 4, 2, 1, 3, 4,
	7, 3, 8, 2, 3, 4, 0, 6, 4, 2, 1, 4, 0, 2, 1, 1, 4, 0, 2, 1, 6, 5, 2, 4, 2, 4, 7, 3, 5, 2,
	8, 6, 3, 4, 7, 8, 6, 8, 6, 8, 8, 7, 1, 4, 6, 4, 2, 1, 3, 5, 2, 3, 4, 6, 5, 2, 3, 5, 0, 2,
	1, 4, 7, 5, 7, 5, 2, 3, 5, 0, 6, 5, 2, 3, 8, 6, 5, 0, 0, 6, 8, 6, 5, 0, 2, 3, 1, 4, 0, 0,
	7, 3, 5, 2, 1, 3, 8, 7, 3, 3, 1, 8, 7, 8, 6, 1, 3, 5, 0, 0, 2, 4, 0, 6, 8, 2, 5, 6, 1, 4,
	3, 1, 5, 2, 8, 7, 5, 7, 8, 7, 5, 2, 8, 0, 2, 4, 7, 8, 0, 6, 6, 8, 0, 6, 1, 1, 3, 5, 2, 5,
	0, 2, 1, 3, 4, 6, 1, 5, 0, 7}
var numer_port = 8000
var address = fmt.Sprintf(":%d", numer_port)

func main() {
	listener, err := net.Listen("tcp", address) // Listen pour écouter sur un certain port
	if err != nil {
		fmt.Println("Erreur lors de l'écoute :", err)
		return
	}
	defer listener.Close() // S'assurer que la connexion soit fermée
	fmt.Printf("Serveur en écoute sur le port %d...\n", numer_port)

	for {
		conn, err := listener.Accept() // Accepter une connexion entrante
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation :", err)
			continue
		}
		fmt.Println("Nouvelle connexion acceptée")

		go handleConnection(conn, matrix) // Traiter la connexion dans une goroutine
	}
}

func handleConnection(conn net.Conn, matrix []int) {
	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	for {
		// Lire l'en-tête pour la taille des données
		var dataSize int32
		err := binary.Read(conn, binary.BigEndian, &dataSize)
		if err != nil {
			fmt.Println("Client déconnecté ou erreur de lecture :", err)
			conn.Close()
			break
		}

		data := make([]byte, dataSize)
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println("Erreur lors de la lecture des données :", err)
			break
		}

		// Lire les données selon la taille spécifiée
		if dataSize == 4 { // Vérifier si la taille correspond à un entier encodé
			number := int(binary.BigEndian.Uint32(data))
			fmt.Println("Nombre reçu :", number)
		} else {
			// Si ce n'est pas un entier, interpréter comme une chaîne
			message := string(data)
			fmt.Println("Message reçu :", message)

			// Si le message est "quit", arrêter la boucle
			if message == "quit" {
				fmt.Println("Le client a demandé à se déconnecter.")
				break
			}
		}
		// Répondre au client
		response := fmt.Sprintf("Message traité \n")
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Erreur lors de l'envoi de la réponse :", err)
			break
		}
		data, err = json.Marshal(matrix)
		if err != nil {
			fmt.Println("Erreur lors de la sérialisation JSON :", err)
			return
		}

		// Envoyer les données
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi des données :", err)
			return
		}

		fmt.Println("Matrice envoyée au client.")

	}
}
