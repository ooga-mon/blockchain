package database

import (
	"testing"

	"encoding/hex"
	"time"
)

func TestNewBlock(t *testing.T) {
	blockParentHash := [32]byte{}
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	const blockHash = "9ff955fe117ddc51a11409ba6a82b25c21f4d6c23c8a90dc98b240f5284d1d71"
	const blockPayload = "test1"
	payload := Payload{[]string{blockPayload}}
	block := NewBlock(blockParentHash, time, 0, payload)

	hash := block.Hash()
	if hex.EncodeToString(hash[:]) != blockHash {
		t.Errorf("block hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), blockHash)
	}
	if block.BlockHeader.ParentHash != blockParentHash {
		t.Errorf("block parentHash is improperly set. Should be empty")
	}
	if block.Payload.Data[0] != blockPayload {
		t.Errorf("block payload is improperly set. Input: %s, Expected: %s.", block.Payload.Data[0], blockPayload)
	}
}
