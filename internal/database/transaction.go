package database

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

type SignedTransaction struct {
	Transaction
	Sig []byte `json:"signature"`
}

func NewSignedTransaction(tx Transaction, sig []byte) SignedTransaction {
	return SignedTransaction{tx, sig}
}

func (t Transaction) Hash() (Hash, error) {
	jsonTx, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(jsonTx), nil
}

func (t SignedTransaction) Hash() (Hash, error) {
	jsonTx, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(jsonTx), nil
}

func (t Transaction) Encode() ([]byte, error) {
	return json.Marshal(t)
}

func (t SignedTransaction) IsAuthentic() (bool, error) {
	hashTx, err := t.Transaction.Hash()
	if err != nil {
		return false, err
	}

	recoveredPubKey, err := crypto.SigToPub(hashTx[:], t.Sig)
	if err != nil {
		return false, err
	}

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])

	return recoveredAccount.Hex() == t.From.Hex(), nil
}
