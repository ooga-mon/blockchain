package node

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ooga-mon/blockchain/internal/database"
	"github.com/ooga-mon/blockchain/internal/wallet"
)

func TestNodeStart(t *testing.T) {
	node := NewNode(DefaultIP, DefaultPort, DefaultPeerIP, DefaultPeerPort)

	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)

	err := node.Start(ctx)
	if err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func TestMine(t *testing.T) {
	testUser1 := wallet.NewWallet()
	testUser2 := wallet.NewWallet()
	tx := database.NewTransaction(testUser1.PublicAddress, testUser2.PublicAddress, "testPayload", 1)
	signedTx, err := testUser1.SignTransaction(tx)
	if err != nil {
		t.Fatal(err)
	}

	node := NewNode(DefaultIP, DefaultPort, DefaultPeerIP, DefaultPeerPort)

	// stage adding a transaction before the blocking node.Start call
	go func() {
		time.Sleep(2 * time.Second)

		err := node.addPendingTransaction(signedTx)
		if err != nil {
			t.Log(err.Error())
		}
	}()

	// Timeout needs to be longer than the total duration of this test including mining difficulty
	ctx, stop := context.WithTimeout(context.Background(), 20*time.Second)

	go func() {
		// always add the same sleep for each go function to make setup semi-serialized.
		time.Sleep(2 * time.Second)
		// Periodically check if we mined a block.
		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-ticker.C:
				if len(node.db.Blocks) == 2 {
					stop()
					return
				}
			}
		}
	}()

	err = node.Start(ctx)
	if err != http.ErrServerClosed {
		t.Fatal(err)
	}

	if len(node.db.Blocks) != 2 {
		t.Errorf("Block was not added. Input: %d, Expected: %d", len(node.db.Blocks), 2)
	}
}
