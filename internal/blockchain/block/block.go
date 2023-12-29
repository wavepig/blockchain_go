package block

import (
	"blockchain_go/internal/blockchain/merkle"
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// 块保留块标头
type Block struct {
	Timestamp     int64          // 当前的时间戳（当区块被创建时）
	Transactions  []*Transaction // Transactions是区块中包含的实际有价值的交易信息
	PrevBlockHash []byte         // 存储前一个区块的哈希
	Hash          []byte         // 区块的哈希
	Nonce         int            //
	Height        int            // 高度以方便和网络上节点比较
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

// HashTransactions 返回区块中交易所有的哈希值 Sum256(hashs)
// 比特币使用了一个更加复杂的技术：它将一个块里面包含的所有交易表示为一个  Merkle tree ，
// 然后在工作量证明系统中使用树的根哈希（root hash）。这个方法能够让我们快速检索一个块里面是否包含了某笔交易，
// 即只需 root hash 而无需下载所有交易即可完成判断。
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := merkle.NewMerkleTree(transactions)

	return mTree.RootNode.Data
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

// NewBlock 创建并返回一个区块Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock 创建并返回创世Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}
