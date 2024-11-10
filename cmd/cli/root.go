package main

import (
	"os"
)

func main() {
	rootCmd.Flags().
		StringVarP(&targetModel, "model", "m", "", "Specify the model")

	// Execute the root command
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
