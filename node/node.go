package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

const DefaultIP = "127.0.0.1"
const DefaultPort = 8081

type ConnectionInfo struct {
	IP   string
	Port uint64

	connected bool
}

func NewConnectionInfo(ip string, port uint64, connected bool) ConnectionInfo {
	return ConnectionInfo{ip, port, connected}
}

type Node struct {
	// for now we will keep the blockchain in memory and persist it to a DB/persistent storage at a later time
	db   database.Blockchain
	info ConnectionInfo

	peers map[string]ConnectionInfo
}

func NewNode(ip string, port uint64) *Node {
	peers := map[string]ConnectionInfo{}

	node := &Node{database.NewBlockchain(), NewConnectionInfo(ip, port, true), peers}

	return node
}

func (n *Node) Start(ctx context.Context) error {
	return n.serverHttp(ctx)
}

func (n *Node) serverHttp(ctx context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/blocks", n.getBlockChain)
	handler.HandleFunc("/mine", n.mineBlock)

	server := &http.Server{Addr: fmt.Sprintf(":%d", n.info.Port), Handler: handler}
	// This is to ensure the server is closed if/when the context signals done
	go func() {
		<-ctx.Done()
		server.Close()
		fmt.Println("Server was closed.")
	}()

	fmt.Printf("Listening for requests on port: %d.\n", n.info.Port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}
