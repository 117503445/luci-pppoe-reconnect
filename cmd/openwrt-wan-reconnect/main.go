package main

import (
	"fmt"
	"os"
	"time"

	"github.com/117503445/openwrt-wan-reconnect/internal/cfg"
	"github.com/117503445/openwrt-wan-reconnect/internal/connector"
	"github.com/117503445/openwrt-wan-reconnect/internal/detector"
	"github.com/117503445/openwrt-wan-reconnect/internal/log"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	logger, err := log.NewAtLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}

	defer func() {
		err = logger.Sync()
	}()

	cfg.InitConfig()

	d := detector.GetDetector("http", logger)
	c := connector.GetConnector("ssh", logger)

	for {
		logger.Info("wait until network fail")
		d.WaitUntilFailure()
		logger.Info("try to connect")
		err := c.Connect()
		if err != nil {
			logger.Error("connect failed", zap.String("err", err.Error()))
		} else {
			time.Sleep(time.Minute)
		}
	}
}
