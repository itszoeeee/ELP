module Structure exposing (..)

type Structure = F Int | L Int | R Int | Repeat Int (List Structure)

type alias Data = List Structure

type alias Ordre = List String

convertToOrdre : Data -> Ordre    -- pour data = [Repeat 2 [R 5, Repeat 3 [F 4, L 1]], F 6] on obtient ["R","R","R","R","R","F","F","F","F","L","F","F","F","F","L","F","F","F","F","L","R","R","R","R","R","F","F","F","F","L","F","F","F","F","L","F","F","F","F","L","F","F","F","F","F","F"]
convertToOrdre data =
    List.concatMap structureToList data


structureToList : Structure -> Ordre   -- pour une seule commande ex "L 4"
structureToList structure =
    case structure of
        F n ->
            List.repeat n "F"

        L n ->
            List.repeat n "L"

        R n ->
            List.repeat n "R"

        Repeat n structures ->
            let
                repeatedStructures =
                    List.concatMap structureToList structures
            in
            List.concat (List.repeat n repeatedStructures)
