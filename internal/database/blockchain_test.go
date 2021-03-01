package database

import (
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	blockchain := NewBlockchain()
	genesisEntity := getGenesisTestBlock()

	if blockchain.Blocks[0].BlockHash != genesisEntity.block.BlockHash {
		t.Error("New blockchain genesis hash does not equal original genesis hash")
	}
}

func TestAddBlock(t *testing.T) {
	blockchain := NewBlockchain()
	genesis := getGenesisTestBlock()
	newBlock := getRandomTestNextBlock(genesis)
	blockchain.AddBlock(newBlock.block)

	if len(blockchain.Blocks) != 2 {
		t.Fatalf("blockchain length is incorrect. Input: %d, Expected %d", len(blockchain.Blocks), 2)
	}
	if blockchain.Blocks[0].BlockHash != genesis.block.BlockHash {
		t.Errorf("blockchain genesis hash does not equal original genesis hash after adding new block")
	}
	if blockchain.Blocks[0].BlockHash != blockchain.Blocks[1].Content.ParentHash {
		t.Errorf("New block parentHash does not equal previous blocks hash. Input: %s, Expected: %s", blockchain.Blocks[0].BlockHash.Hex(), blockchain.Blocks[1].BlockHash.Hex())
	}
}

func TestIsValid(t *testing.T) {
	blockchain := NewBlockchain()
	genesis := getGenesisTestBlock()

	if !blockchain.IsValid() {
		t.Fatal("Created new blockchain and should be valid.")
	}

	newBlock := getRandomTestNextBlock(genesis)
	blockchain.AddBlock(newBlock.block)

	if !blockchain.IsValid() {
		t.Fatal("Added block to blockchain and should be valid.")
	}

	// tamper with blockchain
	blockchain.Blocks = []Block{}
	if blockchain.IsValid() {
		t.Fatal("BlockEntities were all removed. Should be an invalid blockchain.")
	}

	blockchain = NewBlockchain()
	blockchain.AddBlock(newBlock.block)

	blockchain.Blocks[0].Content.Number = 1
	if blockchain.IsValid() {
		t.Error("Geneisis block number was altered. Should be an invalid blockchain.")
	}

	blockchain = NewBlockchain()
	blockchain.AddBlock(newBlock.block)
	blockchain.Blocks[1].Content.Number = 5
	if blockchain.IsValid() {
		t.Error("Block entity 2's number was altered. Should be an invalid blockchain.")
	}

	blockchain = NewBlockchain()
	blockchain.AddBlock(newBlock.block)
	blockchain.Blocks[1].Content.ParentHash = [32]byte{}
	if blockchain.IsValid() {
		t.Error("Block entity 2's parentHash was altered. Should be an invalid blockchain.")
	}
}

func TestReplace(t *testing.T) {
	blockchain1 := NewBlockchain()
	blockchain2 := NewBlockchain()
	genesis := getGenesisTestBlock()

	newBlock := getRandomTestNextBlock(genesis)
	blockchain2.AddBlock(newBlock.block)

	if !blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should have been replaced with blockchain2.")
	}
	if len(blockchain1.Blocks) != len(blockchain2.Blocks) {
		t.Errorf("Blockchain1 data should have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.Blocks), len(blockchain2.Blocks))
	}

	newBlock = getRandomTestNextBlock(newBlock)
	blockchain1.AddBlock(newBlock.block)
	if blockchain1.Replace(&blockchain2) {
		t.Error("Blockchain1 should not have been replaced with blockchain2.")
	}
	if len(blockchain1.Blocks) == len(blockchain2.Blocks) {
		t.Errorf("Blockchain1 data should not have been replaced with blockchain2 data. Input: %d, Expected: %d", len(blockchain1.Blocks), len(blockchain2.Blocks))
	}

}
