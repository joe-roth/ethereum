package txn

import (
	"encoding/hex"
	"ethereum/accnt"
	"ethereum/util"
	"math/big"
	"reflect"
	"testing"
)

func hexToBigInt(hex string) *big.Int {
	i, _ := new(big.Int).SetString(hex, 16)
	return i
}

func TestContractAddress(t *testing.T) {
	var tests = []struct {
		bt       BlockTransaction
		expected string
	}{
		{
			bt: BlockTransaction{
				From:  "0x19e7e376e7c213b7e7e7e46cc70a5dd086daff2a",
				Nonce: 1,
				Input: []byte("test"),
			},
			expected: "0x73b647cba2fe75ba05b8e12ef8f8d6327d6367bf",
		},
	}

	for _, test := range tests {
		if addr := test.bt.ContractAddress(); addr != test.expected {
			t.Fatalf("Expected: %s, received: %s", test.expected, addr)
		}
	}
}

func TestEncodeDecode(t *testing.T) {
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

		if encoded := dTxn.Encode(); !reflect.DeepEqual(encoded, traw) {
			t.Fatalf("Expected: %+v, received: %+v", traw, encoded)
		}
	}
}

func TestSender(t *testing.T) {
	var tests = []struct {
		txn    Transaction
		sender string
	}{
		{
			//"transaction" : {
			//"data" : "",
			//"gasLimit" : "0x5208",
			//"gasPrice" : "0x01",
			//"nonce" : "0x00",
			//"r" : "0x48b55bfa915ac795c431978d8a6a992b628d557da5ff759b307d495a36649353",
			//"s" : "0x1fffd310ac743f371de3b9f7f9cb56c0b28ad43601b4ab949f53faa07bd2c804",
			//"to" : "095e7baea6a6c7c4c2dfeb977efac326af552d87",
			//"v" : "0x1b",
			//"value" : "0x0a"
			// RLP HASH
			// [94 180 245 163 60 98 31 50 168 98 45 95 148 59 107 16 41 148 223 228 229 174 187 239 230 155 177 178 170 15 201 62]
			txn: Transaction{
				Nonce:    0,
				GasPrice: big.NewInt(1), // 2E10 doesn't overflow int64, or else this wouldn't work.
				GasLimit: big.NewInt(21000),
				To:       "0x095e7baea6a6c7c4c2dfeb977efac326af552d87",
				Value:    big.NewInt(0x0a),
				Data:     []byte{},
				V:        0x1b, // 27
				R:        hexToBigInt("48b55bfa915ac795c431978d8a6a992b628d557da5ff759b307d495a36649353"),
				S:        hexToBigInt("1fffd310ac743f371de3b9f7f9cb56c0b28ad43601b4ab949f53faa07bd2c804"),
			},
			sender: "0x963f4a0d8a11b758de8d5b99ab4ac898d6438ea6",
		},
	}

	for _, test := range tests {
		s, err := test.txn.Sender()
		if err != nil {
			t.Fatal(err)
		}

		if s != test.sender {
			t.Fatalf("Expected: %s, received: %s", test.sender, s)
		}
	}
}

func TestSign(t *testing.T) {
	// Create Private account.
	priv, err := accnt.NewAccount("cb4aab9577130f5c4622f355e5c6c3cad2661518ac968c34e4f14a9fde071bfd")
	if err != nil {
		t.Fatal(err)
	}

	// Private Account sign transaction.
	txn := Transaction{
		Nonce:    4,
		GasPrice: big.NewInt(1), // 2E10 doesn't overflow int64, or else this wouldn't work.
		GasLimit: big.NewInt(21000),
		To:       "0x095e7baea6a6c7c4c2dfeb977efac326af552d87",
		Value:    big.NewInt(0x0a),
		Data:     []byte{},
	}
	if err := txn.Sign(priv); err != nil {
		t.Fatal(err)
	}

	// Get sender of signed transaction and compare to private account.
	sender, err := txn.Sender()
	if err != nil {
		t.Fatal(err)
	}

	if sender != priv.Address() {
		t.Fatalf("Expected: %s, received: %s", priv.Address(), sender)
	}
}

func TestHash(t *testing.T) {
	var tests = []struct {
		txn  Transaction
		hash string
	}{
		{
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
			hash: "0548a882856e41ff1bb963032b9e683dd8e45fe7b9344ee045ddfa2712441f8e",
		},
	}

	for _, test := range tests {
		if dh := test.txn.Hash(); dh != test.hash {
			t.Fatalf("Expected: %s, received: %s", test.hash, dh)
		}
	}
}
