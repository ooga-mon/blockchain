package database

import (
	"time"
)

const DIFFICULTY = 4

type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

func NewBlockchain() Blockchain {
	return Blockchain{[]Block{loadGenesisBlock()}}
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) MineBlock(Tx Transactions) Block {
	lastBlock := bc.GetLastBlock()
	timestamp := time.Now()
	var nonce uint64 = 0
	isValidHash := false
	var newBlock Block
	for !isValidHash {
		nonce++
		newBlock = NewBlock(lastBlock.BlockHash, timestamp, lastBlock.Content.Number+1, nonce, Tx)
		isValidHash = isSuccessfulMinedHash(newBlock.BlockHash)
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	return newBlock
}

func isSuccessfulMinedHash(hash Hash) bool {
	hashHex := hash.Hex()
	for i := 0; i < DIFFICULTY; i++ {
		if hashHex[i] != '0' {
			return false
		}
	}

	return true
}

func (bc *Blockchain) IsValid() bool {
	// We assume that a new blockchain, at a minimum, contains the genesis block
	if len(bc.Blocks) == 0 {
		return false
	}
	genesis := loadGenesisBlock()
	if bc.Blocks[0].BlockHash != genesis.BlockHash || bc.Blocks[0].Content.Hash() != genesis.BlockHash {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		if bc.Blocks[i].Content.ParentHash != bc.Blocks[i-1].BlockHash {
			return false
		}
		if bc.Blocks[i].Content.Hash() != bc.Blocks[i].BlockHash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) Replace(newChain *Blockchain) bool {
	if len(newChain.Blocks) <= len(bc.Blocks) {
		return false
	}
	if !newChain.IsValid() {
		return false
	}

	bc.Blocks = newChain.Blocks
	return true
}
