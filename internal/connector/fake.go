package connector

import "go.uber.org/zap"

type fakeConnector struct {
	*zap.Logger
}

func (c *fakeConnector) Connect() error {
	c.Logger.Info("fake connect")
	return nil
}
