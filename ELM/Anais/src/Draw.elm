module Draw exposing (..)

import Parsing exposing(..)
import Parser exposing(..)
import Svg exposing (..)
import Svg.Attributes exposing (..)

type alias Ordre = List String

convert: String -> Ordre
convert string=
    convertToOrdre (tcTurtleToData string)

convertToOrdre : Data -> Ordre    -- pour data = [Repeat 2 [R 5, Repeat 3 [F 4, L 1]], F 6] on obtient ["R","R","R","R","R","F","F","F","F","L","F","F","F","F","L","F","F","F","F","L","R","R","R","R","R","F","F","F","F","L","F","F","F","F","L","F","F","F","F","L","F","F","F","F","F","F"]
convertToOrdre data =
    List.concatMap structureToList data


structureToList : Structure -> Ordre   -- pour une seule commande ex structureToList(L 4) ==> ["L","L","L","L"]
structureToList structure =
    case structure of
        A n ->
            List.repeat n "F"

        G n ->
            List.repeat n "L"

        D n ->
            List.repeat n "R"

        Repeat n structures ->
            let
                repeatedStructures =
                    List.concatMap structureToList structures
            in
            List.concat (List.repeat n repeatedStructures)


type Movement
    = F -- Avancer
    | L -- Tourner à gauche
    | R -- Tourner à droite

-- Fonction pour calculer les positions à partir d'une liste de mouvements
-- Définir un degré en radians
degreeToRadian : Float -> Float
degreeToRadian degree =
    degree * pi / 180

-- Fonction pour normaliser un vecteur
normalize : (Float, Float) -> (Float, Float)
normalize (dx, dy) =
    let
        length = sqrt (dx * dx + dy * dy)
    in
    if length > 0 then
        (dx / length, dy / length)
    else
        (dx, dy)

-- Fonction modifiée pour des rotations plus précises
rotate : Movement -> (Float, Float) -> (Float, Float)
rotate movement (dx, dy) =
    case movement of
        L -> 
            let
                angle = degreeToRadian 1
                newDx = cos angle * dx - sin angle * dy
                newDy = sin angle * dx + cos angle * dy
            in
            normalize (newDx, newDy)

        R -> 
            let
                angle = degreeToRadian (-1)
                newDx = cos angle * dx - sin angle * dy
                newDy = sin angle * dx + cos angle * dy
            in
            normalize (newDx, newDy)

        F -> 
            normalize (dx, dy)

calculatePositions : List Movement -> List (Float, Float)
calculatePositions movements =
    let
        -- Fonction pour calculer la position avec un pas constant
        calculate : List (Float, Float) -> (Float, Float) -> (Float, Float) -> List Movement -> List (Float, Float)
        calculate accPositions (x, y) (dx, dy) remainingMovements =
            case remainingMovements of
                [] -> 
                    accPositions

                F :: rest -> 
                    let
                        step = 1.0 
                        (normalizedDx, normalizedDy) = normalize (dx, dy)
                        newPosition = (x + normalizedDx * step, y + normalizedDy * step)
                    in
                    calculate (accPositions ++ [newPosition]) newPosition (dx, dy) rest

                L :: rest -> 
                    let
                        newDirection = rotate L (dx, dy)
                    in
                    calculate accPositions (x, y) newDirection rest

                R :: rest -> 
                    let
                        newDirection = rotate R (dx, dy)
                    in
                    calculate accPositions (x, y) newDirection rest
    in
    calculate [(50, 50)] (50, 50) (0.2, 0) movements


-- Fonction pour transformer une liste de positions en une chaîne formatée pour le SVG
formatPositions : List (Float, Float) -> String
formatPositions posList =
    posList
        |> List.map (\(x, y) -> String.fromFloat x ++ "," ++ String.fromFloat y)
        |> String.join " "

-- Exemple d'utilisation
level_1 : List Movement
level_1 = [ F, F, F, L, F, F, F, R, F, F, F, R, F, F, F ]

positions : String
positions =
    level_1
        |> calculatePositions
        |> formatPositions

convertToMovements : Ordre -> List Movement
convertToMovements ordre =
    List.filterMap (\s ->
        case s of
            "F" -> Just F
            "L" -> Just L
            "R" -> Just R
            _   -> Nothing
    ) ordre

