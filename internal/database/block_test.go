package database

import (
	"testing"

	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestNewBlock(t *testing.T) {
	blockParentHash := [32]byte{}
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	const blockPayload = "test1"
	const nonce = 0
	const number = 0
	from := common.Address{}
	to := common.Address{}

	tx := NewTransaction(from, to, blockPayload, 1)
	payload := []Transaction{tx}
	block := NewBlock(blockParentHash, time, number, nonce, payload)

	if block.Content.ParentHash != blockParentHash {
		t.Errorf("block parentHash is improperly set. Should be empty")
	}
	if block.Content.Tx[0].Data != blockPayload {
		t.Errorf("block payload is improperly set. Input: %s, Expected: %s.", block.Content.Tx[0].Data, blockPayload)
	}
	if block.Content.Number != number {
		t.Errorf("block number is improperly set. Input: %d, Expected: %d.", block.Content.Number, number)
	}
	if block.Content.Nonce != nonce {
		t.Errorf("block nonce is improperly set. Input: %d, Expected: %d.", block.Content.Nonce, nonce)
	}
}
func TestGetGenesisBlockchain(t *testing.T) {
	const genesisHash = "477adc9a6172471ef8e0212570cb765cb8e4b3e90b4db352a7bf6a7cbd3a45d4"
	genesisParentHash := [32]byte{}
	genesis := loadGenesisBlock()
	if genesis.Content.ParentHash != genesisParentHash {
		t.Errorf("Genesis parentHash is improperly set. Should be empty.")
	}
	hash := genesis.BlockHash
	if hex.EncodeToString(hash[:]) != genesisHash {
		t.Errorf("genesis hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), genesisHash)
	}
	if len(genesis.Content.Tx) > 0 {
		t.Errorf("genesis payload should be empty")
	}
}
