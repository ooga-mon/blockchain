package database

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Block struct {
	BlockHash Hash
	Content   Content
}

type Content struct {
	ParentHash Hash
	Timestamp  time.Time
	Number     uint64
	Tx         Transactions
}

type Transactions struct {
	Data []string
}

func NewBlock(parentHash Hash, time time.Time, number uint64, payload Transactions) Block {
	content := Content{parentHash, time, number, payload}
	return Block{content.Hash(), content}
}

func (c Content) Hash() Hash {
	jsonBlock, err := json.Marshal(c)
	if err != nil {
		return Hash{}
	}

	return sha256.Sum256(jsonBlock)
}
