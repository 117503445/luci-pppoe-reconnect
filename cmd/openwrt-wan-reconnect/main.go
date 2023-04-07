package main

import (
	"fmt"
	"os"

	// "time"

	"github.com/117503445/openwrt-wan-reconnect/internal/cfg"
	"github.com/117503445/openwrt-wan-reconnect/internal/checker"

	// "github.com/117503445/openwrt-wan-reconnect/internal/connector"
	// "github.com/117503445/openwrt-wan-reconnect/internal/detector"
	"github.com/117503445/openwrt-wan-reconnect/internal/log"
	"github.com/spf13/viper"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	// "go.uber.org/zap"
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

	checkersCfg := viper.GetStringMap("checkers")
	// var checkers []*checker.Checker
	for name, c := range checkersCfg {
		cc := c.(map[string]interface{})
		checkerCfg := map[string]map[string]interface{}{
			"detector":  cc["detector"].(map[string]interface{}),
			"connector": cc["connector"].(map[string]interface{}),
		}
		logger.Info("init checker", zap.String("name", name))
		checker := checker.NewChecker(checkerCfg, logger)
		go checker.StartCheck()
		// checkers = append(checkers, )
	}

	select {}
}
