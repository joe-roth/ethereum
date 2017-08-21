package contract

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os/exec"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
)

type Contract struct {
	Abi     map[string]Function
	Address string
	Bin     []byte
}

type Function struct {
	Type     string
	Name     string
	Inputs   []Param
	Outputs  []Param
	Constant bool
	Payable  bool
}

func (f Function) Id() []byte {
	// TODO very naive implementation only works with function with no args
	return crypto.Keccak256([]byte(f.Name + "()"))[:4]
}

type Param struct {
	Name string
	Type string
}

func New(abi, address string) (Contract, error) {
	funcs := make([]Function, 0)
	if err := json.Unmarshal([]byte(abi), &funcs); err != nil {
		return Contract{}, err
	}

	c := Contract{
		Address: address,
		Abi:     make(map[string]Function),
	}
	for _, function := range funcs {
		c.Abi[function.Name] = function
	}

	return c, nil
}

func Compile(filename string) (Contract, error) {
	// Contract binary
	cmdOut, err := exec.Command("solc",
		"--combined-json", "abi,bin",
		filename,
	).Output()
	if err != nil {
		return Contract{}, err
	}

	// Extract abi/bin from solc output
	var data struct {
		Contracts map[string]struct {
			Abi string
			Bin string
		}
		Version string
	}
	if err := json.Unmarshal(cmdOut, &data); err != nil {
		return Contract{}, err
	}

	// Take first contract output.
	for _, contract := range data.Contracts {
		cont, err := New(contract.Abi, "")
		if err != nil {
			return Contract{}, err
		}

		bin, err := hex.DecodeString(contract.Bin)
		if err != nil {
			return Contract{}, err
		}
		cont.Bin = bin

		return cont, nil
	}

	return Contract{}, errors.New("no contract in solc output")
}

// TODO for now, only unmarshals into a string
func (c Contract) UnmarshalResponse(funcName string, resp []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid input")
	}
	e := rv.Elem()
	if e.Kind() != reflect.String {
		return errors.New("unexpected input type")
	}

	// Parse the response into an array of interfaces.
	buf := bytes.NewBuffer(resp)

	function := c.Abi[funcName]
	//outputs := make([]interface{}, len(function.Outputs))
	for _, output := range function.Outputs {
		switch output.Type {
		case "string":
			// Next 32 of resp will show location of the string.
			loc := make([]byte, 32)
			if n, err := buf.Read(loc); n != 32 {
				return errors.New("not enough bytes in resp")
			} else if err != nil {
				return errors.New("can't read bytes")
			}

			// TODO: naive, what if loc exceeds uint64
			loc_64 := binary.BigEndian.Uint64(loc[24:])

			// Take 32 bytes from location, and that is length.
			// TODO: also, limited to uint64 here
			stringLength := resp[loc_64 : loc_64+32]
			stringLength_64 := binary.BigEndian.Uint64(stringLength[24:])

			stringOut := resp[loc_64+32 : loc_64+32+stringLength_64]

			e.SetString(string(stringOut))
		case "address":
			// Next 32 of resp will be address.
			address := make([]byte, 32)
			if n, err := buf.Read(address); n != 32 {
				return errors.New("not enough bytes in resp")
			} else if err != nil {
				return errors.New("can't read bytes")
			}

			e.SetString("0x" + hex.EncodeToString(address[12:]))
		default:
			return errors.New("unsupported output type")
		}
	}

	return nil
}
