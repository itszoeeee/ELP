package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

var matrix = [][]int{{4, 6, 4, -1}, {1, 3, 4, 0}, {3, 1, 1, 3}, {4, 0, 0, 6}}
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

	var conn net.Conn
	for {
		conn, err = listener.Accept() // Accepter une connexion entrante
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation :", err)
			continue
		}
		fmt.Println("Nouvelle connexion acceptée")
		break
	}

	defer conn.Close()
	fmt.Println("Client connecté :", conn.RemoteAddr())

	Largeur, err := receiveInt(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	Longueur, err := receiveInt(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Largeur: %d Longueur: %d", Largeur, Longueur)

	data, err := json.Marshal(matrix)
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

func receiveInt(conn net.Conn) (int, error) {
	// Lire la taille des données
	sizeBuffer := make([]byte, 4)
	_, err := conn.Read(sizeBuffer)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la lecture de la taille : %v", err)
	}
	size := binary.BigEndian.Uint32(sizeBuffer)

	// Lire les données
	dataBuffer := make([]byte, size)
	_, err = conn.Read(dataBuffer)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la lecture des données : %v", err)
	}

	// Convertir les bytes en entier
	value := int(binary.BigEndian.Uint32(dataBuffer))

	// Envoyer une confirmation au client
	confirmationMessage := fmt.Sprintf("Entier reçu : %d", value)
	_, err = conn.Write([]byte(confirmationMessage))
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'envoi de la confirmation : %v", err)
	}

	return value, nil
}
