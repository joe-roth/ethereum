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

func TestDecodeContract(t *testing.T) {
	// Unsigned contract encoded with RLP
	unsignedRLP := []byte{249, 1, 131, 1, 133, 4, 168, 23, 200, 0, 131, 71, 141, 226, 128, 128, 185, 1, 112, 96, 96, 96, 64, 82, 52, 21, 97, 0, 15, 87, 96, 0, 128, 253, 91, 91, 97, 1, 81, 128, 97, 0, 31, 96, 0, 57, 96, 0, 243, 0, 96, 96, 96, 64, 82, 96, 0, 53, 124, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 144, 4, 99, 255, 255, 255, 255, 22, 128, 99, 45, 89, 220, 18, 20, 97, 0, 62, 87, 91, 96, 0, 128, 253, 91, 52, 21, 97, 0, 73, 87, 96, 0, 128, 253, 91, 97, 0, 81, 97, 0, 205, 86, 91, 96, 64, 81, 128, 128, 96, 32, 1, 130, 129, 3, 130, 82, 131, 129, 129, 81, 129, 82, 96, 32, 1, 145, 80, 128, 81, 144, 96, 32, 1, 144, 128, 131, 131, 96, 0, 91, 131, 129, 16, 21, 97, 0, 146, 87, 128, 130, 1, 81, 129, 132, 1, 82, 91, 96, 32, 129, 1, 144, 80, 97, 0, 118, 86, 91, 80, 80, 80, 80, 144, 80, 144, 129, 1, 144, 96, 31, 22, 128, 21, 97, 0, 191, 87, 128, 130, 3, 128, 81, 96, 1, 131, 96, 32, 3, 97, 1, 0, 10, 3, 25, 22, 129, 82, 96, 32, 1, 145, 80, 91, 80, 146, 80, 80, 80, 96, 64, 81, 128, 145, 3, 144, 243, 91, 97, 0, 213, 97, 1, 17, 86, 91, 96, 64, 128, 81, 144, 129, 1, 96, 64, 82, 128, 96, 27, 129, 82, 96, 32, 1, 127, 72, 101, 108, 108, 111, 32, 102, 114, 111, 109, 32, 97, 32, 115, 109, 97, 114, 116, 32, 99, 111, 110, 116, 114, 97, 99, 116, 0, 0, 0, 0, 0, 129, 82, 80, 144, 80, 91, 144, 86, 91, 96, 32, 96, 64, 81, 144, 129, 1, 96, 64, 82, 128, 96, 0, 129, 82, 80, 144, 86, 0, 161, 101, 98, 122, 122, 114, 48, 88, 32, 69, 32, 181, 156, 128, 210, 60, 150, 98, 230, 67, 156, 182, 78, 154, 19, 90, 148, 222, 135, 99, 28, 92, 59, 243, 170, 170, 60, 168, 112, 210, 241, 0, 41, 128, 128, 128}

	// Create txn
	tx, err := Decode(unsignedRLP)
	if err != nil {
		t.Fatal(err)
	}

	// Encode txn
	raw := tx.Encode()

	// Ensure that encoding is equal to original.
	if !reflect.DeepEqual(raw, unsignedRLP) {
		t.Fatal("not equal")
	}
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
