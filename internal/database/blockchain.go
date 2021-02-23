package database

import (
	"time"
)

type Blockchain struct {
	blocks []BlockEntity
}

func NewBlockchain() Blockchain {
	return Blockchain{[]BlockEntity{LoadGenesisBlockEntity()}}
}

func (bc *Blockchain) AddBlock(payload Payload) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(lastBlock.key, time.Now(), lastBlock.value.blockHeader.number+1, payload)
	bc.blocks = append(bc.blocks, NewBlockEntity(newBlock))
}

func (bc *Blockchain) IsValid() bool {
	// We assume that a new blockchain, at a minimum, contains the genesis block
	if len(bc.blocks) == 0 {
		return false
	}
	genesis := LoadGenesisBlockEntity()
	if bc.blocks[0].key != genesis.key || bc.blocks[0].value.Hash() != genesis.key {
		return false
	}

	for i := 1; i < len(bc.blocks); i++ {
		if bc.blocks[i].value.blockHeader.parentHash != bc.blocks[i-1].key {
			return false
		}
		if bc.blocks[i].value.Hash() != bc.blocks[i].key {
			return false
		}
	}

	return true
}

func (bc *Blockchain) Replace(newChain *Blockchain) bool {
	if len(newChain.blocks) <= len(bc.blocks) {
		return false
	}
	if !newChain.IsValid() {
		return false
	}

	bc.blocks = newChain.blocks
	return true
}
