package txn

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

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

func TestDecodeRLP(t *testing.T) {
	// From https://github.com/ethereum/wiki/wiki/RLP
	//The string "dog" = [ 0x83, 'd', 'o', 'g' ]
	//The list [ "cat", "dog" ] = [ 0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g' ]
	//The empty string ('null') = [ 0x80 ]
	//The empty list = [ 0xc0 ]
	//The encoded integer 15 ('\x0f') = [ 0x0f ]
	//The encoded integer 1024 ('\x04\x00') = [ 0x82, 0x04, 0x00 ]
	//The set theoretical representation of three, [ [], [[]], [ [], [[]] ] ] = [ 0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0 ]
	//The string "Lorem ipsum dolor sit amet, consectetur adipisicing elit" = [ 0xb8, 0x38, 'L', 'o', 'r', 'e', 'm', ' ', ... , 'e', 'l', 'i', 't' ]
	var tests = []struct {
		input   []byte
		decoded interface{}
	}{
		{[]byte{0x83, byte('d'), byte('o'), byte('g')}, []byte("dog")},
		{[]byte{0xc8, 0x83, byte('c'), byte('a'), byte('t'), 0x83, byte('d'), byte('o'), byte('g')}, [][]byte{[]byte("cat"), []byte("dog")}},
		{[]byte{0x80}, []byte{}},
		{[]byte{0xc0}, []byte{}},
		{[]byte{0x0f}, []byte{0x0f}},
		{[]byte{0x82, 0x04, 0x00}, []byte{0x04, 0x00}},
		//"[[] [[]] [[] [[]]]]"
		{[]byte{0xc7, 0xc0, 0xc1, 0xc0, 0xc3, 0xc0, 0xc1, 0xc0}, []interface{}{[]interface{}{}, [][]interface{}{{}}, []interface{}{[]interface{}{}, [][]interface{}{{}}}}},
		//"[[] [[]] [[] [[]]]]"
		{[]byte{0xb8, 0x38, byte('L'), byte('o'), byte('r'), byte('e'), byte('m'), byte(' '), byte('i'), byte('p'), byte('s'), byte('u'), byte('m'), byte(' '), byte('d'), byte('o'), byte('l'), byte('o'), byte('r'), byte(' '), byte('s'), byte('i'), byte('t'), byte(' '), byte('a'), byte('m'), byte('e'), byte('t'), byte(','), byte(' '), byte('c'), byte('o'), byte('n'), byte('s'), byte('e'), byte('c'), byte('t'), byte('e'), byte('t'), byte('u'), byte('r'), byte(' '), byte('a'), byte('d'), byte('i'), byte('p'), byte('i'), byte('s'), byte('i'), byte('c'), byte('i'), byte('n'), byte('g'), byte(' '), byte('e'), byte('l'), byte('i'), byte('t')}, []byte("Lorem ipsum dolor sit amet, consectetur adipisicing elit")},
	}

	for _, test := range tests {
		decoded, err := DecodeRLP(bytes.NewBuffer(test.input))
		if err != nil {
			t.Fatalf("Error decoding %+v: %s", test.input, err)
		}

		if fmt.Sprintf("%+v", decoded) != fmt.Sprintf("%+v", test.decoded) {
			t.Logf("Expected: %T, received: %T", test.decoded, decoded)
			t.Fatalf("Expected: %+v, received: %+v", test.decoded, decoded)
		}
	}
}

func TestEncodeRLP(t *testing.T) {
	var tests = []struct {
		raw     [][]byte
		encoded []byte
	}{
		{
			raw: [][]byte{
				[]byte{},
				[]byte{4, 168, 23, 200, 0},
				[]byte{82, 8},
				[]byte{133, 114, 105, 166, 60, 171, 190, 63, 120, 6, 90, 137, 134, 213, 68, 34, 253, 73, 240, 139},
				[]byte{1, 90, 241, 215, 139, 88, 196, 0, 0},
				[]byte{},
				[]byte{28},
				[]byte{161, 54, 246, 13, 83, 245, 241, 2, 255, 192, 231, 72, 124, 33, 237, 26, 169, 101, 143, 76, 167, 188, 96, 250, 126, 152, 217, 180, 151, 41, 43, 210},
				[]byte{114, 14, 48, 120, 189, 220, 161, 198, 222, 76, 52, 202, 219, 24, 111, 163, 56, 84, 140, 65, 133, 5, 136, 163, 187, 183, 90, 241, 225, 122, 197, 41},
			},
			encoded: []byte{248, 109, 128, 133, 4, 168, 23, 200, 0, 130, 82, 8, 148, 133, 114, 105, 166, 60, 171, 190, 63, 120, 6, 90, 137, 134, 213, 68, 34, 253, 73, 240, 139, 137, 1, 90, 241, 215, 139, 88, 196, 0, 0, 128, 28, 160, 161, 54, 246, 13, 83, 245, 241, 2, 255, 192, 231, 72, 124, 33, 237, 26, 169, 101, 143, 76, 167, 188, 96, 250, 126, 152, 217, 180, 151, 41, 43, 210, 160, 114, 14, 48, 120, 189, 220, 161, 198, 222, 76, 52, 202, 219, 24, 111, 163, 56, 84, 140, 65, 133, 5, 136, 163, 187, 183, 90, 241, 225, 122, 197, 41},
		},
	}

	for _, test := range tests {
		encoded := EncodeRLP(test.raw)
		if !reflect.DeepEqual(encoded, test.encoded) {
			t.Fatalf("Expected: %+v, received: %+v", test.encoded, encoded)
		}
	}
}
