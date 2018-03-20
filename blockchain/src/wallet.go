package main

import (
	"crypto/ecdsa"
)

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
}
