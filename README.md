# ELP
Wave function collapse WFC

Amélioration :
- Faire des connecteurs pour ne pas hard coder les règles de connexions entre les options
- Optimiser la bande passante du calcul image sur le client
    - Le client envoie une requête de création de la grille
    - Le client charge les images / Le serveur process la grille (en fonction du curseur de random)
    - Le serveur retourne la grille process (en parallèle)
    - Le client génère l'image

2 possibilités pour paralléliser :
- Générer plusieurs images en modifiant le critère aléatoire blanc / routes.
- Sous-diviser le problème en générant une matrice de bloc où chaque bloc est une image à construire (en parallèle). Il faut d'abord générer les bordures de chaque bloc, pour ensuite recomposer les blocs correctement.

Remarque : la fonction (boucle principale) accélère si on supprime l'affichage