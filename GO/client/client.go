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
	Dim_x    int
	Dim_y    int
	Proba    int
	Div_x    int
	Div_y    int
	NbWorker int
}

func promptInt(prompt string) int {
	var value int
	for {
		fmt.Print(prompt)
		var input string
		fmt.Scanln(&input)

		temp_value, err := strconv.Atoi(input)
		if err != nil || temp_value <= 0 {
			fmt.Println("Erreur : Veuillez entrer un nombre entier positif valide.")
			continue
		}
		value = temp_value
		break
	}
	return value
}

func prompt(prompt_data *prompt_dataItem) {
	prompt_data.Dim_x = promptInt("Entrez la largeur de la grille que vous voulez générer (entier) : ")
	prompt_data.Dim_y = promptInt("Entrez la hauteur de la grille que vous voulez générer (entier) : ")

	prompt_data.Proba = promptInt("Entrez le pourcentage de cases vides ou avec des routes que vous souhaitez (entier entre 0 et 100) : ")
	if prompt_data.Proba > 100 { // Vérifie si la valeur est inférieure à 100
		fmt.Println("Erreur : la probabilité de tirer une case vide doit être inférieure à 100")
		os.Exit(1)
	}

	prompt_data.Div_x = promptInt("Entrez la division sur la largeur de la grille que vous voulez paralléliser (entier) : ")
	if prompt_data.Div_x > prompt_data.Dim_x { // Vérifie que la division sur x est plus petite que la dimension sur x
		fmt.Println("Erreur : la division sur la largeur doit être plus petite que la dimension sur la largeur")
		os.Exit(1)
	}
	if prompt_data.Dim_x/prompt_data.Div_x < 7 { // Vérifie que la taille d'une sous-grille est plus grande que 1 (avec des bordures de 3)
		fmt.Printf("Erreur : la dimension sur x doit être plus grande que %d pour une division de %d sur la largeur\n", 7*prompt_data.Div_x, prompt_data.Div_x)
		os.Exit(1)
	}

	prompt_data.Div_y = promptInt("Entrez la division sur la hauteur de la grille que vous voulez paralléliser (entier) : ")
	if prompt_data.Div_y > prompt_data.Dim_y { // Vérifie que la division sur y est plus petite que la dimension sur y
		fmt.Println("Erreur : la division sur la longueur doit être plus petite que la dimension sur la longueur")
		os.Exit(1)
	}
	if prompt_data.Dim_y/prompt_data.Div_y < 7 { // Vérifie que la taille d'une sous-grille est plus grande que 1 (avec des bordures de 3)
		fmt.Printf("Erreur : la dimension sur y doit être plus grande que %d pour une division de %d sur la longueur\n", 7*prompt_data.Div_y, prompt_data.Div_y)
		os.Exit(1)
	}

	prompt_data.NbWorker = promptInt("Entrez le nombre de threads maximal que vous voulez utiliser (entier) : ")
}

func send_data(conn net.Conn, data prompt_dataItem) {
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
	fmt.Println("Données envoyées avec succès.")
}

func receive_data(conn net.Conn, data *[][]int) {
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
	fmt.Println("\nDonnées reçues avec succès.")
}

func receive_int(conn net.Conn, data *int) {
	// Lecture des données sérialisées de la taille spécifiée
	buffer := make([]byte, 4) // Un int32 nécessite 4 octets
	if _, err := conn.Read(buffer); err != nil {
		fmt.Println("Erreur lors de la réception du pourcentage :", err)
		return
	}
	// Convertir les bytes en entier
	*data = int(binary.BigEndian.Uint32(buffer))
}

func main() {
	// --- Connexion au serveur ---
	_, address := lecture_json("input.JSON")
	conn, err := net.Dial("tcp", address[5])
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecté au serveur.")

	// --- Récupération des données via la prompt ---
	var prompt_data prompt_dataItem // Structure pour stocker les données à envoyer
	prompt(&prompt_data)            // Enregistrement des données via le prompt

	// --- Envoie des données ---
	send_data(conn, prompt_data)

	// --- Affichage de la progression ---
	var percentage int
	for percentage < 100 {
		receive_int(conn, &percentage) // Reception de la progression
		fmt.Printf("\r[")              // Permet de revenir au début de la ligne sans en ajouter une nouvelle, afin de mettre à jour la même ligne de la console
		for j := 0; j < 50; j++ {
			if j < percentage/2 {
				fmt.Print("=")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Printf("] %d%%", percentage)
	}

	// --- Récupération de la grille ---
	var grid [][]int
	receive_data(conn, &grid)

	// --- Exportation de l'image de sortie ---
	var erreur bool = false
	display(grid, prompt_data.Dim_x, prompt_data.Dim_y, &erreur)
	if erreur {
		fmt.Println("Une erreur est survenue dans la génération de l'image, veuillez essayer d'autres paramètres.")
	}
}
