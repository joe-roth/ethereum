package util

import (
	"math/big"
	"testing"
)

func TestEthToWei(t *testing.T) {
	bigInt := func(s string) *big.Int {
		i, _ := new(big.Int).SetString(s, 10)
		return i
	}

	var tests = []struct {
		eth int64
		wei *big.Int
	}{
		{
			eth: int64(25),
			wei: bigInt("25000000000000000000"),
		},
	}

	for _, test := range tests {
		if w := EthToWei(test.eth); w.Cmp(test.wei) != 0 {
			t.Fatalf("Expected: %+v, received: %+v", test.wei, w)
		}
	}
}
