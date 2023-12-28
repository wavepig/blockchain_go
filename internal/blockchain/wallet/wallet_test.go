package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestNewKeyPair(t *testing.T) {
	private, public := newKeyPair()
	indent, _ := json.MarshalIndent(private, "", " ")
	t.Log("private: ", string(indent))
	t.Log("public: ", public)
}

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	t.Log("address: ", string(address))

	ok := ValidateAddress(string(address))
	if !ok {
		t.Error("检验失败")
	}
}

func TestWalletsSaveToFile(t *testing.T) {
	wallets := &Wallets{}
	wallet := wallets.CreateWallet()
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallet)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func TestWalletsCreateWallet(t *testing.T) {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}
