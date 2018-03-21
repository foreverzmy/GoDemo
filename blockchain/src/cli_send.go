package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {

	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}

	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	cbTx := NewCoinbaseTX(from, "")
	txs := []*Transaction{cbTx, tx}

	// 当挖出一个新块时，UTXO 集就会进行更新
	newBlock := bc.MineBlock(txs)
	fmt.Println("Success!")
}
