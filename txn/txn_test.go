package txn

import (
	"encoding/hex"
	"ethereum/util"
	"math/big"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	hexToBigInt := func(hex string) *big.Int {
		i, _ := new(big.Int).SetString(hex, 16)
		return i
	}

	// TODO test when nonce != 0
	var tests = []struct {
		hex string
		txn Transaction
	}{
		{
			hex: "f86d808504a817c80082520894857269a63cabbe3f78065a8986d54422fd49f08b8901" +
				"5af1d78b58c40000801ca0a136f60d53f5f102ffc0e7487c21ed1aa9658f4ca7bc60fa7e98d" +
				"9b497292bd2a0720e3078bddca1c6de4c34cadb186fa338548c41850588a3bbb75af1e17ac529",
			txn: Transaction{
				Nonce:    0,
				GasPrice: big.NewInt(2E10), // 2E10 doesn't overflow int64, or else this wouldn't work.
				GasLimit: big.NewInt(21000),
				To:       "0x857269a63cabbe3f78065a8986d54422fd49f08b",
				Value:    util.EthToWei(25),
				Data:     []byte{},
				V:        28,
				R:        hexToBigInt("a136f60d53f5f102ffc0e7487c21ed1aa9658f4ca7bc60fa7e98d9b497292bd2"),
				S:        hexToBigInt("720e3078bddca1c6de4c34cadb186fa338548c41850588a3bbb75af1e17ac529"),
			},
		},
	}

	for _, test := range tests {
		traw, err := hex.DecodeString(test.hex)
		if err != nil {
			t.Fatal(err)
		}

		dTxn, err := Decode(traw)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(dTxn, test.txn) {
			t.Fatalf("Expected: %+v, received: %+v", test.txn, dTxn)
		}
	}

	//txnHex := "f86d808504a817c80082520894857269a63cabbe3f78065a8986d54422fd49f08b89015af1d78b58c40000801ca0a136f60d53f5f102ffc0e7487c21ed1aa9658f4ca7bc60fa7e98d9b497292bd2a0720e3078bddca1c6de4c34cadb186fa338548c41850588a3bbb75af1e17ac529"

	//TX(0548a882856e41ff1bb963032b9e683dd8e45fe7b9344ee045ddfa2712441f8e)
	//Contract: false
	//From:     110a2729f50791547faa797fa6760a3a749f133b
	//To:       857269a63cabbe3f78065a8986d54422fd49f08b
	//Nonce:    0
	//GasPrice: 0x4a817c800
	//GasLimit  0x5208
	//Value:    0x15af1d78b58c40000
	//Data:     0x
	//V:        0x1c
	//R:        0xa136f60d53f5f102ffc0e7487c21ed1aa9658f4ca7bc60fa7e98d9b497292bd2
	//S:        0x720e3078bddca1c6de4c34cadb186fa338548c41850588a3bbb75af1e17ac529
	//Hex:      f86d808504a817c80082520894857269a63cabbe3f78065a8986d54422fd49f08b89015af1d78b58c40000801ca0a136f60d53f5f102ffc0e7487c21ed1aa9658f4ca7bc60fa7e98d9b497292bd2a0720e3078bddca1c6de4c34cadb186fa338548c41850588a3bbb75af1e17ac529
}
