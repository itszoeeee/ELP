module Structure exposing (..)

type Structure = F Int | L Int | R Int | Repeat Int (List Structure)

type alias Data = List Structure

type alias Ordre = List String

convertToOrdre : Data -> Ordre
convertToOrdre data =
    List.concatMap structureToList data


structureToList : Structure -> Ordre
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
