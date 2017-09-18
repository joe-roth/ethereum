package main

import (
	"bytes"
	"encoding/hex"
	"ethereum/util"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

// INPUTS
// - leveldb data dir
// - stateRoot
// - output file
// OUTPUTS
// - print num addresses written
// - write balance,nonce,contract as new line in file

var (
	dbDir       string
	root        string
	outFileName string

	emptyCodeHash = crypto.Keccak256([]byte{})
	ldb           *ethdb.LDBDatabase
	outFile       *os.File
	numProcessed  int
)

func init() {
	dir, _ := os.Getwd()
	flag.StringVar(&dbDir, "db", dir, "Leveldb directory")
	flag.StringVar(&root, "root", "", "Hex formatted, '0x' prefixed root state key")
	flag.StringVar(&outFileName, "o", "output", "File to store the output")
	flag.Parse()
}

func main() {
	if len(root) < 2 {
		panic("please enter hex encoded ldb key")
	}

	if root[:2] != "0x" {
		panic("key must have '0x' prefix")
	}

	key, err := hex.DecodeString(root[2:])
	if err != nil {
		panic("unable to decode key: " + err.Error())
	}

	// Check for existence of CURRENT file in dbDir.  If file doesn't exist, db
	// has not been initialized, and we should exit.
	if _, err := os.Stat(filepath.Clean(dbDir + "/CURRENT")); os.IsNotExist(err) {
		panic("leveldb does not exist in " + dbDir)
	}

	ldb, err = ethdb.NewLDBDatabase(dbDir, 128, 1024)
	if err != nil {
		panic("unable to join ldb: " + err.Error())
	}

	outFile, err = os.OpenFile(outFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic("unable to open file: " + err.Error())
	}
	defer outFile.Close()

	// Recursively traverse every node.
	processNode(key)
}

func processNode(hash []byte) {
	r := getNode(hash)
	switch len(r) {
	case 2:
		processLeaf(r[1])
	case 17:
		for _, v := range r {
			if len(v) == 0 {
				continue
			}
			processNode(v)
		}
	}
}

func getNode(hash []byte) [][]byte {
	r, err := ldb.Get(hash)
	if err != nil {
		panic(err)
	}

	out, err := util.DecodeRLP(bytes.NewBuffer(r))
	if err != nil {
		panic("unable to decode rlp: " + err.Error())
	}

	outArr, ok := out.([]interface{})
	if !ok {
		panic("unable to convert to []interface{}")
	}

	vals := make([][]byte, len(outArr))
	for i, v := range outArr {
		outVal, ok := v.([]byte)
		if !ok {
			panic("unable to convert to []byte")
		}
		vals[i] = outVal
	}

	return vals
}

func processLeaf(account []byte) {
	accountRaw, err := util.DecodeRLP(bytes.NewBuffer(account))
	if err != nil {
		panic(err)
	}

	accountArr, ok := accountRaw.([]interface{})
	if !ok {
		panic("unable to convert to []interface{}")
	}

	vals := make([][]byte, len(accountArr))
	for i, v := range accountArr {
		outVal, ok := v.([]byte)
		if !ok {
			panic("unable to convert to []byte")
		}
		vals[i] = outVal
	}

	balance := new(big.Int).SetBytes(vals[1])
	isContract := !bytes.Equal(vals[3], emptyCodeHash)

	if _, err := outFile.WriteString(fmt.Sprintf("%s,%t\n", balance.String(), isContract)); err != nil {
		panic("can't write to file: " + err.Error())
	}
	numProcessed++
	fmt.Println("Processed", numProcessed)
}
