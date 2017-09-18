package main

import (
	"bytes"
	"encoding/hex"
	"ethereum/util"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	dbDir string
	root  string
)

func init() {
	dir, _ := os.Getwd()
	flag.StringVar(&dbDir, "db", dir, "Leveldb directory")
	flag.StringVar(&root, "root", "", "Hex formatted, '0x' prefixed root state key")
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

	ldb, err := ethdb.NewLDBDatabase(dbDir, 128, 1024)
	if err != nil {
		panic("unable to join ldb: " + err.Error())
	}

	r, err := ldb.Get(key)
	if err != nil {
		panic("unable to get key: " + err.Error())
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

	fmt.Printf("%+v\n", vals)
}
