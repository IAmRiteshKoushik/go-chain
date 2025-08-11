# Go-Chain: A Simple Blockchain

This project is a minimalist implementation of a blockchain written in Go. It 
serves as a hands-on learning tool to explore the core concepts of blockchain:

- Blocks: The fundamental units of data storage
- Proof-of-Work: A consensus mechanism for securing the chain
- Hashing: The cryptographic process used to link blocks and ensure data integrity

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
