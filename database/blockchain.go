package database

type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

func NewBlockchain() Blockchain {
	return Blockchain{[]Block{loadGenesisBlock()}}
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) AddBlock(newBlock Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) IsValid() bool {
	// We assume that a new blockchain, at a minimum, contains the genesis block
	if len(bc.Blocks) == 0 {
		return false
	}
	genesis := loadGenesisBlock()
	if bc.Blocks[0].BlockHash != genesis.BlockHash || bc.Blocks[0].Content.Hash() != genesis.BlockHash {
		return false
	}

	for i := 1; i < len(bc.Blocks); i++ {
		if bc.Blocks[i].Content.ParentHash != bc.Blocks[i-1].BlockHash {
			return false
		}
		if bc.Blocks[i].Content.Hash() != bc.Blocks[i].BlockHash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) Replace(newChain *Blockchain) bool {
	if len(newChain.Blocks) <= len(bc.Blocks) {
		return false
	}
	if !newChain.IsValid() {
		return false
	}

	bc.Blocks = newChain.Blocks
	return true
}
