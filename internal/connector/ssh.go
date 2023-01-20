package connector

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type sshConnector struct {
	*zap.Logger
}

func (c *sshConnector) Connect() error {
	config := &ssh.ClientConfig{
		User: viper.GetString("connector.username"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("connector.password")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", viper.GetString("connector.host"), config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if err := session.Run("/sbin/ifup wan"); err != nil {
		return err
	}

	return nil
}
