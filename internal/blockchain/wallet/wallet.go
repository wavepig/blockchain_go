package wallet

import (
	"blockchain_go/pkg/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const addressChecksumLen = 4

// Wallet 存储私钥和公钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet 创建并返回钱包
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// GetAddress 将一个公钥转换成一个 Base58 地址
func (w *Wallet) GetAddress() []byte {
	// 使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	pubKeyHash := HashPubKey(w.PublicKey)
	// 给哈希加上地址生成算法版本的前缀
	versionedPayload := append([]byte{version}, pubKeyHash...)
	// 对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节。
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	// 使用 Base58 对 version+PubKeyHash+checksum 组合进行编码
	address := utils.Base58Encode(fullPayload)

	return address
}

// HashPubKey 哈希公钥 key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	// 实现了RIPEMD-160哈希算法
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress 检查地址是否有效
func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum 生成公钥的校验和
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// ECDSA 基于椭圆曲线
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	// 生成随机公钥和私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
