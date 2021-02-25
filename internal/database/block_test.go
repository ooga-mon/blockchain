package database

import (
	"testing"

	"encoding/hex"
	"time"
)

func TestNewBlock(t *testing.T) {
	blockParentHash := [32]byte{}
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	const blockHash = "19b1cf0231b28e49ba6228e646cb9ffd30166b0816a388448ec7095264a1e235"
	const blockPayload = "test1"
	payload := Transactions{[]string{blockPayload}}
	block := NewBlock(blockParentHash, time, 0, payload)

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
}
