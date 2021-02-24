package database

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Block struct {
	BlockHeader BlockHeader
	Payload     Payload
}

type BlockHeader struct {
	ParentHash Hash
	Timestamp  time.Time
	Number     uint64
}

type Payload struct {
	Data []string
}

type BlockEntity struct {
	Key   Hash
	Value Block
}

func NewBlock(parentHash Hash, time time.Time, number uint64, payload Payload) Block {
	header := BlockHeader{parentHash, time, number}
	return Block{header, payload}
}

func NewBlockEntity(block Block) BlockEntity {
	return BlockEntity{block.Hash(), block}
}

func (b Block) Hash() Hash {
	jsonBlock, err := json.Marshal(b)
	if err != nil {
		return Hash{}
	}

	return sha256.Sum256(jsonBlock)
}
