package client

import (
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
