package main

import (
	"github.com/myouju/phantom"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(phantom.AssignableAnalyzer)
}
