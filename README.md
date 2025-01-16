# ELP
Wave function collapse WFC

ALGORITHME 

Entropie = nombre d'options disponibles pour remplir une case.
On choisit une case de la grille avec un nombre d'options minimum (la première choisie est donc une case random) et on la remplit aléatoirement. Une fois remplie, elle est "collapsed" et on ne la modifiera plus. 
Puis, on prend une autre case qui n'a pas encore été "collapsed", avec aussi une entropie minimum (normalement une voisine de la première). On remplit cette case en prenant en compte ses quatre cases voisines, puis on la marque comme "collapsed". 
Pour remplir une case correctement, on observe toutes ses cases voisines, avec chacune qui lui donne une liste d'options possibles. Donc au final on prend une option d'image qui appartient aux quatre listes en même temps (si plusieurs choix prendre aléatoirement), et ainsi de suite jusqu'à remplir toute la grille.


DOCSTRINGS DES DIFFERENTES FONCTIONS

- rotation(img, int) :
Prend en paramètres une image appelée par son src path, et un entier qui indique à l'algorithme l'orientation de la pièce choisie parmi 1 = "T_up", 2 = "T_right", 3 = "T_down", 4 = "T_left", 5 = "C_up", 6 = "C_right", 7 = "C_down", 8 = "C_left"
Renvoie une liste de toutes les pièces disponibles, y compris blank (0), dans le bon ordre à partir de la pièce d'entrée

- loadImage :
importe une image à partir d'un chemin source et la décode/l'affiche

- createTiles() (Tiles[9]image.Image, err error) : 
Tiles[] est la liste renvoyée par rotation (avec 9 pièces différentes)
Renvoie l'image correspondant à chaque pièce possible selon l'indice d'orientation correspondant 

- createEmptyImage(int width, int height)
Crée une image vide de dimension width*height

- placeImageInMatrix(dst *image.RGBA, src image.Image, gridX, gridY, cellSize int)
prend en paramètre une image et la dessine dans la case de taille cellSize, d'indice (gridX,gridY)

- checkValid(option *[]int, valid []int) :
Prend en paramètres une liste des neuf options possibles et une liste des options valables en fonction des cases voisines
Utilise un pointeur pour trouver des options valables en fonction des quatre voisins


AMELIORATION

Régler la fréquence d'apparition du blank en lui attribuant une probabilité supérieure lors du tirage aléatoire
/!\ Pourcentage de probabilité à régler en fonction du nombre d'options disponibles

- func weightedRandom(items []WeightedItem)
Prend en paramètres une liste de WeightedItem (structure {value, weight} )
Renvoie un tirage aléatoire d'une des value, avec une probabilité d'apparition weight 

- func proba(liste []string)
Prend en paramètres la liste des options disponibles pour remplir chaque case et la proba qu'on veut appliquer à blank (p = 50 par ex, en %)
Si blank est une option, renvoie une liste de WeightedItem avec un poids de p pour blank et (100-p)/(n-1) pour les autres options (avec n le nombre d'options disponibles valables)

On applique ensuite a cette liste de valeurs pondérées la fonction WeightedRandom


PARALLELISATION

on subdivise la grille pour faire les sous-grilles en meme temps, avec colonnes et lignes intermédiaires pour les rassembler 
ex : grille de 75*75 
on génère indépendemment (1-37)(1-37), (39-75)(1-37), (1-37)(39-75), (39-75)(39-75)
puis pour les rassembler : on génère la colonne 38 grâce aux colonnes 37 et 39, et on génère la ligne 38 grâce aux lignes 37 et 39 

