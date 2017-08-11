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

// Transaction is a transaction created by user.
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

// BlockTransaction is a representation of a transaction saved on the
// blockchain.
type BlockTransaction struct {
	BlockHash        string
	BlockNumber      uint64
	From             string
	Gas              *big.Int
	GasPrice         *big.Int
	Hash             string
	Input            []byte
	Nonce            uint64
	To               string
	TransactionIndex uint64
	Value            *big.Int
	V                int
	R                *big.Int
	S                *big.Int
}

// Transaction Receipt
type TransactionReceipt struct {
	BlockHash         string
	BlockNumber       uint64
	ContractAddress   string
	CumulativeGasUsed *big.Int
	From              string
	GasUsed           *big.Int
	Logs              []string
	LogsBloom         string
	Root              string
	To                string
	TransactionHash   string
	TransactionIndex  uint64
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
		To: func(i []byte) string {
			if len(i) < 20 {
				return ""
			}
			return "0x" + hex.EncodeToString(vals[3])
		}(vals[3]),
		Value: new(big.Int).SetBytes(vals[4]),
		Data:  vals[5],
		V: func(i []byte) int {
			if len(i) == 0 {
				return 0
			}
			return int(i[0])
		}(vals[6]),
		R: new(big.Int).SetBytes(vals[7]),
		S: new(big.Int).SetBytes(vals[8]),
	}, nil
}

func (t *Transaction) Sender() (string, error) {
	if t.V != 27 && t.V != 28 {
		return "", errors.New("protected txns not yet supported")
	}

	pub, err := accnt.Recover(t.sigHash(), accnt.Signature{
		R: t.R,
		S: t.S,
		V: t.V-27 == 1,
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
		bigIntBytes(t.GasPrice),
		bigIntBytes(t.GasLimit),
		func(addr string) []byte {
			if len(addr) < 42 {
				return []byte{}
			}
			b, err := hex.DecodeString(addr[2:])
			if err != nil {
				panic(err)
			}
			return b
		}(t.To),
		bigIntBytes(t.Value),
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
	return hex.EncodeToString(crypto.Keccak256(t.Encode()))
}

func (t Transaction) Encode() []byte {
	return EncodeRLP([][]byte{
		intToArr(t.Nonce),
		t.GasPrice.Bytes(),
		t.GasLimit.Bytes(),
		func(addr string) []byte {
			// 0x = 2 chars
			// 20 byte address in hex == 40 chars
			if len(addr) < 42 {
				return []byte{}
			}
			b, err := hex.DecodeString(addr[2:])
			if err != nil {
				return []byte{}
			}
			return b
		}(t.To),
		bigIntBytes(t.Value),
		t.Data,
		intToArr(uint64(t.V)),
		bigIntBytes(t.R),
		bigIntBytes(t.S),
	})
}

func bigIntBytes(in *big.Int) []byte {
	if in == nil {
		return []byte{}
	}
	return in.Bytes()
}
