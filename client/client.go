package client

import (
	"encoding/hex"
	"ethereum/txn"
	"ethereum/util"
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

func (c client) GetTransaction(hash string) (txn.BlockTransaction, error) {
	//var raw json.RawMessage
	var rawBlockTxn = struct {
		BlockHash        string `json:"blockHash"`
		BlockNumber      string `json:"blockNumber"`
		From             string `json:"from"`
		Gas              string `json:"gas"`
		GasPrice         string `json:"gasPrice"`
		Hash             string `json:"hash"`
		Input            string `json:"input"`
		Nonce            string `json:"nonce"`
		To               string `json:"to"`
		TransactionIndex string `json:"transactionIndex"`
		Value            string `json:"value"`
		V                string `json:"v"`
		R                string `json:"r"`
		S                string `json:"s"`
	}{}
	err := c.Call(&rawBlockTxn, "eth_getTransactionByHash", hash)
	if err != nil {
		return txn.BlockTransaction{}, err
	}

	return txn.BlockTransaction{
		BlockHash:   rawBlockTxn.BlockHash,
		BlockNumber: util.HexToUint64(rawBlockTxn.BlockNumber),
		From:        rawBlockTxn.From,
		Gas:         util.HexToBigInt(rawBlockTxn.Gas),
		GasPrice:    util.HexToBigInt(rawBlockTxn.GasPrice),
		Hash:        rawBlockTxn.Hash,
		Input: func(h string) []byte {
			b, _ := hex.DecodeString(h[2:])
			return b
		}(rawBlockTxn.Input),
		Nonce:            util.HexToUint64(rawBlockTxn.Nonce),
		To:               rawBlockTxn.To,
		TransactionIndex: util.HexToUint64(rawBlockTxn.TransactionIndex),
		Value:            util.HexToBigInt(rawBlockTxn.Value),
		V:                int(util.HexToUint64(rawBlockTxn.V)),
		R:                util.HexToBigInt(rawBlockTxn.R),
		S:                util.HexToBigInt(rawBlockTxn.S),
	}, nil
}

// Input a signed transaction, return transaction hash.
func (c client) SendTransaction(t txn.Transaction) (string, error) {
	var result string
	err := c.Call(&result, "eth_sendRawTransaction", "0x"+hex.EncodeToString(t.Encode()))
	return result, err
}