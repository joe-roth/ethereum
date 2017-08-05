package client

import (
	"encoding/hex"
	"ethereum/txn"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type client struct {
	*rpc.Client
}

func Dial(url string) (client, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return client{}, err
	}
	return client{c}, nil
}

// Always uses "latest" block.
func (c client) GetBalance(addr string) (*big.Int, error) {
	var result hexutil.Big
	err := c.Call(&result, "eth_getBalance", addr, "latest")
	return (*big.Int)(&result), err
}

func (c client) GetTransactionCount(addr string) (uint64, error) {
	var result hexutil.Uint64
	err := c.Call(&result, "eth_getTransactionCount", addr, "latest")
	return uint64(result), err
}

func (c client) SendTransaction(t txn.Transaction) (string, error) {
	var result []byte
	err := c.Call(&result, "eth_sendRawTransaction", hex.EncodeToString(t.Encode()))
	fmt.Printf("result = %+v\n", result)
	return "", err
}
