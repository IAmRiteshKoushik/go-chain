package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var maxNonce = math.MaxInt64 // 9_223_372_036_854_775_807 possibilities

// The difficulty at which a block was mined. This is a measure of how
// challenging it was for miners to solve the cryptographic puzzle needed
// to validate the block. It is usually dynamically adjusted by the
// network based on - number of miners on chain, time taken for previous set
// of blocks and other factors.

// Here, we are keeping it constant.
const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%v\"\n", pow.block.Transactions)
	for nonce < maxNonce { //
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate a block's proof of work against it's nonce
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:]) // converts the byte array to a bigInt pointer

	// Cmp compares x and y and returns:
	//   - -1 if x < y;
	//   - 0 if x == y;
	//   - +1 if x > y.
	isValid := hashInt.Cmp(pow.target) == -1
	// The int obtained from the hash should be less than the target to be valid

	return isValid
}
