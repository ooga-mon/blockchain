package wallet

import (
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ooga-mon/blockchain/internal/database"
)

func TestSign(t *testing.T) {
	wal := NewWallet()

	msg := []byte("simple payload")

	sig, err := signMessage(msg, wal.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	recoveredAccount, err := recoverPublicAccount(msg, sig)
	if err != nil {
		t.Fatal(err)
	}

	if wal.PublicAddress.Hex() != recoveredAccount.Hex() {
		t.Fatalf("Message signature does not matched from address. Input: %s, Expected: %s", wal.PublicAddress.Hex(), recoveredAccount.Hex())
	}
}

func TestSignTransaction(t *testing.T) {
	user1 := NewWallet()
	user2 := NewWallet()

	tx := database.NewTransaction(user1.PublicAddress, user2.PublicAddress, "simple payload", 1)
	signedTx, err := user1.signTransaction(tx)
	if err != nil {
		t.Error(err)
		return
	}

	verifySignedTransaction(t, signedTx, user1)
}

func TestCreateTransaction(t *testing.T) {
	user1 := NewWallet()
	user2 := NewWallet()

	signedTx, err := user1.CreateTransaction(user2.PublicAddress, "simple payload")
	if err != nil {
		t.Error(err)
		return
	}

	verifySignedTransaction(t, signedTx, user1)
}

func verifySignedTransaction(t *testing.T, signedTx database.SignedTransaction, user Wallet) {
	encodedTx, err := signedTx.Transaction.Encode()
	if err != nil {
		t.Fatal(err)
	}
	recoveredAccount, err := recoverPublicAccount(encodedTx, signedTx.Sig)
	if err != nil {
		t.Fatal(err)
	}

	if recoveredAccount.Hex() != signedTx.From.Hex() {
		t.Fatalf("Message signature does not matched from address. Input: %s, Expected: %s", signedTx.From.Hex(), recoveredAccount.Hex())
	}

	if recoveredAccount.Hex() != user.PublicAddress.Hex() {
		t.Fatalf("recovered account does not matched wallet public address. Input: %s, Expected: %s", recoveredAccount.Hex(), user.PublicAddress.Hex())
	}
}

func recoverPublicAccount(msg, sig []byte) (common.Address, error) {
	msgHash := sha256.Sum256(msg)

	recoveredPubKey, err := crypto.SigToPub(msgHash[:], sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("unable to verify message signature. %s", err.Error())
	}

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])

	return recoveredAccount, nil
}
