package database

import (
	"crypto/sha256"
	"encoding/json"
	"time"
)

type Block struct {
	BlockHash Hash    `json:"Hash"`
	Content   Content `json:"Content"`
}

type Content struct {
	ParentHash Hash                `json:"Parent_Hash"`
	Timestamp  time.Time           `json:"Timestamp"`
	Number     uint64              `json:"Number"`
	Nonce      uint64              `json:"Nonce"`
	Tx         []SignedTransaction `json:"Transactions"`
}

func NewBlock(parentHash Hash, time time.Time, number uint64, nonce uint64, payload []SignedTransaction) Block {
	content := Content{parentHash, time, number, nonce, payload}
	return Block{content.Hash(), content}
}

func (c Content) Hash() Hash {
	jsonBlock, err := json.Marshal(c)
	if err != nil {
		return Hash{}
	}

	return sha256.Sum256(jsonBlock)
}

func loadGenesisBlock() Block {
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	payload := []SignedTransaction{}
	return NewBlock([32]byte{}, time, 0, 0, payload)
}
