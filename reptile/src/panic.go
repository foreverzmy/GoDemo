package main

import (
	"log"
)

// Panic log.Panic
func Panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}
