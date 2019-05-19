package main

import (
	"fmt"
	"os"

	"go.smartmachine.io/go-grpc-api/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		_,_ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}