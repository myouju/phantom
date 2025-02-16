package main

import (
	"github.com/tenntenn/typeannotation"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(typeannotation.Analyzer) }
