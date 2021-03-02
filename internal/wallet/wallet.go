package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ooga-mon/blockchain/internal/database"
)

type Wallet struct {
	Balance       int64
	PublicAddress common.Address
	PrivateKey    *ecdsa.PrivateKey
}

func NewWallet() Wallet {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := elliptic.Marshal(crypto.S256(), publicKey.X, publicKey.Y)
	publicKeyHash := crypto.Keccak256(publicKeyBytes[1:])

	address := common.BytesToAddress(publicKeyHash[12:])

	return Wallet{0, address, privateKey}
}

func (w *Wallet) SignTransaction(tx database.Transaction) (database.SignedTransaction, error) {
	encodedTransaction, err := tx.Encode()
	if err != nil {
		return database.SignedTransaction{}, err
	}

	sig, err := signMessage(encodedTransaction, w.PrivateKey)
	if err != nil {
		return database.SignedTransaction{}, err
	}

	return database.NewSignedTransaction(tx, sig), nil
}

func signMessage(msg []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	msgHash := sha256.Sum256(msg)

	return crypto.Sign(msgHash[:], privateKey)
}
