package txn

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
)

type Contract struct {
	Functions map[string]Function
	Address   string
}

type Function struct {
	Type     string
	Name     string
	Inputs   []Param
	Outputs  []Param
	Constant bool
	Payable  bool
}

func (f Function) id() []byte {
	// TODO very naive implementation only works with function with no args
	return crypto.Keccak256([]byte(f.Name + "()"))[:4]
}

type Param struct {
	Name string
	Type string
}

func NewContract(abi string) (Contract, error) {
	funcs := make([]Function, 0)
	if err := json.Unmarshal([]byte(abi), &funcs); err != nil {
		return Contract{}, err
	}

	var c Contract
	c.Functions = make(map[string]Function)
	for _, function := range funcs {
		c.Functions[function.Name] = function
	}

	return c, nil
}

type CallMessage struct {
	Data Data   `json:"data"`
	To   string `json:"to"`
}

type Data []byte

func (d Data) MarshalText() ([]byte, error) {
	return []byte("0x" + hex.EncodeToString(d)), nil
}

func (c Contract) CallMessage(funcName string) CallMessage {
	return CallMessage{
		To:   c.Address,
		Data: c.Functions[funcName].id(),
	}
}
