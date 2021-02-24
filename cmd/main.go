package main

import (
	"flag"

	"github.com/ooga-mon/blockchain/node"
)

func main() {
	flagIP := flag.String("ip", node.DefaultIP, "your node's public IP in the P2P network")
	flagPort := flag.Uint64("port", node.DefaultPort, "your node's public port in the P2P network")
	flag.Parse()

	node := node.NewNode(*flagIP, *flagPort)
	node.Start()
}
