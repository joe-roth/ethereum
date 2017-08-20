package client

import (
	"ethereum/contract"
	"ethereum/txn"
	"ethereum/util"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCallContract(t *testing.T) {
	var tests = []struct {
		cont        contract.Contract
		funcName    string
		rpcRequest  string
		rpcResponse string
		expected    interface{}
	}{
		{
			cont: contract.Contract{
				Abi: map[string]contract.Function{
					"displayMessage": contract.Function{
						Type: "function",
						Name: "displayMessage",
						Outputs: []contract.Param{
							{Name: "", Type: "string"},
						},
						Constant: true,
						Payable:  false,
					},
				},
				Address: "0xa10a3b175f0f2641cf41912b887f77d8ef34fae8",
			},
			funcName: "displayMessage",
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"Data":"0x2d59dc12","To":` +
				`"0xa10a3b175f0f2641cf41912b887f77d8ef34fae8"},"latest"]}`,
			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":"0x00000000000000000000000000000000000000000000000000` +
				`00000000000020000000000000000000000000000000000000000000000000000000000000002e48656c6c6f2066726f6d2` +
				`06120736d61727420636f6e74726163742063726561746564206279206a6f65212121210000000000000000000000000000` +
				`00000000"}`,
			expected: "Hello from a smart contract created by joe!!!!",
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		var resp string
		if err := c.CallContract(test.cont, test.funcName, nil, &resp); err != nil {
			t.Fatal(err)
		}

		if resp != test.expected {
			t.Fatalf("Expected: %v, received: %v", test.expected, resp)
		}
	}
}

func TestGetTransactionReceipt(t *testing.T) {
	var tests = []struct {
		hash        string
		rpcRequest  string
		rpcResponse string
		expected    txn.TransactionReceipt
	}{
		{
			hash: "0xb0e27987021a059af5f01f17330d2a3c886ddb2e16b62db421ca937c061ebd40",
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionReceipt","params":` +
				`["0xb0e27987021a059af5f01f17330d2a3c886ddb2e16b62db421ca937c061ebd40"]}`,

			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":{"blockHash":` +
				`"0x32dc6ca00fd4b11831bda72c37e91109961604fa7f6d85c85781f22b2d7b49db","blockNumber":"0x3",` +
				`"contractAddress":"0x73b647cba2fe75ba05b8e12ef8f8d6327d6367bf","cumulativeGasUsed":"0x2501b",` +
				`"from":"0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a","gasUsed":"0x2501b","logs":[],"logsBloom":` +
				`"0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000` +
				`0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000` +
				`0000000000000000000","root":"0xe3e256a66a08e9e1c70c3726f5fa08be2c618d04c00216a3ef8a21053956a2d1","to":null,` +
				`"transactionHash":"0xb0e27987021a059af5f01f17330d2a3c886ddb2e16b62db421ca937c061ebd40","transactionIndex":"0x0"}}`,
			expected: txn.TransactionReceipt{
				BlockHash:         "0x32dc6ca00fd4b11831bda72c37e91109961604fa7f6d85c85781f22b2d7b49db",
				BlockNumber:       3,
				ContractAddress:   "0x73b647cba2fe75ba05b8e12ef8f8d6327d6367bf",
				CumulativeGasUsed: big.NewInt(151579),
				From:              "0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a",
				GasUsed:           big.NewInt(151579),
				Logs:              []string{},
				LogsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
					"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
					"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
					"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
					"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
					"000000000000000000",
				Root:             "0xe3e256a66a08e9e1c70c3726f5fa08be2c618d04c00216a3ef8a21053956a2d1",
				To:               "",
				TransactionHash:  "0xb0e27987021a059af5f01f17330d2a3c886ddb2e16b62db421ca937c061ebd40",
				TransactionIndex: 0,
			},
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		btx, err := c.GetTransactionReceipt(test.hash)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(btx, test.expected) {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, btx)
		}
	}
}

