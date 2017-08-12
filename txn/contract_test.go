package txn

import (
	"reflect"
	"testing"
)

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
