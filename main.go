package main

import (
	broker "github.com/abdulkadiryaman/hrotti/broker"
	"github.com/hashicorp/consul/command/agent"
	"github.com/mitchellh/cli"
	"os"
	"os/signal"
	"syscall"
	"github.com/hashicorp/consul/command"
	"fmt"
	"net"
)

func bootstrapConsul(dc string, bootstrap bool) {
	private, err := GetPrivateIP()
	if err != nil {
		fmt.Fprintf("err: %v", err)
	}

	var agentArgs []string
	if bootstrap {
		agentArgs = []string{"-server", "-bootstrap", "-node", private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	}else {
		agentArgs = []string{"-server", "-node", private.String(), "-dc", dc, "-data-dir", "/tmp/consul"}
	}
	ui := &cli.BasicUi{Writer: os.Stdout}
	agentCommand := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}
	agentCommand.Run(agentArgs)

	joinArgs := []string{"172.17.0.2"}

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

var privateBlocks []*net.IPNet

// Returns if the given IP is in a private block
func isPrivateIP(ip_str string) bool {
	ip := net.ParseIP(ip_str)
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// GetPrivateIP is used to return the first private IP address
// associated with an interface on the machine
func GetPrivateIP() (net.IP, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
	}

	// Find private IPv4 address
	for _, rawAddr := range addresses {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}

		if ip.To4() == nil {
			continue
		}
		if !isPrivateIP(ip.String()) {
			continue
		}

		return ip, nil
	}

	return nil, fmt.Errorf("No private IP address found")
}

func main() {
	go bootstrapConsul("dc1", true)
	bootstrapMqttServer()
}
