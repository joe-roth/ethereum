package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"ethereum/accnt"
	"ethereum/client"
	"ethereum/txn"
	"fmt"
	"math/big"
	"os/exec"
)

var cl client.Client

func main() {

}

func getAccounts() []accnt.Private {
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

	accnts := make([]accnt.Private, len(pks))
	for i, pk := range pks {
		ac, err := accnt.NewAccount(pk)
		if err != nil {
			panic(err)
		}
		accnts[i] = ac
	}

	return accnts
}

func deployContract(filename string) {
	c, err := compileContract(filename)
	if err != nil {
		panic(err)
	}

	// Create txn of contract.
	t := txn.Transaction{
		GasPrice: big.NewInt(2E10),    // 2E10 doesn't overflow int64, or else this wouldn't work.
		GasLimit: big.NewInt(3150795), // 3150799 is gas limit
		Data: func(h string) []byte {
			d, err := hex.DecodeString(h)
			if err != nil {
				panic("unable to hex decode contract bin:" + err.Error())
			}
			return d
		}(c.Bin),
	}

	accounts := getAccounts()
	account := accounts[1]

	bal, err := cl.GetBalance(account.Address())
	if err != nil {
		panic(err)
	}

	nonce, err := cl.GetTransactionCount(account.Address())
	if err != nil {
		panic(err)
	}
	t.Nonce = nonce

	fmt.Printf("Using account: %s\n\tbalance: %s\n\tnonce: %d\n", account.Address(), bal, nonce)

	if err := t.Sign(account); err != nil {
		panic(err)
	}

	hash, err := cl.SendTransaction(t)
	if err != nil {
		panic(err)
	}

	fmt.Println("Contract deployed!")
	fmt.Println("Txn hash:", hash)
}

type Contract struct {
	Abi string `json:"abi"`
	Bin string `json:"bin"`
}

func compileContract(filename string) (Contract, error) {
	// Contract binary
	// Call `solc filename.sol -bin` for binary data
	fmt.Println("Compiling:", filename)
	cmdName := "solc"
	cmdArgs := []string{filename, "--combined-json", "abi,bin"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return Contract{}, err
	}

	// Extract abi/bin from solc output
	var data struct {
		Contracts map[string]Contract `json:"contracts"`
		Version   string              `json:"version"`
	}
	if err := json.Unmarshal(cmdOut, &data); err != nil {
		return Contract{}, err
	}

	if len(data.Contracts) < 1 {
		return Contract{}, errors.New("no contract in solc output")
	}

	// Take first contract output.
	var c Contract
	for k, contract := range data.Contracts {
		fmt.Println("Using contract:", k)
		c = contract
		break
	}

	return c, nil
}
