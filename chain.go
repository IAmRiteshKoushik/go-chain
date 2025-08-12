package main

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockchain() *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
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
		} else {
			tip = b.Get([]byte("l")) // stores the hash of the tip / HEAD
		}

		return nil
	})

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

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)

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
