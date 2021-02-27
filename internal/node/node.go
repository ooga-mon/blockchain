package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ooga-mon/blockchain/internal/database"
)

const DefaultIP = "127.0.0.1"
const DefaultPort = 8080
const DefaultPeerIP = "127.0.0.1"
const DefaultPeerPort = 8080

type connectionInfo struct {
	IP   string `json:"IP"`
	Port uint64 `json:"Port"`

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
	db    database.Blockchain
	info  connectionInfo
	state state

	peers map[string]connectionInfo
}

func NewNode(ip string, port uint64, bootstrapPeerIP string, bootstrapPeerPort uint64) *Node {
	peers := map[string]connectionInfo{}
	if bootstrapPeerIP != ip || bootstrapPeerPort != port {
		fmt.Printf("bootstrap parms IP:%s, Port:%d\n", bootstrapPeerIP, bootstrapPeerPort)
		bootstrapCon := newConnectionInfo(bootstrapPeerIP, bootstrapPeerPort, false)
		peers[bootstrapCon.tcpAddress()] = bootstrapCon
	}

	state := newState()
	conInfo := newConnectionInfo(ip, port, true)

	node := &Node{database.NewBlockchain(), conInfo, state, peers}

	return node
}

func (n *Node) Mine(tx []database.Transaction) database.Block {
	prevBlock := n.db.GetLastBlock()
	pendingBlock := newPendingBlock(prevBlock.BlockHash, prevBlock.Content.Number+1, tx)
	newBlock := n.mineBlock(pendingBlock)
	n.state.setLastMineTime()
	fmt.Printf("Difficulty set at %d\n", n.state.curDifficulty)
	n.db.AddBlock(newBlock)
	return newBlock
}

func (n *Node) addPeer(con connectionInfo) {
	if !n.containsPeer(con) {
		fmt.Printf("Found new peer %s\n", con.tcpAddress())
		n.peers[con.tcpAddress()] = con
	}
}

func (n *Node) containsPeer(peer connectionInfo) bool {
	if n.info.IP == peer.IP && n.info.Port == peer.Port {
		return true
	}

	_, found := n.peers[peer.tcpAddress()]

	return found
}

func (n *Node) removePeer(peer connectionInfo) {
	delete(n.peers, peer.tcpAddress())
}

func (n *Node) Start(ctx context.Context) error {
	go n.postSync()

	return n.serverHttp(ctx)
}

func (n *Node) serverHttp(ctx context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/blocks", n.handlerGetBlockChain)
	handler.HandleFunc("/mine", n.handlerMineBlock)
	handler.HandleFunc("/node/status", n.handlerPostStatus)

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