func TestGetTransaction(t *testing.T) {
	var tests = []struct {
		hash        string
		rpcRequest  string
		rpcResponse string
		expected    txn.BlockTransaction
	}{
		{
			hash: "0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719",
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByHash","params":` +
				`["0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719"]}`,
			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":{"blockHash":"0x12601e9203cd8b29eb2317d6f645b14b2acafc2564eb24276e75cb4ec66` +
				`67a4d","blockNumber":"0x12ca","from":"0x9d39856f91822ff0bdc2e234bb0d40124a201677","gas":"0x5208","gasPrice":"0x4a817c800",` +
				`"hash":"0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719","input":"0x","nonce":"0x1","to":"0x2c65492bb82` +
				`0552334ba59b4fbb626f35a95e566","transactionIndex":"0x0","value":"0x15af1d78b58c40000","v":"0x1c","r":"0x2083a43ac72ca892e2` +
				`2e805003926850a0da13d9aeb0ef1c4405de35a67d8447","s":"0x6502a37bd5cd629128dd889ee8916acd8d2193f1301c00d9064cbada640a9b58"}}`,
			expected: txn.BlockTransaction{
				BlockHash:        "0x12601e9203cd8b29eb2317d6f645b14b2acafc2564eb24276e75cb4ec6667a4d",
				BlockNumber:      4810,
				From:             "0x9d39856f91822ff0bdc2e234bb0d40124a201677",
				Gas:              big.NewInt(21000),
				GasPrice:         big.NewInt(2E10), // 2E10 doesn't overflow int64, or else this wouldn't work.
				Hash:             "0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719",
				Input:            []byte{},
				Nonce:            1,
				To:               "0x2c65492bb820552334ba59b4fbb626f35a95e566",
				TransactionIndex: 0,
				Value:            util.EthToWei(25),
				V:                28,
				R: new(big.Int).SetBytes([]byte{32, 131, 164, 58, 199, 44, 168, 146, 226, 46, 128, 80, 3, 146, 104, 80, 160, 218, 19,
					217, 174, 176, 239, 28, 68, 5, 222, 53, 166, 125, 132, 71}),
				S: new(big.Int).SetBytes([]byte{101, 2, 163, 123, 213, 205, 98, 145, 40, 221, 136, 158, 232, 145, 106, 205, 141, 33,
					147, 241, 48, 28, 0, 217, 6, 76, 186, 218, 100, 10, 155, 88}),
			},
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		btx, err := c.GetTransaction(test.hash)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(btx, test.expected) {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, btx)
		}
	}
}

func TestSendTransaction(t *testing.T) {
	var tests = []struct {
		transaction txn.Transaction
		rpcRequest  string
		rpcResponse string
		expected    string
	}{
		{
			transaction: txn.Transaction{
				Nonce:    1,
				GasPrice: big.NewInt(2E10), // 2E10 doesn't overflow int64, or else this wouldn't work.
				GasLimit: big.NewInt(21000),
				To:       "0x2c65492bb820552334ba59b4fbb626f35a95e566",
				Value:    util.EthToWei(25),
				Data:     []byte{},
				V:        28,
				R: new(big.Int).SetBytes([]byte{32, 131, 164, 58, 199, 44, 168, 146, 226, 46, 128, 80, 3, 146, 104, 80, 160, 218, 19,
					217, 174, 176, 239, 28, 68, 5, 222, 53, 166, 125, 132, 71}),
				S: new(big.Int).SetBytes([]byte{101, 2, 163, 123, 213, 205, 98, 145, 40, 221, 136, 158, 232, 145, 106, 205, 141, 33,
					147, 241, 48, 28, 0, 217, 6, 76, 186, 218, 100, 10, 155, 88}),
			},
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_sendRawTransaction","params":` +
				`["0xf86d018504a817c800825208942c65492bb820552334ba59b4fbb626f35a95e56689015af1d` +
				`78b58c40000801ca02083a43ac72ca892e22e805003926850a0da13d9aeb0ef1c4405de35a67d84` +
				`47a06502a37bd5cd629128dd889ee8916acd8d2193f1301c00d9064cbada640a9b58"]}`,
			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":"0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719"}`,
			expected:    "0xd866f3672a3cef05f66dec56d30562bbffcc42aa11b54450e6973d52c89d1719",
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		hash, err := c.SendTransaction(test.transaction)
		if err != nil {
			t.Fatal(err)
		}

		if hash != test.expected {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, hash)
		}
	}

}

func TestGetTransactionCount(t *testing.T) {
	var tests = []struct {
		address     string
		rpcRequest  string
		rpcResponse string
		expected    uint64
	}{
		{
			address: "0x9d39856f91822ff0bdc2e234bb0d40124a201677",
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionCount","params":` +
				`["0x9d39856f91822ff0bdc2e234bb0d40124a201677","latest"]}`,
			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":"0x1"}`,
			expected:    uint64(1),
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		count, err := c.GetTransactionCount(test.address)
		if err != nil {
			t.Fatal(err)
		}

		if count != test.expected {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, count)
		}
	}
}

func TestGetBalance(t *testing.T) {
	var tests = []struct {
		address     string
		rpcRequest  string
		rpcResponse string
		expected    *big.Int
	}{
		{
			address: "0x9d39856f91822ff0bdc2e234bb0d40124a201677",
			rpcRequest: `{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0x9d39856f91822ff0bdc2e234bb0d40124a201677",` +
				`"latest"]}`,
			rpcResponse: `{"jsonrpc":"2.0","id":1,"result":"0x34dad6"}`,
			expected:    big.NewInt(3463894),
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			if d := string(data); d != test.rpcRequest {
				t.Fatalf("Expected: %s, received: %s", test.rpcRequest, d)
			}

			if _, err := w.Write([]byte(test.rpcResponse)); err != nil {
				t.Fatal(err)
			}
		}))
		defer ts.Close()

		c, err := Dial(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		b, err := c.GetBalance(test.address)
		if err != nil {
			t.Fatal(err)
		}

		if b.String() != test.expected.String() {
			t.Fatalf("Expected: %+v, received: %+v", test.expected, b)
		}
	}
}
