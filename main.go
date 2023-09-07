package main

import (
	"fmt"
	"os"
)

const (
	exitOK int = 0
	exitNG int = 1
)

func main() {
	cli := &CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	exitCode := exitOK
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintln(cli.Stderr, "Error:", err)
		exitCode = exitNG
	}

	os.Exit(exitCode)
}
