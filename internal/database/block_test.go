package database

import (
	"testing"

	"encoding/hex"
	"time"
)

func TestNewBlock(t *testing.T) {
	blockParentHash := [32]byte{}
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	const blockHash = "c5c86affbc4f0a6e74a9446d2839b8c87765bee5143cc8b6c56da5ef7f6a61cf"
	const blockPayload = "test1"
	const nonce = 0
	const number = 0
	payload := Transactions{[]string{blockPayload}}
	block := NewBlock(blockParentHash, time, number, nonce, payload)

	hash := block.BlockHash
	if hex.EncodeToString(hash[:]) != blockHash {
		t.Errorf("block hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), blockHash)
	}
	if block.Content.ParentHash != blockParentHash {
		t.Errorf("block parentHash is improperly set. Should be empty")
	}
	if block.Content.Tx.Data[0] != blockPayload {
		t.Errorf("block payload is improperly set. Input: %s, Expected: %s.", block.Content.Tx.Data[0], blockPayload)
	}
	if block.Content.Number != number {
		t.Errorf("block number is improperly set. Input: %d, Expected: %d.", block.Content.Number, number)
	}
	if block.Content.Nonce != nonce {
		t.Errorf("block nonce is improperly set. Input: %d, Expected: %d.", block.Content.Nonce, nonce)
	}
}
