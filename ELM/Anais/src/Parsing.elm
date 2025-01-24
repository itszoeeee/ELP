module Parsing exposing(..)

type Structure = Instruction Int
type Instruction = Repeat [Structure]| Forward | Droite | Gauche