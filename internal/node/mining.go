package node

import (
	"time"

	"github.com/ooga-mon/blockchain/internal/database"
)

const STARTING_DIFFICULTY = 4
const MINING_RATE = 5 * time.Second
const MAX_DIFFICULTY = 6
const MIN_DIFFICULTY = 2

type state struct {
	curDifficulty     int
	lastMineTimestamp time.Time
}

func newState() state {
	return state{STARTING_DIFFICULTY, time.Now()}
}

func (s *state) adjustDifficulty(time time.Time) {
	if time.Sub(s.lastMineTimestamp) < MINING_RATE {
		if s.curDifficulty != MAX_DIFFICULTY {
			s.curDifficulty++
		}
	} else if s.curDifficulty > MIN_DIFFICULTY {
		s.curDifficulty--
	}
}

func (s *state) setLastMineTime() {
	s.lastMineTimestamp = time.Now()
}

type pendingBlock struct {
	parentHash database.Hash
	number     uint64
	tx         database.Transactions
}

func newPendingBlock(parentHash database.Hash, number uint64, tx database.Transactions) pendingBlock {
	return pendingBlock{parentHash, number, tx}
}

func (n *Node) mineBlock(block pendingBlock) database.Block {
	timestamp := time.Now()
	var nonce uint64 = 0
	var newBlock database.Block
	n.state.adjustDifficulty(timestamp)
	for !isSuccessfulMinedHash(newBlock.BlockHash, n.state.curDifficulty) {
		nonce++
		newBlock = database.NewBlock(block.parentHash, timestamp, block.number, nonce, block.tx)
	}

	return newBlock
}

func isSuccessfulMinedHash(hash database.Hash, difficulty int) bool {
	if hash.IsEmpty() {
		return false
	}
	hashHex := hash.Hex()
	for i := 0; i < difficulty; i++ {
		if hashHex[i] != '0' {
			return false
		}
	}

	return true
}
