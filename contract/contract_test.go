package contract

import (
	"ethereum/accnt"
	"ethereum/txn"
	"fmt"
	"reflect"
	"testing"
)

func TestDeploy(t *testing.T) {
	// Given a contract.
	ct, err := Compile("test_data/escrow.sol")
	if err != nil {
		t.Fatal(err)
	}

	a1, _ := accnt.NewAddress("0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a")
	a2, _ := accnt.NewAddress("0x1563915e194d8cfba1943570603f7606a3115508")

	var tx txn.Transaction
	if err := ct.Deploy(&tx, a1, a2); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("tx.Data = %+v\n", tx.Data)
}

func TestCompile(t *testing.T) {
	var tests = []struct {
		filename string
		expected Contract
	}{
		{
			filename: "test_data/helloWorld.sol",
			expected: Contract{
				Abi: map[string]Function{
					"displayMessage": Function{
						Type:   "function",
						Name:   "displayMessage",
						Inputs: make([]Param, 0),
						Outputs: []Param{
							{Name: "", Type: "string"},
						},
						Constant: true,
						Payable:  false,
					},
				},
				Bin: []byte{96, 96, 96, 64, 82, 52, 21, 97, 0, 15, 87, 96, 0, 128, 253, 91, 91, 97, 1, 120, 128, 97, 0, 31, 96, 0, 57, 96, 0, 243, 0, 96, 96, 96, 64, 82, 96, 0, 53, 124, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 144, 4, 99, 255, 255, 255, 255, 22, 128, 99, 45, 89, 220, 18, 20, 97, 0, 62, 87, 91, 96, 0, 128, 253, 91, 52, 21, 97, 0, 73, 87, 96, 0, 128, 253, 91, 97, 0, 81, 97, 0, 205, 86, 91, 96, 64, 81, 128, 128, 96, 32, 1, 130, 129, 3, 130, 82, 131, 129, 129, 81, 129, 82, 96, 32, 1, 145, 80, 128, 81, 144, 96, 32, 1, 144, 128, 131, 131, 96, 0, 91, 131, 129, 16, 21, 97, 0, 146, 87, 128, 130, 1, 81, 129, 132, 1, 82, 91, 96, 32, 129, 1, 144, 80, 97, 0, 118, 86, 91, 80, 80, 80, 80, 144, 80, 144, 129, 1, 144, 96, 31, 22, 128, 21, 97, 0, 191, 87, 128, 130, 3, 128, 81, 96, 1, 131, 96, 32, 3, 97, 1, 0, 10, 3, 25, 22, 129, 82, 96, 32, 1, 145, 80, 91, 80, 146, 80, 80, 80, 96, 64, 81, 128, 145, 3, 144, 243, 91, 97, 0, 213, 97, 1, 56, 86, 91, 96, 96, 96, 64, 81, 144, 129, 1, 96, 64, 82, 128, 96, 46, 129, 82, 96, 32, 1, 127, 72, 101, 108, 108, 111, 32, 102, 114, 111, 109, 32, 97, 32, 115, 109, 97, 114, 116, 32, 99, 111, 110, 116, 114, 97, 99, 116, 32, 99, 114, 101, 97, 129, 82, 96, 32, 1, 127, 116, 101, 100, 32, 98, 121, 32, 106, 111, 101, 33, 33, 33, 33, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 129, 82, 80, 144, 80, 91, 144, 86, 91, 96, 32, 96, 64, 81, 144, 129, 1, 96, 64, 82, 128, 96, 0, 129, 82, 80, 144, 86, 0, 161, 101, 98, 122, 122, 114, 48, 88, 32, 123, 114, 216, 78, 207, 106, 78, 177, 155, 106, 59, 170, 105, 195, 234, 3, 147, 121, 77, 111, 171, 253, 110, 18, 230, 172, 191, 108, 176, 167, 176, 9, 0, 41},
			},
		},
	}

	for _, test := range tests {
		cont, err := Compile(test.filename)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(cont, test.expected) {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, cont)
		}
	}
}

func TestUnmarshalResponse(t *testing.T) {
	var tests = []struct {
		abi      string
		funcname string
		data     []byte
		expected interface{}
	}{
		{
			abi: `[{"constant":true,"inputs":[],"name":"displayMessage","outputs":[{"name":"","type":"string"}],"payable":false,` +
				`"type":"function"}]`,
			funcname: "displayMessage",
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 45, 72, 101, 108, 108, 111, 32, 102, 114, 111, 109, 32, 97, 32, 115, 109, 97, 114, 116, 32, 99, 111, 110, 116, 114, 97, 99, 116, 32, 99, 114, 101, 97, 116, 101, 100, 32, 102, 114, 111, 109, 32, 103, 101, 116, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: "Hello from a smart contract created from geth",
		},
	}

	for _, test := range tests {
		contract, err := New(test.abi, "")
		if err != nil {
			t.Fatal(err)
		}

		var v string
		v = "test"
		if err := contract.UnmarshalResponse(test.funcname, test.data, &v); err != nil {
			t.Fatal(err)
		}

		if v != test.expected {
			t.Fatalf("Expected: %v, received: %v", test.expected, v)
		}
	}
}
