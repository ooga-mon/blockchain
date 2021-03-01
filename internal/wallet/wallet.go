package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ooga-mon/blockchain/internal/database"
)

type Wallet struct {
	balance    int64
	privateKey *ecdsa.PrivateKey
}

func NewWallet() Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	return Wallet{0, privateKey}
}

func (w *Wallet) SignTransaction(tx database.Transaction, privateKey *ecdsa.PrivateKey) (database.SignedTransaction, error) {
	encodedTransaction, err := tx.Encode()
	if err != nil {
		return database.SignedTransaction{}, err
	}

	sig, err := sign(encodedTransaction, privateKey)
	if err != nil {
		return database.SignedTransaction{}, err
	}

	return database.NewSignedTransaction(tx, sig), nil
}

func sign(msg []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	msgHash := sha256.Sum256(msg)

	return crypto.Sign(msgHash[:], privateKey)
}
