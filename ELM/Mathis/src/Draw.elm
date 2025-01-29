module Draw exposing (main)

import Svg exposing (..)
import Svg.Attributes exposing (..)

type Movement
    = F -- Avancer
    | L -- Tourner à gauche
    | R -- Tourner à droite

-- Fonction pour calculer les positions à partir d'une liste de mouvements
calculatePositions : List Movement -> List (Int, Int)
calculatePositions movements =
    let
        -- Fonction pour gérer les rotations
        rotate : Movement -> (Int, Int) -> (Int, Int)
        rotate movement (dx, dy) =
            case movement of
                L -> 
                    -- Tourner à gauche de 90° (anti-horaire)
                    ( -dy, dx )

                R -> 
                    -- Tourner à droite de 90° (horaire)
                    ( dy, -dx )

                F -> 
                    -- Ne modifie pas la direction
                    (dx, dy)

        -- Fonction pour calculer la position
        calculate : List (Int, Int) -> (Int, Int) -> (Int, Int) -> List Movement -> List (Int, Int)
        calculate accPositions (x, y) (dx, dy) remainingMovements =
            case remainingMovements of
                [] ->
                    accPositions

                F :: rest ->
                    -- Avancer en fonction de la direction
                    let
                        newPosition = (x + dx, y + dy)
                    in
                    calculate (accPositions ++ [newPosition]) (x + dx, y + dy) (dx, dy) rest

                L :: rest ->
                    -- Tourner à gauche
                    let
                        newDirection = rotate L (dx, dy)
                    in
                    calculate accPositions (x, y) newDirection rest

                R :: rest ->
                    -- Tourner à droite
                    let
                        newDirection = rotate R (dx, dy)
                    in
                    calculate accPositions (x, y) newDirection rest
    in
    -- On commence à la position (0, 0) avec une direction (1, 0) (vers la droite)
    calculate [(50, 50)] (50, 50) (5, 0) movements

-- Fonction pour transformer une liste de positions en une chaîne formatée pour le SVG
formatPositions : List (Int, Int) -> String
formatPositions posList =
    posList
        |> List.map (\(x, y) -> String.fromInt x ++ "," ++ String.fromInt y)
        |> String.join " "

-- Exemple d'utilisation
level_1 : List Movement
level_1 = [ F, F, F, L, F, F, F, R, F, F, F, R, F, F, F ]

positions : String
positions =
    level_1
        |> calculatePositions
        |> formatPositions

-- Affichage SVG
main =
  svg
    [ width "600"
    , height "600"
    , viewBox "0 0 120 120"
    ]
    [ rect
        [ x "10"
        , y "10"
        , width "100"
        , height "100"
        , rx "15"
        , ry "15"
        , fill "white"
        , strokeWidth "3"
        , stroke "rgb(50, 200, 200)"
        ]
        []
      , polyline
        [ points positions
        , fill "none"
        , strokeWidth "1"
        , stroke "black"
        ]
        []
    ]