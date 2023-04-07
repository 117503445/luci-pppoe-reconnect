package checker

import (
	"time"

	"github.com/117503445/openwrt-wan-reconnect/internal/connector"
	"github.com/117503445/openwrt-wan-reconnect/internal/detector"
	"go.uber.org/zap"
)

type Checker struct {
	logger *zap.Logger
	d      detector.Detector
	c      connector.Connector
}

func (c *Checker) StartCheck() {
	for {
		c.logger.Info("wait until network fail")
		c.d.WaitUntilFailure()
		c.logger.Info("try to connect")
		err := c.c.Connect()
		if err != nil {
			c.logger.Error("connect failed", zap.String("err", err.Error()))
		} else {
			time.Sleep(time.Minute)
		}
	}
}

func NewChecker(cfg map[string]map[string]interface{}, logger *zap.Logger) *Checker {
	logger = logger.Named("checker")

	dCfg, ok := cfg["detector"]
	if !ok {
		panic("detector not found")
	}
	cCfg, ok := cfg["connector"]
	if !ok {
		panic("connector not found")
	}

	return &Checker{
		logger: logger,
		d:      detector.GetDetector(dCfg, logger.Named("detector")),
		c:      connector.GetConnector(cCfg, logger.Named("connector")),
	}
}
