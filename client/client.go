package client

import (
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
