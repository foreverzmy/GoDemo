package main

// Block 由区块头和交易两部分构成
// Timestamp, PrevBlockHash, Hash 属于区块头（block header）
type Block struct {
	Timestamp    int64  // Timestamp: 当前时间戳，也就是区块创建的时间
	PreBlockHash []byte // PrevBlockHash: 前一个块的哈希
	Hash         []byte // Hash: 当前块的哈希
	Data         []byte // Data: 区块实际存储的信息，比特币中也就是交易
}
