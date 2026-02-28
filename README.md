# Go-Chain: A Simple Blockchain

This project is a minimalist implementation of a blockchain written in Go. It 
serves as a hands-on learning tool to explore the core concepts of blockchain:

- Blocks: The fundamental units of data storage
- Proof-of-Work: A consensus mechanism for securing the chain
- Hashing: The cryptographic process used to link blocks and ensure data integrity

## Understanding the DB Structure
We know that the Bitcoin Core uses two "buckets" to store data:
1. `blocks` stores metadata describing all the blocks in the chain
2. `chainstate` stores the state of a chain, which is all currently unspent 
transaction outputs and some metadata

In `blocks`, the `key -> value` pairs are:
1. 'b' + 32-byte block hash -> block index record
2. 'f' + 4-byte file number -> file information record
3. 'l' -> 4-byte file number: the last block file number used
4. 'R' -> 1-byte boolean: whether we are in the process of reindexing
5. 'F' + 1-byte flag name length + flag name string -> 1 byte boolean: various 
    flags can be on or off
6. 't' + 32-byte transaction hash -> transaction index record

In `chainstate`, the `key -> value` pairs are:
1. 'c' + 32-byte transaction hash -> unspent transaction output record for that 
    transaction
2. 'B' -> 32-byte block hash: the block hash up to whcih the database represents
    the unspent transaction outputs

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/IAmRiteshKoushik/go-chain
cd go-chain
```
2. Run the application:
```bash
go run .
```
You will see output in the console as new blocks are mined and added to the 
chain.

## Project Structure
- `main.go`: The application's entry point, where the blockchain is initialized
and new blocks are created
- `chain.go`: Defines the blockchain's data structure, and the genesis block
- `block.go`: Defines the structure of a block and how it's hash is computed
- `go.mod`: The Go module file

## Contributing
This is a personal project that does not have a formal contribution process. Feel
free to fork the repository and use it as a foundation for your own 
experiments with blockchain technology.
