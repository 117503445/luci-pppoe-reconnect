package detector

import "go.uber.org/zap"

type Detector interface {
	WaitUntilFailure()
}

func GetDetector(name string, logger *zap.Logger) Detector {
	switch name {
	case "http":
		return &httpDetector{
			Logger: logger,
		}
	}
	panic("Unsupported Detector")
}
