package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type prompt_dataItem struct {
	Dim_x    int
	Dim_y    int
	Proba    int
	Div_x    int
	Div_y    int
	NbWorker int
}

type clientProgress struct {
	nb_cell_collapsed int64
	stopChan          chan struct{}
}

var numer_port = 8000
var address = fmt.Sprintf(":%d", numer_port)

func receive_data(conn net.Conn, data *prompt_dataItem) {
	// Lecture de la taille des données JSON
	buffer_size := make([]byte, 4)
	if _, err := conn.Read(buffer_size); err != nil {
		fmt.Println("Erreur lors de la lecture de la taille :", err)
		return
	}
	data_size := binary.BigEndian.Uint32(buffer_size)

	// Lecture des données sérialisées de la taille spécifiée
	data_serial := make([]byte, data_size)
	if _, err := conn.Read(data_serial); err != nil {
		fmt.Println("Erreur lors de la lecture des données :", err)
		return
	}
	// Désérialisation des données JSON
	if err := json.Unmarshal(data_serial, data); err != nil {
		fmt.Println("Erreur lors de la désérialisation des données :", err)
		return
	}
	fmt.Println("Données reçues avec succès du client :", conn.RemoteAddr())
}

func send_data(conn net.Conn, data [][]int) {
	// Sérialisation en JSON de data
	data_serial, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erreur de sérialisation JSON :", err)
		return
	}
	// Envoyer la taille des données JSON
	buffer_size := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer_size, uint32(len(data_serial)))
	if _, err := conn.Write(buffer_size); err != nil {
		fmt.Println("Erreur lors de l'envoi de la taille des données :", err)
		return
	}
	// Envoyer les données sérialisées
	if _, err := conn.Write(data_serial); err != nil {
		fmt.Println("Erreur lors de l'envoi des données :", err)
		return
	}
	fmt.Println("Données envoyées avec succès au client :", conn.RemoteAddr())
}

func send_int(conn net.Conn, data int) {
	// Convertir l'entier en bytes
	buffer := make([]byte, 4) // Un int32 nécessite 4 octets
	binary.BigEndian.PutUint32(buffer, uint32(data))

	// Envoyer du pourcentage
	_, err := conn.Write(buffer)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi du pourcentage :", err)
		return
	}
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close() // Ferme la connexion client
	defer wg.Done()    // Ferme la tâche associée au client

	client_progress := &clientProgress{
		nb_cell_collapsed: 0,
		stopChan:          make(chan struct{}),
	}

	// --- Reception des données ---
	var prompt_data prompt_dataItem // Structure pour stocker les données à recevoir
	receive_data(conn, &prompt_data)

	var grid_TCP [][]int
	dim_x := prompt_data.Dim_x
	dim_y := prompt_data.Dim_y
	proba := prompt_data.Proba
	div_x := prompt_data.Div_x
	div_y := prompt_data.Div_y
	nbWorkers := prompt_data.NbWorker

	if dim_x > 0 || dim_y > 0 || div_x > 0 || div_y > 0 || nbWorkers > 0 {
		stopChan := make(chan struct{}) // Créer un canal pour arrêter le rapporteur de progression
		go progress(client_progress.stopChan, conn, dim_x, dim_y, &client_progress.nb_cell_collapsed)

		grid_process(&grid_TCP, dim_x, dim_y, proba, div_x, div_y, nbWorkers, &client_progress.nb_cell_collapsed) // Calcul de la grille

		close(stopChan) // Ferme le canal de progression
		fmt.Println("Génération de la grille terminée pour le client :", conn.RemoteAddr())

		close(client_progress.stopChan)
		send_data(conn, grid_TCP) // Envoie de la grille au client
	} else {
		fmt.Println("Données reçues invalides ou incomplètes pour :", conn.RemoteAddr())
	}

}

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
		// Incrementer le waitGroup pour chaque client et lancer une goroutine pour gerer la connexion
		wg.Add(1)
		go handleClient(conn, &wg)
	}
}
