package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter    *rate.Limiter
		lastActive time.Time
	}

	var (
		mux     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mux.Lock()

			for ip, client := range clients {
				if time.Since(client.lastActive) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mux.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.cfg.limiter.enabled {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mux.Lock()
			if _, found := clients[host]; !found {
				clients[host] = &client{limiter: rate.NewLimiter(rate.Limit(app.cfg.limiter.rps), app.cfg.limiter.burst)}
			}

			clients[host].lastActive = time.Now()

			if !clients[host].limiter.Allow() {
				mux.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mux.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
