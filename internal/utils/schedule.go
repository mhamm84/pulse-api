package utils

import (
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ScheduleTaskRunner struct {
	ticker time.Ticker
	delay  time.Duration
	quit   chan int
	logger *jsonlog.Logger
}

func NewScheduleTaskRunner(initialDelay time.Duration, delay time.Duration, logger *jsonlog.Logger) ScheduleTaskRunner {
	return ScheduleTaskRunner{
		ticker: *time.NewTicker(initialDelay),
		delay:  delay,
		quit:   make(chan int),
		logger: logger,
	}
}

func (tr *ScheduleTaskRunner) Start(task func()) error {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer func() {
			tr.stopTheClock()
			tr.logger.PrintDebug("ScheduleTaskRunner stopped", nil)
		}()
		firstEx := true
		for {
			select {
			case <-tr.ticker.C:
				if firstEx {
					tr.ticker.Stop()
					tr.ticker = *time.NewTicker(tr.delay)
					firstEx = false
				}
				tr.logger.PrintDebug("ScheduleTaskRunner running task", nil)
				go task()
				break

			case <-tr.quit:
				return

			case <-sig:
				tr.Close()
				break
			}
		}
	}()
	return nil
}

func (tr *ScheduleTaskRunner) Close() {
	go func() {
		tr.quit <- 1
	}()
}

func (tr *ScheduleTaskRunner) stopTheClock() {
	tr.ticker.Stop()
}
