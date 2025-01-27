package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// var old_matrix = [][]int{{4, 6, 4, -1}, {1, 3, 4, 0}, {3, 1, 1, 3}, {4, 0, 0, 6}}
var matrix = [][]int{{10, 9, 1, 5, 6, 9, 3, 3, 3, 10, 9, 1, 1, 10, 4}, {8, 6, 8, 7, 8, 2, 1, 9, 1, 3, 1, 10, 8, 0, 2}, {4, 7, 4, 6, 9, 5, 7, 5, 0, 2, 8, 7, 4, 0, 6}, {9, 4, 6, 3, 1, 3, 9, 8, 7, 5, 6, 1, 4, 0, 7}, {1, 1, 10, 5, 7, 9, 5, 6, 5, 7, 10, 8, 6, 3, 5}, {8, 0, 7, 10, 1, 1, 10, 3, 10, 4, 0, 2, 8, 2, 10}, {4, 0, 2, 3, 3, 8, 7, 1, 8, 6, 3, 5, 6, 1, 3}, {6, 3, 5, 2, 5, 2, 4, 0, 6, 3, 5, 0, 7, 10, 9}, {0, 6, 3, 9, 10, 9, 4, 7, 10, 1, 3, 3, 1, 3, 1}, {0, 0, 2, 1, 3, 4, 2, 5, 7, 10, 9, 1, 3, 4, 0}, {10, 3, 1, 10, 5, 6, 5, 7, 9, 3, 4, 7, 5, 2, 10}, {7, 5, 7, 8, 7, 8, 0, 2, 5, 2, 5, 6, 3, 4, 7}, {9, 3, 1, 5, 6, 1, 10, 4, 7, 4, 7, 8, 6, 5, 2}, {2, 5, 0, 7, 8, 0, 7, 1, 5, 6, 5, 2, 10, 8, 11}, {6, 8, 0, 2, 5, 0, 2, 8, 0, 0, 0, 6, 8, 6, 4}}
var numer_port = 8000
var address = fmt.Sprintf(":%d", numer_port)

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Erreur lors de l'écoute :", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Serveur en écoute sur le port %d...\n", numer_port)

	// WaitGroup pour attendre que toutes les goroutines se terminent
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation :", err)
			continue
		}
		fmt.Println("Nouvelle connexion acceptée de", conn.RemoteAddr())
		//incrementer le waitGroup pour chaque client et lancer une goroutine pour gerer la connexion
		wg.Add(1)
		go handleClient(conn, &wg)
	}
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	Largeur, err := receiveInt(conn)
	if err != nil {
		fmt.Println("Erreur de réception de la largeur :", err)
		return
	}

	Longueur, err := receiveInt(conn)
	if err != nil {
		fmt.Println("Erreur de réception de la longueur :", err)
		return
	}

	fmt.Printf("Client %s - Largeur: %d Longueur: %d\n", conn.RemoteAddr(), Largeur, Longueur)

	data, err := json.Marshal(matrix)
	if err != nil {
		fmt.Println("Erreur lors de la sérialisation JSON :", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi des données :", err)
		return
	}

	fmt.Printf("Matrice envoyée au client %s.\n", conn.RemoteAddr())
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
