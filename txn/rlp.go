package txn

import (
	"bytes"
	"encoding/binary"
	"io"
)

// DecodeRLP decoded an RLP encoded byte array.
// TODO should return [][]byte, but difficult because of recursion.
func DecodeRLP(in io.Reader) (interface{}, error) {
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
		l := arrToInt(ll)
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
		l := arrToInt(ll)
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
			// when string length is >55 bytes
			// write 0xf7+length of length of string
			// write length of string
			// write string
			l := IntToArr(uint64(h))
			_ = out.WriteByte(byte(0xb7 + len(l)))
			_, _ = out.Write(l)
			_, _ = out.Write(arr)
		}
	}

	ol := out.Len()
	if out.Len() <= 55 {
		return append([]byte{byte(0xc0 + ol)}, out.Bytes()...)
	}

	ola := IntToArr(uint64(ol))
	r := append([]byte{byte(0xf7 + len(ola))}, ola...)
	return append(r, out.Bytes()...)
}

// returns the left-trimmed byte array of the big endian encoding of the given
// uint64
func IntToArr(i uint64) []byte {
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

func arrToInt(a []byte) uint64 {
	if len(a) > 8 {
		return 0
	}

	return binary.BigEndian.Uint64(append(make([]byte, 8-len(a)), a...))
}
