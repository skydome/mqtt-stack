package main

import (
	"fmt"
	broker "github.com/abdulkadiryaman/hrotti/broker"
	command "github.com/hashicorp/consul/command"
	"github.com/hashicorp/consul/command/agent"
	consul "github.com/hashicorp/consul/consul"
	"github.com/mitchellh/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func bootstrapConsul(dc string, bootstrap bool) {
	private, err := consul.GetPrivateIP()
	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	var agentArgs []string
	if bootstrap {
		agentArgs = []string{"-server", "-bootstrap", "-node", "skydome-mqtt-" + private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	} else {
		agentArgs = []string{"-server", "-node", "skydome-mqtt-" + private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	}
	ui := &cli.BasicUi{Writer: os.Stdout}
	agentCommand := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}
	agentCommand.Run(agentArgs)
}

func Join(args []string) {
	log.Println("Sleeping for 3 seconds to consul come up....")
	time.Sleep(3000 * time.Millisecond)
	ui := &cli.BasicUi{Writer: os.Stdout}
	joinCommand := &command.JoinCommand{
		Ui: ui,
	}
	joinCommand.Run(args)
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
	var bootstrap bool
	if len(os.Args) < 2 {
		bootstrap = true
	} else {
		bootstrap = os.Args[1] != "false"
	}

	go bootstrapConsul("dc1", bootstrap)

	if !bootstrap {
		joinArgs := []string{"172.17.0.2"}
		go Join(joinArgs)
	}

	bootstrapMqttServer()
}
