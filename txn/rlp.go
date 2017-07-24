package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

// DecodeRLP decoded an RLP encoded byte array.
func DecodeRLP(in io.Reader) (interface{}, error) {
	//h := make([]byte, 1)
	t := make([]byte, 1)
	if _, err := in.Read(t); err != nil {
		return nil, err
	}
	h := t[0]

	switch {
	case h < 0x80:
		return h, nil
	case h <= 0xb7:
		s := make([]byte, h-0x80)
		if _, err := in.Read(s); err != nil {
			return nil, err
		}
		return s, nil
	case h <= 0xbf:
		ll := make([]byte, h-0xb7)
		if _, err := in.Read(ll); err != nil {
			return nil, err
		}
		l, _ := binary.Uvarint(ll)
		s := make([]byte, l)
		if _, err := in.Read(s); err != nil {
			return nil, err
		}
		return s, nil
	case h <= 0xf7:
		s := make([]byte, h-0xc0)
		if _, err := in.Read(s); err != nil {
			return nil, err
		}
		b := bytes.NewBuffer(s)
		var list []interface{}
		for b.Len() > 0 {
			s, err := DecodeRLP(b)
			if err != nil {
				return nil, err
			}
			list = append(list, s)
		}
		return list, nil
	default:
		ll := make([]byte, h-0xf7)
		if _, err := in.Read(ll); err != nil {
			return nil, err
		}
		l, _ := binary.Uvarint(ll)
		s := make([]byte, l)
		if _, err := in.Read(s); err != nil {
			return nil, err
		}
		b := bytes.NewBuffer(s)
		var list []interface{}
		for b.Len() > 0 {
			s, err := DecodeRLP(b)
			if err != nil {
				return nil, err
			}
			list = append(list, s)
		}
		return list, nil
	}
}
