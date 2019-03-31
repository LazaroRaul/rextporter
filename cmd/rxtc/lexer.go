package main

import (
	"github.com/simelo/rextporter/src/rxt"
	"github.com/simelo/rextporter/src/rxt/grammar"
)

func main() {
	grammar.LexTheRxt(&rxt.TokenWriter{}, "LEX")
}
