package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/117503445/luci-pppoe-reconnect/internal/cfg"
	"github.com/117503445/luci-pppoe-reconnect/internal/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"

	"golang.org/x/crypto/ssh"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	logger, err := log.NewAtLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}

	defer func() {
		err = logger.Sync()
	}()

	cfg.InitConfig()

	logger.Info("Hello world!", zap.String("connector.name", viper.GetString("connector.name")))

	config := &ssh.ClientConfig{
		User: viper.GetString("connector.username"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("connector.password")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", viper.GetString("connector.host"), config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run("/sbin/ifup wan"); err != nil {
		panic(err)
	}
	fmt.Println(buff.String())

	return err
}
