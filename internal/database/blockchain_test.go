package database

import (
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	genesisEntity := LoadGenesisBlockEntity()

	if blockchain.blocks[0].key != genesisEntity.key {
		t.Error("New blockchain genesis hash does not equal original genesis hash")
	}
}

func TestAddBlock(t *testing.T) {
	blockchain := NewBlockchain()
	genesisEntity := LoadGenesisBlockEntity()
	const newBlockPayload = "payload1"

	payload := Payload{[]string{newBlockPayload}}
	blockchain.AddBlock(payload)

	if len(blockchain.blocks) != 2 {
		t.Fatalf("blockchain length is incorrect. Input: %d, Expected %d", len(blockchain.blocks), 2)
	}
	if blockchain.blocks[0].key != genesisEntity.key {
		t.Error("blockchain genesis hash does not equal original genesis hash after adding new block")
	}
	if blockchain.blocks[0].key != blockchain.blocks[1].value.blockHeader.parentHash {
		t.Error("new block parentHash does not equal previous blocks hash")
	}
}

func TestIsValid(t *testing.T) {
	blockchain := NewBlockchain()

	if !blockchain.IsValid() {
		t.Fatal("Created new blockchain and should be valid.")
	}

	payload := Payload{[]string{"payload"}}
	blockchain.AddBlock(payload)

	if !blockchain.IsValid() {
		t.Fatal("Added block to blockchain and should be valid.")
	}

	// tamper with blockchain
	blockchain.blocks = []BlockEntity{}
	if blockchain.IsValid() {
		t.Fatal("BlockEntities were all removed. Should be an invalid blockchain.")
	}

	blockchain = NewBlockchain()
	blockchain.AddBlock(payload)

	blockchain.blocks[0].value.blockHeader.number = 1
	if blockchain.IsValid() {
		t.Error("Geneisis block number was altered. Should be an invalid blockchain.")
	}

	blockchain.blocks[0].value.blockHeader.number = 0
	blockchain.blocks[1].value.blockHeader.number = 5
	if blockchain.IsValid() {
		t.Error("Block entity 2's number was altered. Should be an invalid blockchain.")
	}

	blockchain.blocks[1].value.blockHeader.number = 1
	blockchain.blocks[1].value.blockHeader.parentHash = [32]byte{}
	if blockchain.IsValid() {
		t.Error("Block entity 2's parentHash was altered. Should be an invalid blockchain.")
	}
}

func TestReplace(t *testing.T) {
	blockchain1 := NewBlockchain()
	blockchain2 := NewBlockchain()

	payload := Payload{[]string{"payload"}}
	blockchain2.AddBlock(payload)

	if !blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should have been replaced with blockchain2.")
	}
	if len(blockchain1.blocks) != len(blockchain2.blocks) {
		t.Errorf("Blockchain1 data should have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.blocks), len(blockchain2.blocks))
	}

	blockchain1.AddBlock(payload)
	if blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should not have been replaced with blockchain2.")
	}
	if len(blockchain1.blocks) == len(blockchain2.blocks) {
		t.Errorf("Blockchain1 data should not have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.blocks), len(blockchain2.blocks))
	}

}
