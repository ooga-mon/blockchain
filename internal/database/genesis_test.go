package database

import (
	"testing"

	"encoding/hex"
)

func TestGetGenesisBlockchain(t *testing.T) {
	const genesisHash = "c812e13be23ef807807739cad4ac5cf28c6535c56b8c1397267b1d825f890501"
	genesisParentHash := [32]byte{}
	genesis := LoadGenesisBlockEntity()
	if genesis.Value.BlockHeader.ParentHash != genesisParentHash {
		t.Errorf("Genesis parentHash is improperly set. Should be empty.")
	}
	hash := genesis.Value.Hash()
	if hex.EncodeToString(hash[:]) != genesisHash {
		t.Errorf("genesis hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), genesisHash)
	}
	if genesis.Value.Hash() != genesis.Key {
		t.Errorf("genesis block hash and key of genesis block do not match")
	}
	if len(genesis.Value.Payload.Data) > 0 {
		t.Errorf("genesis payload should be empty")
	}
}
