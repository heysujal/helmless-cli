package main

import (
	"fmt"
	"os"

	"github.com/helmless/helmless-cli/cmd/helmless"
)

func main() {
	if err := helmless.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
