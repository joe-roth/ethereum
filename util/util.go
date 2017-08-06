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
