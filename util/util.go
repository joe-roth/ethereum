package util

import "math/big"

func EthToWei(eth int64) *big.Int {
	// wei = eth * 10^18
	e := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return new(big.Int).Mul(big.NewInt(eth), e)
}
