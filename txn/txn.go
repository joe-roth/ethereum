package txn

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"ethereum/accnt"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type Transaction struct {
	Nonce    uint64 // i think this can also be a bigint, but not sure.
	GasPrice *big.Int
	GasLimit *big.Int
	To       string // `0xHEX` format
	Value    *big.Int
	Data     []byte
	V        int
	R        *big.Int
	S        *big.Int
}

func Decode(raw []byte) (Transaction, error) {
	r, err := DecodeRLP(bytes.NewBuffer(raw))
	if err != nil {
		return Transaction{}, err
	}

	// RLP array
	a, ok := r.([]interface{})
	if !ok {
		return Transaction{}, fmt.Errorf("decoded raw txn not rlp array")
	}

	vals := make([][]uint8, len(a))
	for i := range a {
		vals[i], ok = a[i].([]uint8)
		if !ok {
			return Transaction{}, fmt.Errorf("error decoding rlp array")
		}
	}

	// this decodes the nonce, which is variable length
	putVarUint64 := func(buf []byte) uint64 {
		if l := len(buf); l < 8 {
			buf = append(make([]byte, 8-l), buf...)
		}
		return binary.BigEndian.Uint64(buf)
	}

	return Transaction{
		Nonce:    putVarUint64(vals[0]),
		GasPrice: new(big.Int).SetBytes(vals[1]),
		GasLimit: new(big.Int).SetBytes(vals[2]),
		To:       "0x" + hex.EncodeToString(vals[3]),
		Value:    new(big.Int).SetBytes(vals[4]),
		Data:     vals[5],
		V:        int(vals[6][0]),
		R:        new(big.Int).SetBytes(vals[7]),
		S:        new(big.Int).SetBytes(vals[8]),
	}, nil
}

func (t *Transaction) Sender() (string, error) {
	if t.V != 27 && t.V != 28 {
		return "", errors.New("protected txns not yet supported")
	}

	pub, err := accnt.Recover(t.sigHash(), accnt.Signature{
		R: t.R,
		S: t.S,
		V: func(i int) bool {
			if i == 1 {
				return true
			}
			return false
		}(t.V - 27),
	})
	if err != nil {
		return "", err
	}

	return pub.Address(), nil
}

// returns hashed RLP of txn which must be signed.  Does not support EIP155.
func (t Transaction) sigHash() []byte {
	return crypto.Keccak256(EncodeRLP([][]byte{
		intToArr(t.Nonce),
		t.GasPrice.Bytes(),
		t.GasLimit.Bytes(),
		func(addr string) []byte {
			b, err := hex.DecodeString(addr[2:])
			if err != nil {
				panic(err)
			}
			return b
		}(t.To),
		t.Value.Bytes(),
		t.Data,
	}))
}

// Will populate V,R,S fields.
func (t *Transaction) Sign(priv accnt.Private) error {
	sig, err := priv.Sign(t.sigHash())
	if err != nil {
		return err
	}

	t.R = sig.R
	t.S = sig.S
	t.V = func(i bool) int {
		if i {
			return 28
		}
		return 27
	}(sig.V)

	return nil
}

func (t Transaction) Hash() string {
	return "test"
}

func (t Transaction) Encode() []byte {
	return EncodeRLP([][]byte{
		intToArr(t.Nonce),
		t.GasPrice.Bytes(),
		t.GasLimit.Bytes(),
		func(addr string) []byte {
			b, err := hex.DecodeString(addr[2:])
			if err != nil {
				panic(err)
			}
			return b
		}(t.To),
		t.Value.Bytes(),
		t.Data,
		intToArr(uint64(t.V)),
		t.R.Bytes(),
		t.S.Bytes(),
	})
}
