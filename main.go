package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/consul"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"io/ioutil"
	"log"
	broker "github.com/abdulkadiryaman/hrotti/broker"
)

func tmpDir() string {
	dir, err := ioutil.TempDir("", "consul")
	if err != nil {
		log.Fatal("err: %v", err)
	}
	return dir
}
func getConfiguration() (string, *consul.Config) {
	dir := tmpDir()
	config := consul.DefaultConfig()
	config.Bootstrap = true
	config.Datacenter = "dc1"
	config.DataDir = dir

	// Adjust the ports
	rpcPort := 8300
	config.NodeName = fmt.Sprintf("Node %d", rpcPort)
	config.RPCAddr = &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: rpcPort,
	}
	lanPort := 8301
	lanIp := "127.0.0.1"
	config.SerfLANConfig.MemberlistConfig.BindAddr = lanIp
	config.SerfLANConfig.MemberlistConfig.BindPort = lanPort
	config.SerfLANConfig.MemberlistConfig.SuspicionMult = 2
	config.SerfLANConfig.MemberlistConfig.ProbeTimeout = 50*time.Millisecond
	config.SerfLANConfig.MemberlistConfig.ProbeInterval = 100*time.Millisecond
	config.SerfLANConfig.MemberlistConfig.GossipInterval = 100*time.Millisecond


	private, err := consul.GetPrivateIP()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	wanPort := 8302
	config.SerfWANConfig.MemberlistConfig.BindAddr = private.String()
	config.SerfWANConfig.MemberlistConfig.BindPort = wanPort
	config.SerfWANConfig.MemberlistConfig.SuspicionMult = 2
	config.SerfWANConfig.MemberlistConfig.ProbeTimeout = 50*time.Millisecond
	config.SerfWANConfig.MemberlistConfig.ProbeInterval = 100*time.Millisecond
	config.SerfWANConfig.MemberlistConfig.GossipInterval = 100*time.Millisecond

	config.RaftConfig.HeartbeatTimeout = 40*time.Millisecond
	config.RaftConfig.ElectionTimeout = 40*time.Millisecond

	config.ReconcileInterval = 100*time.Millisecond
	return dir, config
}


func bootstrapConsul(dc string, bootstrap bool) (string, *consul.Server) {
	dir, config := getConfiguration()
	config.Datacenter = dc
	config.Bootstrap = bootstrap
	server, err := consul.NewServer(config)
	if err != nil {
		log.Fatal("err: %v", err)
	}

	server.RPC()
	private, err := consul.GetPrivateIP()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	time.Sleep(10 * time.Second)
	// Try to join
	log.Println("Joining to : ", private)
	if _, err := server.JoinWAN([]string{private.String()}); err != nil {
		log.Fatal("err: %v", err)
	}
	log.Println("Joined to : ", private)
	return dir, server
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
	go bootstrapConsul("dc1", true)
	bootstrapMqttServer()
}
