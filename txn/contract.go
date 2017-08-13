package txn

import (
	"bytes"
	"encoding/binary"
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

func (c Contract) ProcessResponse(funcName string, resp []byte) []interface{} {
	function := c.Functions[funcName]

	buf := bytes.NewBuffer(resp)

	outputs := make([]interface{}, len(function.Outputs))
	for i, output := range function.Outputs {
		switch output.Type {
		case "string":
			// Next 32 of resp will show location of the string.
			loc := make([]byte, 32)
			if n, err := buf.Read(loc); n != 32 {
				panic("not enough bytes in resp")
			} else if err != nil {
				panic("can't read bytes")
			}

			// TODO: naive, what if loc exceeds uint64
			loc_64 := binary.BigEndian.Uint64(loc[24:])

			// Take 32 bytes from location, and that is length.
			// TODO: also, limited to uint64 here
			stringLength := resp[loc_64 : loc_64+32]
			stringLength_64 := binary.BigEndian.Uint64(stringLength[24:])

			stringOut := resp[loc_64+32 : loc_64+32+stringLength_64]
			outputs[i] = string(stringOut)
		default:
			panic("unsupported output type")
		}
	}

	return outputs
}
