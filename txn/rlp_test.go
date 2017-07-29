package txn

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDecodeRLP(t *testing.T) {
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
