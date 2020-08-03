package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/onflow/flow-go-sdk/crypto"
)

func GeneratePrivateKey(sigAlgoName string) string {
	seed := make([]byte, crypto.MinSeedLength)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	sigAlgo := crypto.StringToSignatureAlgorithm(sigAlgoName)
	privateKey, err := crypto.GeneratePrivateKey(sigAlgo, seed)
	if err != nil {
		panic(err)
	}

	privateKeyHex := hex.EncodeToString(privateKey.Encode())

	return privateKeyHex
}

func main() {
	privateKey := GeneratePrivateKey("ECDSA_P256")
	fmt.Println(privateKey)
}
