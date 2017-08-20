package contract

import (
	"testing"
)

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
