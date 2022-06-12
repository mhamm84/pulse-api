package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		uri string
	}
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
	flag.StringVar(&cfg.db.uri, "db_uri", os.Getenv("PULSE_MONGO_URI"), "connection uri of the MongoDB server")
	flag.Parse()

	myFigure := figure.NewColorFigure("Pulse API", "", "green", true)
	myFigure.Print()

	client, err := openMongo(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := application{
		cfg:    cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		app.logger.Printf("Error serving app: %s", err.Error())
	}

}

/*
 * Connect to MongoDB
 */
func openMongo(cfg config) (*mongo.Client, error) {
	// Create a new client and connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.db.uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged MongoDB.")
	return client, nil
}
