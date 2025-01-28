module Parsing exposing(..)

import Parser exposing (..)
type Structure = F Int | L Int | R Int | Repeat Int (List Structure)

type alias Data = List Structure

--Fonction principale
tcTurtleToData : String -> Result (List Parser.DeadEnd) Data
tcTurtleToData input =
    run parserGlobal input

parserGlobal : Parser Data
parserGlobal =
    Parser.sequence
        { start = "["
        , separator = ","
        , end = "]"
        , spaces = spaces
        , item = optionsParser 
        , trailing = Parser.Optional
        }

optionsParser : Parser Structure --On essaie de parse avec un des 4 parser (f,r,l ou repeat)
optionsParser =
    oneOf
        [ fParser
        , lParser
        , rParser
        , repeatParser
        ]

fParser : Parser Structure
fParser =
    succeed F
        |= (symbol "Forward" |> andThen (\_ -> spaces) |> andThen (\_ -> int))

lParser : Parser Structure
lParser =
    succeed L
        |= (symbol "Left" |> andThen (\_ -> spaces) |> andThen (\_ -> int))

rParser : Parser Structure
rParser =
    succeed R
        |= (symbol "Right" |> andThen (\_ -> spaces) |> andThen (\_ -> int))

repeatParser : Parser Structure
repeatParser =
    succeed Repeat
        |= (symbol "Repeat"
                |> andThen (\_ -> spaces)
                |> andThen (\_ -> int)
           )
        |= (spaces
                |> andThen
                    (\_ ->
                        parserGlobal   --"recurence" pour relancer parserGlobal dans le [] du repeat
                    )
           )