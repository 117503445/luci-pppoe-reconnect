package connector

import (
	"encoding/json"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type sshConnector struct {
	*zap.Logger
	cfg *sshConnectorConfig
}

type sshConnectorConfig struct {
	Host     string
	Username string
	Password string
}

func (c *sshConnector) Connect() error {
	config := &ssh.ClientConfig{
		User: c.cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.cfg.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", c.cfg.Host, config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.Run("/sbin/ifup wan")
	return err
}

func newSSHConnector(cfg map[string]interface{}, logger *zap.Logger) *sshConnector {
	data, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	var config sshConnectorConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &sshConnector{
		Logger: logger,
		cfg:    &config,
	}
}
