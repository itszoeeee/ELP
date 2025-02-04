module Main exposing (main)

import Browser
import Parser exposing(..)
import Html exposing (Html, div, h1, p, text, input, button)
import Html.Attributes exposing (class, value, placeholder)
import Html.Events exposing (onInput, onClick)
import Svg exposing (svg, polyline, rect)
import Svg.Attributes exposing (width, height, viewBox, points, fill, strokeWidth, stroke, x, y, rx, ry)
import Draw exposing(..)

-- MODEL
type alias Model =
    { title : String
    , content : String
    , userInput : String
    , result : String
    , svgPath : String
    }

init : Model
init =
    { title = "Dessiner avec TcTurtle"
    , content = "Entrez vos instructions dans le langage TcTurtle"
    , userInput = ""
    , result = ""
    , svgPath = ""
    }

-- UPDATE
type Msg
    = UpdateInput String
    | ProcessInput

result : String -> String
result str =
    let
        parsedMovements = convert str
    in
    formatPositions (calculatePositions (convertToMovements parsedMovements))


update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newInput ->
            { model
                | userInput = newInput
            }
        
        ProcessInput ->
            let
                parsedMovements = convert model.userInput
                svgPathString = formatPositions (calculatePositions (convertToMovements parsedMovements))
            in
            { model
                | result = "Instructions traitées."
                , svgPath = svgPathString
            }

-- VIEW
view : Model -> Html Msg
view model =
    div [ class "container" ]
        [ h1 [ class "title" ] [ text model.title ]
        , p [ class "content" ] [ text model.content ]
        , div [ class "input-section" ]
            [ input 
                [ placeholder "Ex: [Left 3, Forward 10] ou [Repeat 2 [Left 3, Right 5]]"
                , value model.userInput
                , onInput UpdateInput
                , class "input-field"
                ] []
            , button [ onClick ProcessInput, class "process-button" ] 
                [ text "Valider" ]
            ]
        , p [ class "result" ] [ text model.result ]
        , svg
            [ width "600"
            , height "600"
            , viewBox "0 0 120 250"
            ]
            [ rect
                [ x "10"
                , y "10"
                , rx "15"
                , ry "15"
                , width "150"
                , height "150"
                , fill "white"
                , strokeWidth "3"
                , stroke "rgb(200, 50, 162)"
                ]
                []
            , polyline
                [ points model.svgPath
                , fill "none"
                , strokeWidth "0.5"
                , stroke "black"
                ]
                []
            ]
        ]


-- MAIN
main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }

