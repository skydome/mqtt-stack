package main

import (
	"fmt"
	broker "github.com/abdulkadiryaman/hrotti/broker"
	"github.com/hashicorp/consul/command/agent"
	"github.com/mitchellh/cli"
	"os"
	"os/signal"
	"syscall"
)

func bootstrapConsul2(dc string, bootstrap bool) {
	var args
	if bootstrap {
		args = []string{"-server", "-bootstrap", "-node", "canawar", "-dc", dc, "-data-dir", "/tmp/consul"}
	}else{
		args = []string{"-server", "-node", "canawar", "-dc", dc, "-data-dir", "/tmp/consul"}
	}
	ui := &cli.BasicUi{Writer: os.Stdout}
	command := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}
	command.Run(args)
}

func bootstrapMqttServer() {
	listener := broker.NewListenerConfig("tcp://0.0.0.0:1883")

	h := broker.NewHrotti(100)

	fmt.Println("starting app")

	h.AddListener("self", listener)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	h.Stop()
}

func main() {
	go bootstrapConsul2("dc1", true)
	bootstrapMqttServer()
}