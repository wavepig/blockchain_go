package block

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// 区块链保持区块序列
type Blockchain struct {
	tip []byte // 追后端的hash 值
	DB  *bolt.DB
}

// AddBlock将提供的数据保存为区块链中的一个块
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// 这是另一种（只读）类型的BoltDB事务
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Block added. Hash: %x\n", newBlock.Hash)
}

// NewBlockchain用genesis Block创建新的区块链
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}

			// 存储以genesis Block为起点的区块链 key = hash
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return err
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// BlockchainIterator用于对区块链块进行迭代
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.DB}

	return bci
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
