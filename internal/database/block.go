package database

import (
	"crypto/sha256"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	blockHeader BlockHeader
	payload     Payload
}

type BlockHeader struct {
	parentHash [32]byte
	timestamp  time.Time
	number     uint64
}

type Payload struct {
	data []string
}

type BlockEntity struct {
	key   [32]byte
	value Block
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
	return b.blockHeader.toBytes() + b.payload.toBytes()
}

func (bh BlockHeader) toBytes() string {
	return string(bh.parentHash[:]) + time.Time.String(bh.timestamp) + strconv.FormatUint(bh.number, 10)
}

func (p Payload) toBytes() string {
	sb := strings.Builder{}
	for _, val := range p.data {
		sb.WriteString(val)
	}
	return sb.String()
}
