package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Hero-exe/mdcheck/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if exitErr, ok := err.(app.ExitCodeError); ok {
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
