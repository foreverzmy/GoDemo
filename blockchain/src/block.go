package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block 由区块头和交易两部分构成，表示区块链中的一个区块
// Timestamp, PrevBlockHash, Hash 属于区块头（block header）
type Block struct {
	Timestamp     int64          // 当前时间戳，也就是区块创建的时间
	Transactions  []*Transaction // 交易
	PrevBlockHash []byte         // 前一个块的哈希
	Hash          []byte         // 当前块的哈希
	Nonce         int            // 对工作量证明进行验证时用到
}

// NewBlock 用于生成新块，参数需要 transactions 与 PrevBlockHash
// 当前块的哈希会基于 transactions 和 PrevBlockHash 计算得到
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run() // 运行工作量证明找到有效哈希

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock 生成创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// HashTransactions 计算区块里所有交易的哈希
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}

	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize 将 Block 序列化为一个字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Print(err)
	}
	return result.Bytes()
}

// DeserializeBlock 将字节数组反序列化为一个 Block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		log.Print(err)
	}

	return &block
}
