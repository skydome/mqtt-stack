package main

import (
	broker "github.com/abdulkadiryaman/hrotti/broker"
	"github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/command/agent"
	consul "github.com/hashicorp/consul/consul"
	"github.com/mitchellh/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func bootstrapConsul(dc string, bootstrap bool) {
	private, err := consul.GetPrivateIP()
	if err != nil {
		log.Fatal("err: %v", err)
	}

	var agentArgs []string
	if bootstrap {
		agentArgs = []string{"-server", "-bootstrap", "-node", private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	} else {
		agentArgs = []string{"-server", "-node", private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	}
	ui := &cli.BasicUi{Writer: os.Stdout}
	agentCommand := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}
	agentCommand.Run(agentArgs)

	joinArgs := []string{"192.168.2.90"}

	joinCommand := &command.JoinCommand{
		Ui: ui,
	}
	joinCommand.Run(joinArgs)

}

func bootstrapMqttServer() {
	listener := broker.NewListenerConfig("tcp://0.0.0.0:1883")

	h := broker.NewHrotti(100)

	h.AddListener("self", listener)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	h.Stop()
}

func main() {
	go bootstrapConsul("dc1", true)
	bootstrapMqttServer()
}
