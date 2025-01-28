module Main exposing (main)

import Browser
import Parser exposing(..)
import Html exposing (Html, div, h1, p, text, input, button)
import Html.Attributes exposing (class, value, placeholder)
import Html.Events exposing (onInput, onClick)
import Parsing exposing(..)

-- MODEL
type alias Model =
    { title : String
    , content : String
    , userInput : String
    , result : String
    }

init : Model
init =
    { title = "Dessiner avec TcTurtle"
    , content = "Entre vos instructions dans le langage TcTurtle"
    , userInput = ""
    , result = ""
    }

-- UPDATE
type Msg
    = UpdateInput String
    | ProcessInput

-- Fonction qui traite l'input (vous pouvez la modifier selon vos besoins)
processString : String -> String
processString str =
    "traitons l'entree : " ++ str

update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newInput ->
            { model | userInput = newInput }
        
        ProcessInput ->
            { model | result = processString model.userInput }

-- VIEW
view : Model -> Html Msg
view model =
    div [ class "container" ]
        [ h1 [ class "title" ] [ text model.title ]
        , p [ class "content" ] [ text model.content ]
        , div [ class "input-section" ]
            [ input 
                [ placeholder "Ex:[Left 3, Forward 10] ou [Repeat 2 [Left 3, Right 5]]"
                , value model.userInput
                , onInput UpdateInput
                , class "input-field"
                ] []
            , button [ onClick ProcessInput, class "process-button" ] 
                [ text "valider" ]
            ]
        , p [ class "result" ] [ text model.result ]
        ]

-- MAIN
main : Program () Model Msg
main =
    Browser.sandbox
        { init = init
        , update = update
        , view = view
        }

