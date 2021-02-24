package database

import (
	"time"
)

type Blockchain struct {
	Blocks []BlockEntity
}

func NewBlockchain() Blockchain {
	return Blockchain{[]BlockEntity{LoadGenesisBlockEntity()}}
}

func (bc *Blockchain) AddBlock(payload Payload) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(lastBlock.Key, time.Now(), lastBlock.Value.BlockHeader.Number+1, payload)
	bc.Blocks = append(bc.Blocks, NewBlockEntity(newBlock))
}

func (bc *Blockchain) IsValid() bool {
	// We assume that a new blockchain, at a minimum, contains the genesis block
	if len(bc.Blocks) == 0 {
		return false
	}
	genesis := LoadGenesisBlockEntity()
	if bc.Blocks[0].Key != genesis.Key || bc.Blocks[0].Value.Hash() != genesis.Key {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		if bc.Blocks[i].Value.BlockHeader.ParentHash != bc.Blocks[i-1].Key {
			return false
		}
		if bc.Blocks[i].Value.Hash() != bc.Blocks[i].Key {
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
