package main

import (
	"encoding/json"
)

with open('config.json', 'r') as JSON:
    intent = json.load(JSON)