package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

// Wallets 保存多个钱包
type Wallets struct {
	Wallets map[string]*Wallet
}

type WalletsJson struct {
	Wallets map[string]*WalletJson
}

type WalletJson struct {
	PrivateKey PrivateKey
	PublicKey  []byte
}

type PrivateKey struct {
	//Curve   elliptic.Curve `json:"Curve"`
	X, Y, D *big.Int
}

// NewWallets 创建钱包并从文件中填充（如果存在）
func NewWallets(nodeID string) (*Wallets, error) {
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile(nodeID)

	return wallets, err
}

// CreateWallet 将钱包添加到钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets = make(map[string]*Wallet)
	ws.Wallets[address] = wallet

	return address
}

// GetAddresses 返回存储在钱包文件中的地址数组
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetWallet 通过地址返回钱包
func (ws *Wallets) GetWallet(address string) (Wallet, error) {
	wallet, ok := ws.Wallets[address]
	if ok {
		return *wallet, nil
	}

	return Wallet{}, errors.New("get wallet address nil")
}

// LoadFromFile 从文件中加载钱包
func (ws *Wallets) LoadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var wallets WalletsJson
	err = json.Unmarshal(fileContent, &wallets)
	if err != nil {
		return err
	}
	//gob.Register(elliptic.P256())
	//decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	//err = decoder.Decode(&wallets)
	//if err != nil {
	//	log.Panic(err)
	//}

	w := make(map[string]*Wallet)
	for k, v := range wallets.Wallets {
		w[k] = &Wallet{
			PublicKey: v.PublicKey,
			PrivateKey: ecdsa.PrivateKey{
				ecdsa.PublicKey{
					elliptic.P256(),
					v.PrivateKey.X,
					v.PrivateKey.Y,
				},
				v.PrivateKey.D,
			},
		}
	}

	ws.Wallets = w

	return nil
}

// SaveToFile 将钱包保存到文件
// 使用json序列化方式 方便查看与学习
func (ws *Wallets) SaveToFile() {
	//var content bytes.Buffer
	//gob.Register(elliptic.P256())
	//
	//encoder := gob.NewEncoder(&content)
	//err := encoder.Encode(ws)
	content, err := json.Marshal(ws)
	if err != nil {
		log.Panic(err)
	}
	//content.Bytes()

	err = ioutil.WriteFile(walletFile, content, 0644)
	if err != nil {
		log.Panic(err)
	}
}
