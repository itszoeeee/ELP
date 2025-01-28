module Ordre exposing (..)

import Parsing exposing(..)
import Parser exposing(..)

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
