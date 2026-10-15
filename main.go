package main

import (
	"fmt"
	"os"
)

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("utilisation : go run . entree.txt sortie.txt")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("lecture : %w", err)
	}
	result, err := Transform(string(data))
	if err != nil {
		return err
	}
	if err := os.WriteFile(args[1], []byte(result), 0644); err != nil {
		return fmt.Errorf("écriture : %w", err)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
