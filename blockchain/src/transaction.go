// 在比特币中，是先有蛋，然后才有鸡。输入引用输出的逻辑，是经典的“蛋还是鸡”问题：输入先产生输出，然后输出使得输入成为可能。在比特币中，最先有输出，然后才有输入。换而言之，第一笔交易只有输出，没有输入。

// 当矿工挖出一个新的块时，它会向新的块中添加一个 coinbase 交易。coinbase 交易是一种特殊的交易，它不需要引用之前一笔交易的输出。它“凭空”产生了币（也就是产生了新币），这是矿工获得挖出新块的奖励，也可以理解为“发行新币”。

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction 由交易 ID，输入和输出构成
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXInput 包含 3 部分
// Txid: 一个交易输入引用了之前一笔交易的一个输出, 存储之前交易的 ID
// Vout: 一笔交易可能有多个输出，Vout 为输出的索引
// ScriptSig: 提供解锁输出 Txid:Vout 的数据
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TXOutput 包含两部分
// Value: 有多少币，就是存储在 Value 里面
// ScriptPubKey: 对输出进行锁定
// 在当前实现中，ScriptPubKey 将仅用一个字符串来代替
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// IsCoinbase 判断是否是 coinbase 交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID 设置 id
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// CanUnlockOutputWith 这里的 unlockingData 可以理解为地址
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith 是个可以不锁定
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX 构建 coinbase 交易，该没有输入，只有一个输出
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// NewUTXOTransaction 创建一笔新的交易
func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	// 找到足够的未花费输出
	acc, validOutpits := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not Enough Funds.")
	}

	// Build a list of inputs
	for txid, outs := range validOutpits {

		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}

	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})

	// 如果 UTXO 总数超过所需，则产生找零
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx

}
