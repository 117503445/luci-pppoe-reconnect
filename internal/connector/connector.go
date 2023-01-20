package connector

import "go.uber.org/zap"

type Connector interface {
	Connect() error
}

func GetConnector(name string, logger *zap.Logger) Connector {
	switch name {
	case "fake":
		return &fakeConnector{
			Logger: logger,
		}
	case "ssh":
		return &sshConnector{
			Logger: logger,
		}
	}
	panic("Unsupported Connector")
}
