module Structure exposing (..)

type Structure = F Int | L Int | R Int | Repeat Int (List Structure)

type alias Data = List Structure

type alias Ordre = List String

convertToOrdre : Data -> Ordre    -- pour une liste ex [L 4, R 3] ==> [L, L, L, L, R, R, R]
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
