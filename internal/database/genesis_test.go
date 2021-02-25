package database

import (
	"testing"

	"encoding/hex"
)

func TestGetGenesisBlockchain(t *testing.T) {
	const genesisHash = "457ba46d82bdfe88a7f7b0c52001122674cc3dde79ab7a94f2b8b40bf80a3b2c"
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
