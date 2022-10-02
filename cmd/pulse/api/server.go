package api

import (
	"context"
	"fmt"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.Port),
		Handler:      http.TimeoutHandler(app.routes(), 5*time.Second, "timeout limit of request reached"),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx := context.TODO()
		utils.Logger(ctx).Info("signal caught",
			zap.String("signal", s.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		utils.Logger(ctx).Info("completing background tasks...",
			zap.String("addr", srv.Addr),
		)
		app.wg.Wait()
		shutdownError <- nil
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	utils.Logger(context.TODO()).Info("stopped Pulse server",
		zap.String("addr", srv.Addr),
		zap.String("env", app.cfg.Env),
	)

	return nil
}
