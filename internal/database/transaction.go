package database

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	From    common.Address `json:"From"`
	To      common.Address `json:"To"`
	Time    time.Time      `json:"Timestamp"`
	Data    string         `json:"Data"`
	Payment int            `json:"Payment"`
}

func NewTransaction(from, to common.Address, data string, payment int) Transaction {
	return Transaction{from, to, time.Now(), data, payment}
}
