package main

import (
	"fmt"
	"log"
)

// createBlockChain 创建区块链
func (cli *CLI) createBlockChain(address string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is note valid")
	}
	bc := CreateBlockChain(address)
	bc.db.Close()
	fmt.Println("Done!")
}
