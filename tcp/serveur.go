package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

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

		go handleConnection(conn) // Traiter la connexion dans une goroutine
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	for {
		// Lire l'en-tête pour la taille des données
		var dataSize int32
		err := binary.Read(conn, binary.BigEndian, &dataSize)
		if err != nil {
			fmt.Println("Client déconnecté ou erreur de lecture :", err)
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
	}
}
