package connector

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

type clashConnector struct {
	*zap.Logger
	cfg *clashConnectorConfig
}

type clashConnectorConfig struct {
	Host     string
	Token    string
	Selector string
}

func (c *clashConnector) Connect() error {
	// https://clash.gitbook.io/doc/restful-api/proxies
	// https://req.cool/zh/docs/prologue/introduction/

	c.Logger.Info("clash connect")
	resp := req.SetBearerAuthToken(c.cfg.Token).MustGet(fmt.Sprintf("%s/proxies", c.cfg.Host))
	var data map[string]map[string]map[string]interface{}
	err := json.Unmarshal([]byte(resp.String()), &data)
	if err != nil {
		return err
	}
	names := []string{}
	for _, proxy := range data["proxies"] {
		t := proxy["type"].(string)
		if t == "Vmess" || t == "Trojan" {
			names = append(names, proxy["name"].(string))
		}
	}

	minDelay := math.MaxInt
	minDelayProxy := ""

	for _, name := range names {
		resp := req.SetBearerAuthToken(c.cfg.Token).SetQueryParams(map[string]string{
			"timeout": "3000",
			"url":     "http://google.com",
		}).MustGet(fmt.Sprintf("%s/proxies/%s/delay", c.cfg.Host, name))
		var data map[string]int
		err := json.Unmarshal([]byte(resp.String()), &data)
		if err == nil && data["delay"] != 0 {
			delay := data["delay"]
			if delay < minDelay {
				minDelay = delay
				minDelayProxy = name
			}
		}
	}

	if minDelayProxy != "" {
		c.Logger.Debug("check proxy", zap.String("name", minDelayProxy))
		req.SetBearerAuthToken(c.cfg.Token).SetBody(map[string]string{"name": minDelayProxy}).MustPut(fmt.Sprintf("%s/proxies/%s", c.cfg.Host, c.cfg.Selector))
	}
	return nil
}

func newClashConnector(cfg map[string]interface{}, logger *zap.Logger) *clashConnector {
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	var config clashConnectorConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &clashConnector{
		Logger: logger,
		cfg:    &config,
	}
}
