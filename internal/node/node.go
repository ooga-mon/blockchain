package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

const DefaultIP = "127.0.0.1"
const DefaultPort = 8081
const DefaultPeerIP = "127.0.0.1"
const DefaultPeerPort = 8081

type connectionInfo struct {
	IP   string
	Port uint64

	connected bool
}

func newConnectionInfo(ip string, port uint64, connected bool) connectionInfo {
	return connectionInfo{ip, port, connected}
}

func (con *connectionInfo) tcpAddress() string {
	return fmt.Sprintf("%s:%d", con.IP, con.Port)
}

type Node struct {
	// for now we will keep the blockchain in memory and persist it to a DB/persistent storage at a later time
	db   database.Blockchain
	info connectionInfo

	peers map[string]connectionInfo
}

func NewNode(ip string, port uint64, bootstrapPeerIP string, bootstrapPeerPort uint64) *Node {
	peers := map[string]connectionInfo{}
	if bootstrapPeerIP != ip && bootstrapPeerPort != port {
		bootstrapCon := newConnectionInfo(bootstrapPeerIP, bootstrapPeerPort, false)
		peers[bootstrapCon.tcpAddress()] = bootstrapCon
	}

	node := &Node{database.NewBlockchain(), newConnectionInfo(ip, port, true), peers}

	return node
}

func (n *Node) addPeer(con connectionInfo) {
	n.peers[con.tcpAddress()] = con
}

func (n *Node) Start(ctx context.Context) error {
	return n.serverHttp(ctx)
}

func (n *Node) serverHttp(ctx context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/blocks", n.handlerGetBlockChain)
	handler.HandleFunc("/mine", n.handlerMineBlock)
	handler.HandleFunc("/node/peer", n.handlerAddPeer)

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
