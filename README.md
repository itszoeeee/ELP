# ELP
Wave function collapse WFC

Amélioration :
- Faire des connecteurs -> 
- Optimer la bande passante calcul image sur client

2 possibilités pour paralléliser :
- Générer plusieurs images en modifiant le critère aléatoire blanc / routes.
- Sous-diviser le problème en générant une matrice de bloc où chaque bloc est une image à construire (en parallèle). Il faut d'abord générer les bordures de chaque bloc, pour ensuite recomposer les blocs correctement.
- 