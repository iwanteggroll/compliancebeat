package main

import (
	"os"

	"github.com/iwanteggroll/compliancebeat/cmd"

	_ "github.com/iwanteggroll/compliancebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
