package block

import (
	"github.com/boltdb/bolt"
	"log"
)

// // BlockchainIterator 用于对区块链块进行迭代
type BlockchainIterator struct {
	currentHash []byte
	DB          *bolt.DB
}

// Next 返回下一个块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.DB.View(func(tx *bolt.Tx) error {
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
