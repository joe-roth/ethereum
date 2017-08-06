package client

import (
	"ethereum/txn"
	"ethereum/util"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
