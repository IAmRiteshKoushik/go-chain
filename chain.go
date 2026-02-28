package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

// Genesis data is arbitrary. Even in Bitcoin, the first block's txn
// is arbitrary and contains the following data.
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// Check whether the database exists or not
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// As the blockchain does not already exists, we are creating with by
	// mining a genesis block which contains a coinbase transaction
	chainErr := db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			log.Panic(err)
		}

		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			log.Panic(err)
		}

		tip = genesis.Hash
		return nil
	})
	if chainErr != nil {
		log.Panic(err)
	}

	bc := Blockchain{
		tip: tip,
		db:  db,
	}
	return &bc
}

// In theory, our chain could be huge and we would need some kind of
// iterator to print all the blocks. As we are storing everything in BoltDB,
// the storage inside BoltDB is byte ordered but we wish to print them in
// order of minting blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// The iterator points to the tip of the blockchain and thus blocks will be
// obtained from top to bottom
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}

// MineBlock is the replacement of AddBlock() and it takes mines a new block
// with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if err = b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			log.Panic(err)
		}
		if err = b.Put([]byte("l"), newBlock.Hash); err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// Return a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	// Checking all the transactions of each block
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Check if the output is spent
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		// Iterate till Genesis block from the tip
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// Finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
