package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

/*

-- Lifecycle of a transaction:
1. There is a genesis block that contains a coinbase transaction. There are no
real inputs in a coinbase transaction, so signing is not necessary. The output
of the coinbase transaction contains a hashed public key

	hashedKey = ripemd16(sha256(pubKey))

2.When one sends coins, a transaction is created. Inputs of the transaction
will referenece outputs from previous transaction[s]. Every input will store a
public key (unhashed) and a signature of the whole transaction.

3. Other nodes in the Bitcoin network that receive the transaction will verify
it. Besides other things, they will check that: the hash of the public key in
an input matches the hash of the referenced output to ensure

4.

*/

// The reward for mining is called subsidy. For now, it is a constant but
// in Bitcoin, it is 50 BTC for Genesis and then every 210_000 blocks, the
// reward is halved
const subsidy = 10

type TxInput struct {
	Txid      []byte // reference of previous output
	Vout      int    // index of an output in the transaction
	ScriptSig string // data to be used in output's ScriptPubKey
}

type TxOutput struct {
	Value        int    // the actual coins
	ScriptPubKey string // a puzzle to lock
}

type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// We check for a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Coinbase transaction is when a miner adds a new block. It is a special
// type of transaction which does not require previously existing outputs
// It creates coins out of nowhere
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{
		Txid:      []byte{}, // no transaction ID for first txn
		Vout:      -1,       // No output
		ScriptSig: data,     // arbitrary data is stored for Coinbase txn
	}
	txout := TxOutput{
		Value:        subsidy, // the output is the reward
		ScriptPubKey: to,      // first coins in someone's wallet
	}

	tx := Transaction{
		ID:   nil,
		Vin:  []TxInput{txin},
		Vout: []TxOutput{txout},
	}
	tx.SetID()

	return &tx
}

// The outputs (coins) of a given transaction input be unlocked using this
func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// The value in the TxOutput can be used up for a future input to another txn
func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// Creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Insufficient funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		// Return the change if any
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
