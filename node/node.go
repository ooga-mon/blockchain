package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ooga-mon/blockchain/database"
)

const DefaultIP = "127.0.0.1"
const DefaultPort = 8080
const DefaultPeerIP = "127.0.0.1"
const DefaultPeerPort = 8080
const MiningInterval = time.Second * 10

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
	// For now we will keep all the data in memory. We can do persistent storage as a later feature.
	db    database.Blockchain
	info  connectionInfo
	state state

	txPool map[string]database.SignedTransaction
	peers  map[string]connectionInfo
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
	txPool := map[string]database.SignedTransaction{}

	node := &Node{database.NewBlockchain(), conInfo, state, txPool, peers}

	return node
}

func (n *Node) Start(ctx context.Context) error {
	go n.mine(ctx)
	go n.postSync()

	return n.serverHttp(ctx)
}

func (n *Node) mine(ctx context.Context) {
	ticker := time.NewTicker(MiningInterval)

	for {
		select {
		case <-ticker.C:
			if !n.state.isMining {
				n.state.setMiningFlag(true)
				go n.mineTxPool()
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (n *Node) mineTxPool() {
	signedTxs := n.getNextTransactionToMine()
	if len(signedTxs) > 0 {
		prevBlock := n.db.GetLastBlock()
		pendingBlock := newPendingBlock(prevBlock.BlockHash, prevBlock.Content.Number+1, signedTxs)
		newBlock := n.mineBlock(pendingBlock)
		n.state.setLastMineTime()
		fmt.Printf("Mined new block %s at difficulty %d\n", newBlock.BlockHash.Hex(), n.state.curDifficulty)
		n.db.AddBlock(newBlock)
		n.removePendingTransactions(pendingBlock.tx)
	}
	n.state.setMiningFlag(false)
}

func (n *Node) getNextTransactionToMine() []database.SignedTransaction {
	signedTxs := make([]database.SignedTransaction, len(n.txPool))

	for _, tx := range n.txPool {
		signedTxs = append(signedTxs, tx)
	}

	return signedTxs
}

func (n *Node) addPendingTransaction(tx database.SignedTransaction) error {
	if ok, err := tx.IsAuthentic(); !ok {
		fmt.Println("Non-authentic transaction came in. Not adding to pool.")
		return err
	}

	signedTxHash, err := tx.Hash()
	if err != nil {
		return err
	}

	signedTxHex := signedTxHash.Hex()

	if _, found := n.txPool[signedTxHex]; !found {
		n.txPool[signedTxHex] = tx
	}

	fmt.Println("Added transaction to pool.")
	return nil
}

func (n *Node) removePendingTransactions(txs []database.SignedTransaction) {
	for _, signedTx := range txs {
		hash, err := signedTx.Hash()
		if err != nil {
			continue
		}
		delete(n.txPool, hash.Hex())
	}
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

func (n *Node) serverHttp(ctx context.Context) error {
	handler := http.NewServeMux()

	handler.HandleFunc("/node/blocks", n.handlerGetBlockChain)
	handler.HandleFunc("/node/transaction", n.handlerPostTransaction)
	handler.HandleFunc("/node/status", n.handlerPostStatus)

	server := &http.Server{Addr: fmt.Sprintf(":%d", n.info.Port), Handler: handler}
	// This is to ensure the server is closed if/when the context signals done. Most widely used in testing for early termination.
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	fmt.Printf("Listening for requests on port: %d.\n", n.info.Port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}
