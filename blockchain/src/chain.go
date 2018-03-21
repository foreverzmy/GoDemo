package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	bolt "github.com/coreos/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 21/Mar/2018 is the best day."

// BlockChain 是一个 Block 指针数组
// tip 这个词本身有事物尖端或尾部的意思，这里指的是存储最后一个块的哈希
// 在链的末端可能出现短暂分叉的情况，所以选择 tip 其实也就是选择了哪条链
// db 存储数据库连接
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

// CreateBlockChain 创建一个新的区块链数据库
// address 用来接收挖出创世块的奖励
func CreateBlockChain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("BlockChian already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

// NewBlockChain 创建一个有创世块的链
func NewBlockChain(address string) *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	// 打开一个 BoltDB 文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 读写事物
	err = db.Update(func(tx *bolt.Tx) error {
		// 获取了存储区块的 bucket：
		// 如果存在，就从中读取 l 键；
		// 如果不存在，就生成创世块
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

// FindSpendableOutputs 从 address 中找到至少 amount 的 UTXO
func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspendTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspendTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}

			}
		}
	}
	return accumulated, unspentOutputs
}

// FindTransaction finds a transaction by its ID
// 通过 ID 找到一笔交易（这需要在区块链上迭代所有区块）
func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and
// returns transactions with spent outputs removed
// 找到未花费输出
func (bc *BlockChain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	// FindUnspentTransactions returns a list of transactions containing unspent outputs
	// 找到未花费输出的交易
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 如果交易输出被花费了
				if spentTXOs[txID] != nil {
					for _, spendOutIdx := range spentTXOs[txID] {
						if spendOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				// 如果该交易输出可以被解锁，即可被花费
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator 迭代器，里面存储了当前迭代的块哈希（currentHash）和数据库的连接（db）
func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

// MineBlock mines a new block with the provided transactions
// 利用提供的交易挖掘新块
func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("VerifyTransaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil

	})

	return newBlock
}

// SignTransaction signs inputs of a Transaction
// 传入一笔交易，找到它引用的交易，然后对它进行签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

// VerifyTransaction verifies transaction input signatures
// 对交易进行验证
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs)
}

// dbExists 数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
