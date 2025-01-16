Voici la grille (grid_TCP) qui sera renvoyée par le serveur TCP, tu me diras si ça te va.

Grille renvoyée par le serveur TCP (pour une matrice 20x20):
[3 8 0 -1 8 6 8 7 5 7 3 1 5 2 5 2 5 7 3 3 2 1 3 5 6 8 2 5 7 4 6 8 7 4 7 4 7 4 2 1 5 0 2 8 7 4 2 3 1 1 3 1 1 1 1 1 4 6 1 3 8 7 1 5 2 4 6 5 7 8 6 3 8 7 3 3 5 7 3 4 6 5 7 3 4 2 3 3 5 2 3 5 2 1 1 1 3 1 4 6 8 0 6 4 2 4 6 5 7 1 5 7 1 3 3 8 6 3 1 3 4 0 0 2 1 4 7 3 1 3 8 2 3 4 6 5 0 6 3 4 5 7 3 1 8 2 1 5 7 4 6 5 6 4 0 7 3 3 1 4 8 6 4 0 2 1 8 0 2 5 0 0 0 2 3 4 2 5 0 2 5 7 5 0 2 3 5 7 4 0 7 8 0 2 5 6 1 3 3 5 8 6 8 0 2 4 7 1 5 0 6 1 3 5 7 3 8 2 1 3 5 0 6 3 5 2 1 3 3 8 7 8 6 8 2 4 2 1 3 4 7 3 8 2 3 4 0 6 4 2 1 4 0 2 1 1 4 0 2 1 6 5 2 4 2 4 7 3 5 2 8 6 3 4 7 8 6 8 6 8 8 7 1 4 6 4 2 1 3 5 2 3 4 6 5 2 3 5 0 2 1 4 7 5 7 5 2 3 5 0 6 5 2 3 8 6 5 0 0 6 8 6 5 0 2 3 1 4 0 0 7 3 5 2 1 3 8 7 3 3 1 8 7 8 6 1 3 5 0 0 2 4 0 6 8 2 5 6 1 4 3 1 5 2 8 7 5 7 8 7 5 2 8 0 2 4 7 8 0 6 6 8 0 6 1 1 3 5 2 5 0 2 1 3 4 6 1 5 0 7]

De plus, j'ai un peu réarrangé le code de WFC.go sous forme de fonction. C'est encore un peu brouillon car je n'ai pas fini de structurer mes fonctions avec la parallélisation (j'ai travaillé sur une fonction plus simple pour tester). Donc ce n'est pas super clair, normalement ça marche mais c'est bancal, je vais améliorer ça.

En revanche, j'ai regroupé toutes les fonctions qui me permettent de générer les images dans une fonction qui s'appelle "client". Tu peux bien sûr regarder dedans, car ça peut t'être utile pour faire le client en fonction de la grille. Je précise que la fonction fait plein d'appels à d'autres fonctions au-dessus que tu peux regarder (c'est pour alléger le code). Bien sûr, toutes ces fonctions vont être supprimées du serveur car il ne traite pas les images. Donc fais comme ça te plaît pour le client. Je précise également, pour ne pas utiliser la structure gridItem (et te simplifier la vie), je ne transmets pas la "grid" que j'utilise mais la "grid_TCP" qui est une version simplifiée (voir ci-dessus). Néanmoins, mon code pour afficher les tuiles utilise encore la structure avec grid (je ne l'ai pas adaptée), donc il faudra adapter la fonction qui affiche les tuiles. Seulement, je sais que tu as ton propre système pour lire le JSON, lire les images et les traiter. Donc fais comme tu veux, à partir du moment où le client affiche la bonne image à partir de la grille qu'il reçoit ci-dessus. Moi, je n'utilise plus les images à présent, donc fais ce qui t'arrange le plus.

Pour utiliser la grille (grid_TCP) (IMPORTANT) :
- Les numéros correspondent aux tuiles (toujours selon la correspondance ci-dessous)
- Les "-1" (j'ai volontairement inséré une valeur à l'indice 3 pour que tu testes) sont lorsque le code n'a pas réussi à générer la tuile pour une raison x ou y (notamment à cause de la parallélisation (voir juste après)). Il faut donc soit rien afficher, soit remplacer par une case noire comme tu veux. Cette information évite de passer par l'attribu collapsed et je pense te facilite la vie.

Je tiens aussi à souligner un problème : il va falloir ajouter une dernière tuile (soit la croix, soit un "tout droit") car même si on avait fait toutes les possibilités, il y a une présence de bug à présent avec le système de parallélisation. Je suggère donc d'ajouter la croix car ça évite de te faire tourner la pièce (pour la 3ème fois). Il faudra donc générer la croix comme la case blanche (croix = 9). Bon, après on pourra le faire plus tard, mais autant anticiper le problème maintenant.

BLANK = 0
T_UP = 1
T_RIGHT = 2
T_DOWN = 3
T_LEFT = 4
C_UP = 5
C_RIGHT = 6
C_DOWN = 7
C_LEFT = 8
CROSS = 9