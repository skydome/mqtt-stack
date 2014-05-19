package main

import (
	"flag"
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
)

func bootstrapConsul(dc string, bootstrap bool) {
	private, err := consul.GetPrivateIP()
	if err != nil {
		fmt.Errorf("err: %v", err)
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

func Join(args []string) int {
	var wan bool
	log.Println("Joining to cluster with addresses : ", args)
	cmdFlags := flag.NewFlagSet("join", flag.ContinueOnError)
	cmdFlags.BoolVar(&wan, "wan", false, "wan")
	rpcAddr := command.RPCAddrFlag(cmdFlags)
	if err := cmdFlags.Parse(args); err != nil {
		fmt.Errorf("Error occured : %v", err)
		return 1
	}

	addrs := cmdFlags.Args()
	if len(addrs) == 0 {
		log.Fatal("At least one address to join must be specified.")
		return 1
	}

	client, err := command.RPCClient(*rpcAddr)
	if err != nil {
		log.Fatal("Error connecting to Consul agent: %s", err)
		return 1
	}
	defer client.Close()

	n, err := client.Join(addrs, wan)
	if err != nil {
		log.Fatal("Error joining the cluster: %s", err)
		return 1
	}

	log.Println(
		"Successfully joined cluster by contacting %d nodes.", n)
	return 0
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
