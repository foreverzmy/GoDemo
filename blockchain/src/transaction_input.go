package main

import (
	"bytes"
)

// TXInput 表示交易输入
// Txid: 一个交易输入引用了之前一笔交易的一个输出, 存储之前交易的 ID
// Vout: 一笔交易可能有多个输出，Vout 为输出的索引
// ScriptSig: 提供解锁输出 Txid:Vout 的数据
// PubKey: 公钥
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
