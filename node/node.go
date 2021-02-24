package node

import (
	"fmt"
	"net/http"
	"strconv"

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

func (n *Node) Start() {
	n.serverHttp()
}

func (n *Node) serverHttp() {
	http.HandleFunc("/blocks", n.GetBlockChain)
	http.HandleFunc("/mine", n.MineBlock)
	err := http.ListenAndServe(":"+strconv.FormatUint(n.info.Port, 10), nil)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("Listening for requests on port: %d.", n.info.Port)
}
