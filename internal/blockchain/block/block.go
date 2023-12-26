package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

// 块保留块标头
type Block struct {
	Timestamp     int64  // 当前的时间戳（当区块被创建时）
	Data          []byte // Data是区块中包含的实际有价值的信息
	PrevBlockHash []byte // 存储前一个区块的哈希
	Hash          []byte // 区块的哈希
	Nonce         int
}

// Serialize 序列化块
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock 反序列化块
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// SetHash 计算并设置块哈希
// 获取块字段，将它们连接起来，并在连接的组合上计算SHA-256哈希
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	// Join将s的元素连接起来以创建一个新的字节切片。分隔符sep放置在结果切片中的元素之间。
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewBlock 创建并返回一个区块Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock 创建并返回genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("创世(Genesis)区块", []byte{})
}
