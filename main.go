package main

import (
	"github.com/mvo5/osbuild-cvet/clienterrorsif"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(clienterrorsif.Analyzer)
}
