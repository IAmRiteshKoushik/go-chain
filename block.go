package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
	}

	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	// returns the timestamp as a string
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))

	// headers is a matrix which contains 2 slices
	// [prevHash, Data, timestamp]
	// [empty]
	headers := bytes.Join([][]byte{
		b.PrevBlockHash,
		b.Data,
		timestamp,
	}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}
