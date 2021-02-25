package node

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (n *Node) sync(ctx context.Context) error {
	n.doSync(ctx)

	ticker := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-ticker.C:
			n.doSync(ctx)
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) doSync(ctx context.Context) {
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

		status, err := getPeerStatus(peer)
		if err != nil {
			n.removePeer(peer)
		}

		n.addUnknownPeers(status)

		n.pushKnownPeers(peer)
	}
}

func getPeerStatus(peer connectionInfo) (Status, error) {
	url := fmt.Sprintf("http://%s/node/status", peer.tcpAddress())
	res, err := http.Get(url)
	if err != nil {
		return Status{}, err
	}

	peerStatus := Status{}
	err = readResponse(res, &peerStatus)
	if err != nil {
		return Status{}, err
	}

	return peerStatus, nil
}

func (n *Node) addUnknownPeers(peerStatus Status) error {
	for _, unknownPeer := range peerStatus.KnowPeers {
		if !n.containsPeer(unknownPeer) {
			fmt.Printf("Found new peer %s\n", unknownPeer.tcpAddress())
			n.addPeer(unknownPeer)
		}
	}

	return nil
}

func (n *Node) pushKnownPeers(peer connectionInfo) error {
	if peer.connected {
		return nil
	}

	url := fmt.Sprintf("http://%s/node/peer?ip=%s&port=%d", peer.tcpAddress(), n.info.IP, n.info.Port)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	stdRes := StandardResponse{}
	err = readResponse(res, stdRes)
	if err != nil {
		return err
	}
	if !stdRes.Success {
		return fmt.Errorf(stdRes.Error)
	}

	knownPeer := n.peers[peer.tcpAddress()]
	knownPeer.connected = stdRes.Success

	// reassigns known peer back to map
	n.addPeer(knownPeer)

	return nil
}
