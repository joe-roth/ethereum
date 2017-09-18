package main

import "fmt"

// INPUTS
// - leveldb data dir
// - stateRoot
// - output file
// OUTPUTS
// - print num addresses written
// - write balance,nonce,contract as new line in file

var (
	dbDir string
	root  string
	ldb   *ethdb.LDBDatabase
)

func init() {
	dir, _ := os.Getwd()
	flag.StringVar(&dbDir, "db", dir, "Leveldb directory")
	flag.StringVar(&root, "root", "", "Hex formatted, '0x' prefixed root state key")
	flag.Parse()
}

func main() {
	fmt.Println("test")
}
