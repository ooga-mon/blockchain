package database

import (
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	genesisEntity := loadGenesisBlock()

	if blockchain.Blocks[0].BlockHash != genesisEntity.BlockHash {
		t.Error("New blockchain genesis hash does not equal original genesis hash")
	}
}

func TestMineBlock(t *testing.T) {
	blockchain := NewBlockchain()
	genesisEntity := loadGenesisBlock()
	const newBlockPayload = "payload1"

	payload := Transactions{[]string{newBlockPayload}}
	blockchain.MineBlock(payload)

	if len(blockchain.Blocks) != 2 {
		t.Fatalf("blockchain length is incorrect. Input: %d, Expected %d", len(blockchain.Blocks), 2)
	}
	if blockchain.Blocks[0].BlockHash != genesisEntity.BlockHash {
		t.Error("blockchain genesis hash does not equal original genesis hash after adding new block")
	}
	if blockchain.Blocks[0].BlockHash != blockchain.Blocks[1].Content.ParentHash {
		t.Error("new block parentHash does not equal previous blocks hash")
	}
	hashHex := blockchain.Blocks[1].BlockHash.Hex()
	for i := 0; i < DIFFICULTY; i++ {
		if hashHex[i] != '0' {
			t.Errorf("new block hash does not meet difficulty requirements. Input: %s", blockchain.Blocks[1].BlockHash.Hex())
		}
	}
}

func TestIsValid(t *testing.T) {
	blockchain := NewBlockchain()

	if !blockchain.IsValid() {
		t.Fatal("Created new blockchain and should be valid.")
	}

	payload := Transactions{[]string{"payload"}}
	blockchain.MineBlock(payload)

	if !blockchain.IsValid() {
		t.Fatal("Added block to blockchain and should be valid.")
	}

	// tamper with blockchain
	blockchain.Blocks = []Block{}
	if blockchain.IsValid() {
		t.Fatal("BlockEntities were all removed. Should be an invalid blockchain.")
	}

	blockchain = NewBlockchain()
	blockchain.MineBlock(payload)

	blockchain.Blocks[0].Content.Number = 1
	if blockchain.IsValid() {
		t.Error("Geneisis block number was altered. Should be an invalid blockchain.")
	}

	blockchain.Blocks[0].Content.Number = 0
	blockchain.Blocks[1].Content.Number = 5
	if blockchain.IsValid() {
		t.Error("Block entity 2's number was altered. Should be an invalid blockchain.")
	}

	blockchain.Blocks[1].Content.Number = 1
	blockchain.Blocks[1].Content.ParentHash = [32]byte{}
	if blockchain.IsValid() {
		t.Error("Block entity 2's parentHash was altered. Should be an invalid blockchain.")
	}
}

func TestReplace(t *testing.T) {
	blockchain1 := NewBlockchain()
	blockchain2 := NewBlockchain()

	payload := Transactions{[]string{"payload"}}
	blockchain2.MineBlock(payload)

	if !blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should have been replaced with blockchain2.")
	}
	if len(blockchain1.Blocks) != len(blockchain2.Blocks) {
		t.Errorf("Blockchain1 data should have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.Blocks), len(blockchain2.Blocks))
	}

	blockchain1.MineBlock(payload)
	if blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should not have been replaced with blockchain2.")
	}
	if len(blockchain1.Blocks) == len(blockchain2.Blocks) {
		t.Errorf("Blockchain1 data should not have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.Blocks), len(blockchain2.Blocks))
	}

}
