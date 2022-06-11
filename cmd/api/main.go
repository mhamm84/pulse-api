package main

import (
	"flag"
	"github.com/common-nighthawk/go-figure"
	"log"
	"os"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	cfg    config
	logger *log.Logger
}

func main() {
	var cfg config

	// Initialize a new logger which writes messages to the standard out stream, // prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	flag.IntVar(&cfg.port, "port", 9091, "Pulse API port number")
	flag.StringVar(&cfg.env, "env", "dev", "dev|stg|uat|prod")
	flag.Parse()

	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	app := application{
		cfg:    cfg,
		logger: logger,
	}

	err := app.server()

	panic(err)

}
