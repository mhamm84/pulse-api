package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	// Setup logging
	logger := jsonlog.New(os.Stdout, jsonlog.GetLevel("DEBUG"))

	// Fancy ascii splash when starting the app
	myFigure := figure.NewColorFigure("Pulse API Admin CLI", "", "green", true)
	myFigure.Print()

	logger.PrintInfo("Hi! Welcome!", nil)

	var pulseCmd = &cobra.Command{
		Use:   "pulse",
		Short: "Pulse API CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello Cobra CLI")
		},
	}

	pulseCmd.AddCommand()

	err := pulseCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	fmt.Println("I'm inside initConfig function in cmd/root.go")
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
