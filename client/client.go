package client

import (
	"encoding/hex"
	"ethereum/contract"
	"ethereum/txn"
	"ethereum/util"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	*rpc.Client
}

func Dial(url string) (Client, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return Client{}, err
	}
	return Client{c}, nil
}

func (c Client) CallContract(cont contract.Contract, funcname string, inputs []interface{}, output interface{}) error {
	cm := struct {
		Data util.Data
		To   string
	}{
		Data: util.Data(cont.Abi[funcname].Id()),
		To:   cont.Address,
	}

	var result hexutil.Bytes
	if err := c.Call(&result, "eth_call", cm, "latest"); err != nil {
		return err
	}

	return cont.UnmarshalResponse(funcname, result, output)
}

// Always uses "latest" block.
func (c Client) GetBalance(addr string) (*big.Int, error) {
	var result hexutil.Big
	err := c.Call(&result, "eth_getBalance", addr, "latest")
	return (*big.Int)(&result), err
}

func (c Client) GetTransactionCount(addr string) (uint64, error) {
	var result hexutil.Uint64
	err := c.Call(&result, "eth_getTransactionCount", addr, "latest")
	return uint64(result), err
}

func (c Client) GetTransactionReceipt(hash string) (txn.TransactionReceipt, error) {
	var rawTxnReceipt = struct {
		BlockHash         string   `json:"blockHash"`
		BlockNumber       string   `json:"blockNumber"`
		ContractAddress   string   `json:"contractAddress"`
		CumulativeGasUsed string   `json:"cumulativeGasUsed"`
		From              string   `json:"from"`
		GasUsed           string   `json:"gasUsed"`
		Logs              []string `json:"logs"`
		LogsBloom         string   `json:"logsBloom"`
		Root              string   `json:"root"`
		To                string   `json:"to"`
		TransactionHash   string   `json:"transactionHash"`
		TransactionIndex  string   `json:"transactionIndex"`
	}{}
	err := c.Call(&rawTxnReceipt, "eth_getTransactionReceipt", hash)
	if err != nil {
		return txn.TransactionReceipt{}, err
	}

	return txn.TransactionReceipt{
		BlockHash:         rawTxnReceipt.BlockHash,
		BlockNumber:       util.HexToUint64(rawTxnReceipt.BlockNumber),
		ContractAddress:   rawTxnReceipt.ContractAddress,
		CumulativeGasUsed: util.HexToBigInt(rawTxnReceipt.CumulativeGasUsed),
		From:              rawTxnReceipt.From,
		GasUsed:           util.HexToBigInt(rawTxnReceipt.GasUsed),
		Logs:              rawTxnReceipt.Logs,
		LogsBloom:         rawTxnReceipt.LogsBloom,
		Root:              rawTxnReceipt.Root,
		To:                rawTxnReceipt.To,
		TransactionHash:   rawTxnReceipt.TransactionHash,
		TransactionIndex:  util.HexToUint64(rawTxnReceipt.TransactionIndex),
	}, nil
}

func (c Client) GetTransaction(hash string) (txn.BlockTransaction, error) {
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
func (c Client) SendTransaction(t txn.Transaction) (string, error) {
	var result string
	err := c.Call(&result, "eth_sendRawTransaction", "0x"+hex.EncodeToString(t.Encode()))
	return result, err
}
