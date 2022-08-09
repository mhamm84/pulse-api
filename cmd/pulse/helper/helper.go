package helper

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mhamm84/pulse-api/cmd/config"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/utils"
	"os"
	"time"
)

func OpenDB(cfg *config.DbConfig, retryAttempts int, backoff time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	logger := jsonlog.New(os.Stdout, jsonlog.GetLevel("DEBUG"))

	err := utils.Retry(retryAttempts, backoff, func() error {
		d, err := openDB(*cfg, *logger)
		if err != nil {
			logger.PrintInfo("Error trying to open connection to postgres", map[string]interface{}{
				"err": err,
			})
			return err
		}
		logger.PrintInfo("Setting DB handle", nil)
		db = d
		return nil
	})
	return db, err
}

/*
 * Connect to DB
 */
func openDB(cfg config.DbConfig, logger jsonlog.Logger) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.PrintInfo("connecting and pinging postgres", map[string]interface{}{
		"postgres": cfg.Dsn,
	})
	db, err := sqlx.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, err
	}
	logger.PrintInfo("postgres opened", nil)

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	logger.PrintInfo("postgres ping...", nil)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	logger.PrintInfo("postgres ping success", nil)
	return db, nil
}

//func FolderExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
