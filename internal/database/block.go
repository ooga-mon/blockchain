package database

import (
	"crypto/sha256"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	BlockHeader BlockHeader
	Payload     Payload
}

type BlockHeader struct {
	ParentHash [32]byte
	Timestamp  time.Time
	Number     uint64
}

type Payload struct {
	Data []string
}

type BlockEntity struct {
	Key   [32]byte
	Value Block
}

func NewBlock(parentHash [32]byte, time time.Time, number uint64, payload Payload) Block {
	header := BlockHeader{parentHash, time, number}
	return Block{header, payload}
}

func NewBlockEntity(block Block) BlockEntity {
	return BlockEntity{block.Hash(), block}
}

func (b Block) Hash() [32]byte {
	return sha256.Sum256([]byte(b.contentAsBytes()))
}

func (b Block) contentAsBytes() string {
	return b.BlockHeader.toBytes() + b.Payload.toBytes()
}

func (bh BlockHeader) toBytes() string {
	return string(bh.ParentHash[:]) + time.Time.String(bh.Timestamp) + strconv.FormatUint(bh.Number, 10)
}

func (p Payload) toBytes() string {
	sb := strings.Builder{}
	for _, val := range p.Data {
		sb.WriteString(val)
	}
	return sb.String()
}
