package txn

import (
	"bytes"
	"encoding/binary"
	"io"
)

// DecodeRLP decoded an RLP encoded byte array.
// TODO should return [][]byte
func DecodeRLP(in io.Reader) (interface{}, error) {
	//h := make([]byte, 1)
	t := make([]byte, 1)
	if _, err := in.Read(t); err != nil {
		return nil, err
	}
	h := t[0]

	switch {
	case h <= 0x7f:
		return []byte{h}, nil
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

func EncodeRLP(in [][]byte) []byte {
	out := bytes.NewBuffer([]byte{})

	for _, arr := range in {
		h := len(arr)
		switch {
		case h == 1 && arr[0] < 0x80:
			_, _ = out.Write([]byte{arr[0]})
		case h == 1:
			out.Write([]byte{0x81, arr[0]})
		case h <= 55:
			o := append([]byte{0x80 + byte(h)}, arr...)
			_, _ = out.Write(o)
		default:
			// TODO support when string length is >55 bytes
			panic("unsupported")
		}
	}

	// TODO support when out.Len() overflows uint8
	tl := uint8(out.Len())
	if tl <= 55 {
		return append([]byte{0xc0 + tl}, out.Bytes()...)
	}
	return append([]byte{0xf8, tl}, out.Bytes()...)
}

// returns the left-trimmed byte array of the big endian encoding of the given
// uint64
func intToArr(i uint64) []byte {
	o := make([]byte, 8)
	binary.BigEndian.PutUint64(o, i)
	for i, b := range o {
		if b == 0 {
			continue
		}
		return o[i:]
	}
	return []byte{}
}
