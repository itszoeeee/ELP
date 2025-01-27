module Parsing exposing(..)

type Structure = F Int | L Int | R Int | Repeat Int (List Structure)

type alias Data = List Structure

-- Fonction qui fait la traduction de tcTurtle a la liste de Structure
tcTurtleToData : String -> Result String Data
tcTurtleToData input =
    case parseInput input of
        Ok structures -> Ok structures
        Err error -> Err error

parseInput : String -> Result String Data --retourne ok Data si ca marche et sinon un string error  
parseInput input = parseStructures (String.trim input) --String.trim pour supprimer les espaces

parseStructures : String -> Result String Data
parseStructures str =
    if String.isEmpty str then
        Ok []
    else if String.startsWith "[" str && String.endsWith "]" str then
        let
            inner = str
                |> String.dropLeft 1
                |> String.dropRight 1
                |> String.trim
        in
        parseRepeatOrStructures inner

    else
        str
            |> splitTopLevelCommaSeparated
            |> List.map (String.trim >> parseStructureString) 
            |> combineResults

-- New function to handle complex comma splitting
splitTopLevelCommaSeparated : String -> List String
splitTopLevelCommaSeparated str =
    let
        splitHelper : List Char -> Int -> List Char -> List String
        splitHelper chars bracketDepth currentPart =
            case chars of
                [] ->
                    if List.isEmpty currentPart then
                        []
                    else
                        [ String.fromList (List.reverse currentPart) ]

                '[' :: rest ->
                    splitHelper rest (bracketDepth + 1) ('[' :: currentPart)

                ']' :: rest ->
                    splitHelper rest (bracketDepth - 1) (']' :: currentPart)

                ',' :: rest ->
                    if bracketDepth == 0 then
                        String.fromList (List.reverse currentPart) 
                        :: splitHelper rest bracketDepth []
                    else
                        splitHelper rest bracketDepth (',' :: currentPart)

                c :: rest ->
                    splitHelper rest bracketDepth (c :: currentPart)
    in
    splitHelper (String.toList str) 0 []

parseRepeatOrStructures : String -> Result String Data
parseRepeatOrStructures str =
    let
        parseRepeatHelper parts =
            case parts of
                countStr :: content ->
                    case String.toInt (String.trim countStr) of
                        Just repeatCount ->
                            let
                                cleanContent = 
                                    String.join " " content 
                                    |> String.trim
                                    |> String.replace "]" ""
                            in
                            case parseStructures cleanContent of
                                Ok nested ->
                                    Ok [ Repeat repeatCount nested ]
                                Err err ->
                                    Err err

                        Nothing ->
                            Err ("Invalid repeat count: " ++ String.join " " parts)

                _ ->
                    Err "Invalid Repeat format"
    in
    if String.contains "Repeat" str then
        str
            |> String.replace "Repeat" ""
            |> String.trim
            |> String.split "["
            |> List.filter (\s -> s /= "")
            |> List.map String.trim
            |> parseRepeatHelper
    else
        parseStructures str

-- Parse une seule structure (F, L, R) à partir d'une chaîne
parseStructureString : String -> Result String Structure
parseStructureString str =
    case String.words (String.trim (String.replace "]" "" str)) of
        [ "Forward", count ] ->
            String.toInt count
                |> Maybe.map F
                |> Result.fromMaybe ("Invalid count for Forward: " ++ count)

        [ "Left", count ] ->
            String.toInt count
                |> Maybe.map L
                |> Result.fromMaybe ("Invalid count for Left: " ++ count)

        [ "Right", count ] ->
            String.toInt count
                |> Maybe.map R
                |> Result.fromMaybe ("Invalid count for Right: " ++ count)

        _ ->
            Err ("Unknown structure: " ++ str)

-- Combine une liste de résultats en un seul
combineResults : List (Result String Structure) -> Result String Data
combineResults results =
    results
        |> List.foldr combine (Ok [])

-- Combine deux résultats
combine : Result String Structure -> Result String (List Structure) -> Result String (List Structure)
combine currentResult acc =
    case (currentResult, acc) of
        (Ok item, Ok items) ->
            Ok (item :: items)

        (Err err, _) ->
            Err err

        (_, Err err) ->
            Err err