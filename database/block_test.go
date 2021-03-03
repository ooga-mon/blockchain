package database

import (
	"testing"

	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestNewBlock(t *testing.T) {
	testBlock := getDefaultTestBlock()

	if testBlock.block.Content.ParentHash != testBlock.blockParentHash {
		t.Errorf("block parentHash is improperly set. Should be empty")
	}
	if testBlock.block.Content.Tx[0].Data != testBlock.payloadData {
		t.Errorf("block payload is improperly set. Input: %s, Expected: %s.", testBlock.block.Content.Tx[0].Data, testBlock.payloadData)
	}
	if testBlock.block.Content.Number != testBlock.number {
		t.Errorf("block number is improperly set. Input: %d, Expected: %d.", testBlock.block.Content.Number, testBlock.number)
	}
	if testBlock.block.Content.Nonce != testBlock.nonce {
		t.Errorf("block nonce is improperly set. Input: %d, Expected: %d.", testBlock.block.Content.Nonce, testBlock.nonce)
	}

	testBlockNext := getRandomTestNextBlock(testBlock)
	if testBlockNext.block.Content.ParentHash != testBlockNext.blockParentHash {
		t.Errorf("block parentHash is improperly set. Should be empty")
	}
	if testBlockNext.block.Content.Tx[0].Data != testBlockNext.payloadData {
		t.Errorf("block payload is improperly set. Input: %s, Expected: %s.", testBlockNext.block.Content.Tx[0].Data, testBlockNext.payloadData)
	}
	if testBlockNext.block.Content.Number != testBlockNext.number {
		t.Errorf("block number is improperly set. Input: %d, Expected: %d.", testBlockNext.block.Content.Number, testBlockNext.number)
	}
	if testBlockNext.block.Content.Nonce != testBlockNext.nonce {
		t.Errorf("block nonce is improperly set. Input: %d, Expected: %d.", testBlockNext.block.Content.Nonce, testBlockNext.nonce)
	}
}

func TestGetGenesisBlockchain(t *testing.T) {
	const genesisHash = "477adc9a6172471ef8e0212570cb765cb8e4b3e90b4db352a7bf6a7cbd3a45d4"
	genesis := loadGenesisBlock()
	genesisTest := getGenesisTestBlock()
	if genesis.Content.ParentHash != genesisTest.blockParentHash {
		t.Errorf("Genesis parentHash is improperly set. Should be empty.")
	}
	hash := genesis.BlockHash
	if hex.EncodeToString(hash[:]) != genesisHash {
		t.Errorf("genesis hash is improperly set. Input: %s, Expected: %s.", hex.EncodeToString(hash[:]), genesisHash)
	}
	if len(genesis.Content.Tx) > 0 {
		t.Errorf("genesis payload should be empty. Input: %d, Expected: %d", len(genesis.Content.Tx), 0)
	}
	if genesis.Content.Number != genesisTest.number {
		t.Errorf("block number is improperly set. Input: %d, Expected: %d.", genesis.Content.Number, genesisTest.number)
	}
	if genesis.Content.Nonce != genesisTest.nonce {
		t.Errorf("block nonce is improperly set. Input: %d, Expected: %d.", genesis.Content.Nonce, genesisTest.nonce)
	}
}

type testBlock struct {
	blockParentHash [32]byte
	time            time.Time
	nonce           uint64
	number          uint64
	from            common.Address
	to              common.Address
	tx              Transaction
	signTx          SignedTransaction
	payloadData     string
	payload         []SignedTransaction
	block           Block
}

func (tb *testBlock) generateNewBlock() {
	tx := NewTransaction(tb.from, tb.to, tb.payloadData, 1)
	signedTx := SignedTransaction{tx, []byte{}}
	payload := []SignedTransaction{signedTx}
	tb.block = NewBlock(tb.blockParentHash, tb.time, tb.number, tb.nonce, payload)
}

func getTestBlock(parentHash [32]byte, timestamp time.Time, nonce uint64, blockPayload string, number uint64, from common.Address, to common.Address) testBlock {
	tx := NewTransaction(from, to, blockPayload, 1)
	signedTx := SignedTransaction{tx, []byte{}}
	payload := []SignedTransaction{signedTx}
	block := NewBlock(parentHash, timestamp, number, nonce, payload)
	return testBlock{
		blockParentHash: parentHash,
		time:            timestamp,
		nonce:           nonce,
		number:          number,
		from:            from,
		to:              to,
		tx:              tx,
		signTx:          signedTx,
		payloadData:     blockPayload,
		payload:         payload,
		block:           block,
	}
}

func getGenesisTestBlock() testBlock {
	return testBlock{
		blockParentHash: [32]byte{},
		time:            time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC),
		nonce:           0,
		number:          0,
		from:            common.Address{},
		to:              common.Address{},
		tx:              Transaction{},
		signTx:          SignedTransaction{},
		payloadData:     "",
		payload:         []SignedTransaction{},
		block:           loadGenesisBlock(),
	}
}

func getDefaultTestBlock() testBlock {
	base := getGenesisTestBlock()
	base.payloadData = "test1"
	base.time = time.Now()
	base.generateNewBlock()
	return base
}

func getRandomTestNextBlock(parentTestBlock testBlock) testBlock {
	timestamp := time.Now()
	var nonce uint64 = 10
	const blockPayload = "randomPayload"
	number := parentTestBlock.number + 1
	from := common.Address{}
	to := common.Address{}
	return getTestBlock(parentTestBlock.block.BlockHash, timestamp, nonce, blockPayload, number, from, to)
}
