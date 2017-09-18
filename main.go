package main

import (
	"ethereum/accnt"
	"ethereum/client"
	"ethereum/contract"
	"ethereum/txn"
	"fmt"
	"math/big"
)

var cl client.Client
var accounts []accnt.Private

func init() {
	// 32 bytes
	var pks = []string{
		"1010101010101010101010101010101010101010101010101010101010101010",
		"1111111111111111111111111111111111111111111111111111111111111111",
		"2222222222222222222222222222222222222222222222222222222222222222",
		"3333333333333333333333333333333333333333333333333333333333333333",
		"4444444444444444444444444444444444444444444444444444444444444444",
		"5555555555555555555555555555555555555555555555555555555555555555",
		"6666666666666666666666666666666666666666666666666666666666666666",
		"7777777777777777777777777777777777777777777777777777777777777777",
		"8888888888888888888888888888888888888888888888888888888888888888",
		"9999999999999999999999999999999999999999999999999999999999999999",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}

	accounts = make([]accnt.Private, len(pks))
	for i, pk := range pks {
		ac, err := accnt.NewAccount(pk)
		if err != nil {
			panic(err)
		}
		accounts[i] = ac
	}

	client, err := client.Dial("http://localhost:18545")
	if err != nil {
		panic(err)
	}
	cl = client
}

func main() {
	//callContract()
	//return

	a := accounts[1]

	ctr, err := contract.Compile("escrow.sol")
	if err != nil {
		panic(err)
	}

	nonce, err := cl.GetTransactionCount(a.Address())
	if err != nil {
		panic(err)
	}

	t := txn.Transaction{
		Nonce:    nonce,
		GasPrice: big.NewInt(2E10),    // 2E10 doesn't overflow int64, or else this wouldn't work.
		GasLimit: big.NewInt(3150795), // 3150799 is gas limit
		//Value:    util.EthToWei(10),
		To: "0x560a0c0ca6b0a67895024dae77442c5fd3dc473e",
	}

	//if err := ctr.Deploy(&t, accounts[2].Address(), accounts[3].Address()); err != nil {
	//panic(err)
	//}
	ctr.Call("payoutToSeller", &t)

	if err := t.Sign(a); err != nil {
		panic(err)
	}

	hash, err := cl.SendTransaction(t)
	if err != nil {
		panic(err)
	}

	fmt.Printf("hash = %+v\n", hash)
}

func callContract() {
	ctr, err := contract.Compile("escrow.sol")
	if err != nil {
		panic(err)
	}
	ctr.Address = "0x560a0c0ca6b0a67895024dae77442c5fd3dc473e"

	var resp uint64
	if err := cl.CallContract(ctr, "payoutToSeller", nil, &resp); err != nil {
		panic(err)
	}
	fmt.Printf("resp = %+v\n", resp)
}
