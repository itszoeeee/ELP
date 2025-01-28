# ELM

# Tests
## Level 1
- TcTurtle : [Forward 3]
- Data : F 3
- Ordre : [F, F, F]

## Level 2
- TcTurtle : [Left 3, Right 5]
- Data : [L 3, R 5]
- Ordre : [L, L, L, R, R, R, R, R]

## Level 3
- TcTurtle : [Repeat 2 [Forward 3]]
- Data : [Repeat 2[F 3]]
- Ordre : [F, F, F, F, F, F]

## Level 4
- TcTurtle : [Repeat 2 [Left 3, Right 5]]
- Data : [Repeat 2 [L 3, R 5]] 
- Ordre : [L, L, L, R, R, R, R, R, L, L, L, R, R, R, R, R]

## Level 5
- TcTurtle : [Repeat 2 [Right 5, Repeat 3 [Forward 4, Left 1]], Forward 6]
- Data :  [Repeat 2 [R 5, Repeat 3 [F 4, L 1]], F 6]
- Ordre : [R, R, R, R, R, F, F, F, F, L, F, F, F, F, L, F, F, F, F, L, R, R, R, R, R, F, F, F, F, L, F, F, F, F, L, F, F, F, F, L, F, F, F, F, F, F]
# Objectifs
- Anaïs (Parsing) : Prompt textuel TcTurtle --> Structure
- Zoé (Engine) : Structure --> Ordre
- Mathis (Display) : Ordre --> Dessin

# Types
- TcTurtle : chaîne de texte
- Data : liste de tuples : (Instructions, entier)
- Instructions : type d'entiers :
  F = 0
  L = 1
  R = 2
  Repeat = Data
- Ordre : liste d'Instructions (sans Repeat)
