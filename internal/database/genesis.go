package database

import (
	"time"
)

func LoadGenesisBlockEntity() BlockEntity {
	genesisBlock := loadGenesisBlock()
	return BlockEntity{genesisBlock.Hash(), genesisBlock}
}

func loadGenesisBlock() Block {
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	payload := Payload{}
	return NewBlock([32]byte{}, time, 0, payload)
}
