package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("Starting [%s] Pulse server at Addr: %v", app.cfg.env, srv.Addr)

	err := srv.ListenAndServe()

	return err
}
