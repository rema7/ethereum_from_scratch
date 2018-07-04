package main

import (
	"crypto/ecdsa"
	"crypto/rand"

	"bytes"
	"encoding/gob"
	"io/ioutil"

	"os"

	"github.com/inconshreveable/log15"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

var log = log15.New("module", "wallet")

const walletFile = "wallet.dat"

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func InitWallet() *Wallet {
	wallet := &Wallet{}
	log.Info("init wallet")

	if err := wallet.Load(); err != nil {
		log.Info("creating new wallet")

		priv, pub := newKeyPair()

		wallet.PrivateKey = priv
		wallet.PublicKey = pub

		wallet.Save()
	}

	return wallet
}

func (w *Wallet) Load() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var wallet Wallet
	gob.Register(secp256k1.S256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallet)
	if err != nil {
		return err
	}

	w.PrivateKey = wallet.PrivateKey
	w.PublicKey = wallet.PublicKey

	return nil
}

func (w *Wallet) Save() error {
	var content bytes.Buffer

	gob.Register(secp256k1.S256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Error("can't encode wallet")
		return err
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Error("can't save to file", "file", walletFile)
		return err
	}
	return nil
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := secp256k1.S256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Error("save wallet", "can't save to file", walletFile)
	}
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, pubKey
}

func Keccak256(payload []byte) []byte {

	var result []byte
	hash := sha3.NewKeccak256()
	hash.Write(payload)
	result = hash.Sum(result)
	return result
}
