package main

import (
	"fmt"
	"log"
)

// createBlockChain 创建区块链
func (cli *CLI) createBlockChain(address, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is note valid")
	}
	bc := CreateBlockChain(address, nodeID)
	defer bc.db.Close()

	// 当一个新的区块链被创建以后，就会立刻进行重建索引
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("Done!")
}
