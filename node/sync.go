package node

import (
	"fmt"

	"github.com/ooga-mon/blockchain/database"
)

func (n *Node) broadcastTransaction(tx database.SignedTransaction) {
	go n.doBroadcastTransaction(tx)
}

func (n *Node) doBroadcastTransaction(tx database.SignedTransaction) {
	for _, peer := range n.peers {
		err := n.postTransaction(peer, tx)
		if err != nil {
			fmt.Print("Removing peer from list.")
			n.removePeer(peer)
		}
	}
}

func (n *Node) postTransaction(peer connectionInfo, tx database.SignedTransaction) error {
	url := fmt.Sprintf("http://%s/node/sync/transaction", peer.tcpAddress())
	res, err := WriteRequest(url, tx)
	if err != nil {
		return err
	}

	res.Body.Close()
	return nil
}

func (n *Node) broadcastStatus() {
	go n.doBroadcastStatus()
}

func (n *Node) doBroadcastStatus() {
	fmt.Println("Attempting sync process")
	for _, peer := range n.peers {
		fmt.Printf("Starting sync process with %s\n", peer.tcpAddress())
		// We can be a little defensive here even though there are checks in other places.
		// Just in case a peer add accidentally maps itself. This could cause unwanted behaviour in this function.
		if n.info.IP == peer.IP && n.info.Port == peer.Port {
			continue
		}
		if peer.IP == "" || peer.Port == 0 {
			continue
		}

		peerStatus, err := n.postStatus(peer)
		if err != nil {
			fmt.Print("Removing peer from list.")
			n.removePeer(peer)
			continue
		}

		n.updateStatusDifferences(peerStatus)
	}
}

func (n *Node) postStatus(peer connectionInfo) (Status, error) {
	url := fmt.Sprintf("http://%s/node/sync/status", peer.tcpAddress())
	nodeStatus := Status{n.info, n.peers, n.db}
	res, err := WriteRequest(url, nodeStatus)
	if err != nil {
		return Status{}, err
	}

	peerStatus := Status{}
	err = ReadResponse(res, &peerStatus)
	if err != nil {
		return Status{}, err
	}

	return peerStatus, nil
}

func (n *Node) updateStatusDifferences(peerStatus Status) {
	n.addUnknownPeers(peerStatus)
	n.db.Replace(&peerStatus.Blockchain)
}

func (n *Node) addUnknownPeers(peerStatus Status) error {
	n.addPeer(peerStatus.Info)
	for _, unknownPeer := range peerStatus.KnowPeers {
		n.addPeer(unknownPeer)
	}

	return nil
}
