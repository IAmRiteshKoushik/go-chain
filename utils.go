package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Convert an int64 to an array of bytes. The way this is done is by first
// converting it into a stream of characters and then using bytes to represent
// these characters.
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)

	// binary.Write() is writing to a freshly created byte buffer in the
	// big-Endian format. It is converting an int64 into a byte format
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
