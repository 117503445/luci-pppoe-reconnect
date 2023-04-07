package detector

import "go.uber.org/zap"

type Detector interface {
	WaitUntilFailure()
}

func GetDetector(cfg map[string]interface{}, logger *zap.Logger) Detector {
	name, ok := cfg["type"]
	if !ok {
		panic("detector not found")
	}
	delete(cfg, "type")

	switch name {
	case "http":
		return newHTTPDetector(cfg, logger)
	}
	panic("Unsupported Detector")
}
