package main

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	// UI
	cli := CLI{bc}
	cli.Run()
}
