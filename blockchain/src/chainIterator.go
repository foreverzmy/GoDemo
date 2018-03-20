// BoltDB 允许对一个 bucket 里面的所有 key 进行迭代，
// 但是所有的 key 都以字节序进行存储，
// 而且我们想要以区块能够进入区块链中的顺序进行打印。
// 此外，因为我们不想将所有的块都加载到内存中，
// 我们将会一个一个地读取它们，
// 故而，我们需要一个区块链迭代器（BlockchainIterator）
package main

import (
	"log"

	bolt "github.com/coreos/bbolt"
)

// BlockChainIterator 区块链迭代器
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator 迭代器，里面存储了当前迭代的块哈希（currentHash）和数据库的连接（db）
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

// Next 返回链中的下一个块
func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodeBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodeBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
