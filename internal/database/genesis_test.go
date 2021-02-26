package database

import (
	"testing"

	"encoding/hex"
)

func TestGetGenesisBlockchain(t *testing.T) {
	const genesisHash = "f540a12f198dce0e5407d1bc1e4d2e27d64964af152e58881b60031a0bc3aba3"
	genesisParentHash := [32]byte{}
	genesis := loadGenesisBlock()
	if genesis.Content.ParentHash != genesisParentHash {
		t.Errorf("Genesis parentHash is improperly set. Should be empty.")
	}
	hash := genesis.BlockHash
	if hex.EncodeToString(hash[:]) != genesisHash {
		t.Errorf("genesis hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), genesisHash)
	}
	if len(genesis.Content.Tx.Data) > 0 {
		t.Errorf("genesis payload should be empty")
	}
}
