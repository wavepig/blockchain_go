package wallet

import (
	"encoding/json"
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
}
