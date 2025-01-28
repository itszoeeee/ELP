module Web exposing (main)

import Browser
import Html exposing (div)

import Svg exposing (..)
import Svg.Attributes --exposing (..)

init = 0

view model =
  div [] [text "Test oh oh"]

update model = model

main =
  Browser.sandbox
    {
      init = init,
      view = view,
      update = update
    }