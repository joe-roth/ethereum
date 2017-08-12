package util

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
)

func EthToWei(eth int64) *big.Int {
	// wei = eth * 10^18
	e := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return new(big.Int).Mul(big.NewInt(eth), e)
}

func HexToBigInt(h string) *big.Int {
	return new(big.Int).SetBytes(decodeHexString(h))
}

func HexToUint64(h string) uint64 {
	buf := decodeHexString(h)
	if l := len(buf); l < 8 {
		buf = append(make([]byte, 8-l), buf...)
	}
	return binary.BigEndian.Uint64(buf)
}

func decodeHexString(h string) []byte {
	var b []byte
	if len(h)%2 == 1 {
		b, _ = hex.DecodeString("0" + h[2:])
	} else {
		b, _ = hex.DecodeString(h[2:])
	}
	return b
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

func ArrToInt(a []byte) uint64 {
	if len(a) > 8 {
		return 0
	}

	return binary.BigEndian.Uint64(append(make([]byte, 8-len(a)), a...))
}
