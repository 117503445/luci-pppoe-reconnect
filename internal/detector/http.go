package detector

import (
	"encoding/json"
	"time"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

type httpDetector struct {
	*zap.Logger
	cfg *httpDetectorConfig
}

type httpDetectorConfig struct {
	URL string
}

const MaxTry = 3

func (d *httpDetector) WaitUntilFailure() {
	var counter int
	for {
		_, err := req.Get(d.cfg.URL)
		if err != nil {
			counter++
			d.Logger.Info("get failed", zap.Int("counter", counter), zap.Error(err))
		} else {
			if counter != 0 {
				d.Logger.Info("get success, clear counter")
			}
			counter = 0
		}
		if counter == MaxTry {
			break
		}
		time.Sleep(time.Minute)
	}
}

func newHTTPDetector(cfg map[string]interface{}, logger *zap.Logger) *httpDetector {
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	var config httpDetectorConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	logger.Info("init http detector", zap.String("url", config.URL))

	return &httpDetector{
		Logger: logger,
		cfg:    &config,
	}
}
