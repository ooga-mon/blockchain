package database

import (
	"time"
)

func loadGenesisBlock() Block {
	time := time.Date(2021, time.February, 13, 20, 0, 0, 0, time.UTC)
	payload := Transactions{}
	return NewBlock([32]byte{}, time, 0, payload)
}
