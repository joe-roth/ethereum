package txn

import (
	"reflect"
	"testing"
)

func TestContractUnpack(t *testing.T) {
	abi := `[{"constant":true,"inputs":[],"name":"displayMessage","outputs":[{"name":"","type":"string"}],"payable":false,` +
		`"type":"function"}]`

	contract, err := NewContract(abi)
	if err != nil {
		t.Fatal(err)
	}

	resp := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 45, 72, 101, 108, 108, 111, 32, 102, 114, 111, 109, 32, 97, 32, 115, 109, 97, 114, 116, 32, 99, 111, 110, 116, 114, 97, 99, 116, 32, 99, 114, 101, 97, 116, 101, 100, 32, 102, 114, 111, 109, 32, 103, 101, 116, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	out := contract.ProcessResponse("displayMessage", resp)
	out_s, ok := out[0].(string)
	if !ok {
		panic(out)
	}

	expected := "Hello from a smart contract created from geth"
	if out_s != expected {
		t.Fatalf("Expected: %+v, received: %+v", expected, out_s)
	}
}

func TestContractCall(t *testing.T) {
	var tests = []struct {
		abi          string
		functionName string
		expected     CallMessage
	}{
		{
			abi: `[{"constant":true,"inputs":[],"name":"displayMessage","outputs":[{"name":"","type":"string"}],"payable":false,` +
				`"type":"function"}]`,
			functionName: "displayMessage",
			expected: CallMessage{
				Data: []byte{0x2d, 0x59, 0xdc, 0x12},
			},
		},
	}

	for _, test := range tests {
		contract, err := NewContract(test.abi)
		if err != nil {
			t.Fatal(err)
		}

		cm := contract.CallMessage(test.functionName)
		if !reflect.DeepEqual(cm.Data, test.expected.Data) {
			t.Fatalf("Expected: %+v, received: %+v", test.expected.Data, cm.Data)
		}
	}
}
