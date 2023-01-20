package detector

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type httpDetector struct {
	*zap.Logger
}

func (d *httpDetector) WaitUntilFailure() {
	var counter int
	for {
		_, err := http.Get("https://www.baidu.com")
		if err != nil {
			counter++
			d.Logger.Info("get failed", zap.Int("counter", counter))
		} else {
			if counter != 0 {
				d.Logger.Info("get success, clear counter")
			}
			counter = 0
		}
		if counter == 3 {
			break
		}
		time.Sleep(time.Minute)
	}

}
