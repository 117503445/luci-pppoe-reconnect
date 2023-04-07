package connector

import "go.uber.org/zap"

type Connector interface {
	Connect() error
}

func GetConnector(cfg map[string]interface{}, logger *zap.Logger) Connector {
	name, ok := cfg["type"]
	if !ok {
		panic("connector not found")
	}
	delete(cfg, "type")

	switch name {
	case "fake":
		return &fakeConnector{
			Logger: logger,
		}
	case "ssh":
		return newSSHConnector(cfg, logger)
	case "clash":
		return newClashConnector(cfg, logger)
	}

	panic("Unsupported Connector")
}
