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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
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
	errFileName string

	emptyCodeHash = crypto.Keccak256([]byte{})
	ldb           *ethdb.LDBDatabase
	outFile       *os.File
	errFile       *os.File
	numProcessed  int
)

func init() {
	dir, _ := os.Getwd()
	flag.StringVar(&dbDir, "db", dir, "Leveldb directory")
	flag.StringVar(&root, "root", "", "Hex formatted, '0x' prefixed root state key")
	flag.StringVar(&outFileName, "o", "output", "File to store the output")
	flag.StringVar(&errFileName, "e", "errors", "File to store errors")
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

	ldb, err = ethdb.NewLDBDatabase(dbDir, 128, 128)
	if err != nil {
		panic("unable to join ldb: " + err.Error())
	}

	outFile, err = os.OpenFile(outFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic("unable to open file: " + err.Error())
	}
	defer outFile.Close()

	errFile, err = os.OpenFile(errFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic("unable to open file: " + err.Error())
	}
	defer errFile.Close()

	// Recursively traverse every node.
	processNode(key)
}

func processNode(hash []byte) {
	r, err := getNode(hash)
	if err != nil {
		errFile.WriteString(fmt.Sprintf("Error getting node %+v: %s\n", hash, err))
		return
	}
	switch len(r) {
	case 2:
		if err := processLeaf(r[1]); err != nil {
			errFile.WriteString(fmt.Sprintf("Error processing leaf %+v: %s\n", hash, err))
		}
	case 17:
		for _, v := range r {
			if len(v) == 0 {
				continue
			}
			processNode(v)
		}
	}
}

func getNode(hash []byte) ([][]byte, error) {
	r, err := ldb.Get(hash)
	if err != nil {
		return nil, err
	}

	out, err := util.DecodeRLP(bytes.NewBuffer(r))
	if err != nil {
		return nil, err
	}

	outArr, ok := out.([]interface{})
	if !ok {
		return nil, err
	}

	vals := make([][]byte, len(outArr))
	for i, v := range outArr {
		outVal, ok := v.([]byte)
		if !ok {
			return nil, fmt.Errorf("Unable to process rlp")
		}
		vals[i] = outVal
	}

	return vals, nil
}

type account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash
	CodeHash []byte
}

func processLeaf(in []byte) error {
	var a account
	if err := rlp.Decode(bytes.NewReader(in), &a); err != nil {
		return err
	}

	isContract := !bytes.Equal(a.CodeHash, emptyCodeHash)

	if _, err := outFile.WriteString(fmt.Sprintf("%s,%t\n", a.Balance.String(), isContract)); err != nil {
		return err
	}
	numProcessed++
	fmt.Println("Processed", numProcessed)
	return nil
}
