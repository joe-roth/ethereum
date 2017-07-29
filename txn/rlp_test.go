package txn

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

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
