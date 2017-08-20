package main

import (
	"encoding/hex"
	"ethereum/accnt"
	"ethereum/txn"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/piotrnar/gocoin/lib/secp256k1"
)

func doubleK(account accnt.Private) {
	// Create txn 1.
	t1 := txn.Transaction{
		Nonce:    1,
		GasPrice: big.NewInt(2E10), // 2E10 doesn't overflow int64, or else this wouldn't work.
		GasLimit: big.NewInt(21005),
		To:       "0xb82875007A206D52222887B8Bc21ed309357f878",
		Value:    big.NewInt(1000000000000000),
		Data:     []byte{},
	}

	// Convert txn sig hash to secp256k1 number.
	hash := t1.SigHash()
	var hashNum secp256k1.Number
	hashNum.SetHex(hex.EncodeToString(hash))

	// Generate sig with gocoin.
	var sig secp256k1.Signature
	var recid int
	sig.Sign(
		&secp256k1.Number{*account.PrivateKey.D},
		&hashNum,
		&secp256k1.Number{*big.NewInt(100000)},
		&recid,
	)

	// Append gocoin sig to txn.
	t1.R = &sig.R.Int
	t1.S = &sig.S.Int
	t1.V = recid + 27

	fmt.Printf("t1 = %+v\n", t1)
	fmt.Printf("t1.Hash() = %+v\n", t1.Hash())
	bytes := t1.Encode()
	str := hex.EncodeToString(bytes)
	fmt.Printf("str = %+v\n", str)
}

func main() {
	ac, err := accnt.NewAccount(hex.EncodeToString(crypto.Keccak256([]byte("11fortunefavorsthebold11"))))
	if err != nil {
		panic(err)
	}
	fmt.Printf("ac.Address() = %+v\n", ac.Address())

	doubleK(ac)
}
