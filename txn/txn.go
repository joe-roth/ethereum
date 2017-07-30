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

	hash := crypto.Keccak256(EncodeRLP([][]byte{
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

	pub, err := accnt.Recover(hash, accnt.Signature{
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

func (t *Transaction) Sign() {
	//func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	//h := s.Hash(tx)
	//sig, err := crypto.Sign(h[:], prv)
	//if err != nil {
	//return nil, err
	//}
	//return s.WithSignature(tx, sig)
	//}

	// Get hash of tx

	//func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	//return rlpHash([]interface{}{
	//tx.data.AccountNonce,
	//tx.data.Price,
	//tx.data.GasLimit,
	//tx.data.Recipient,
	//tx.data.Amount,
	//tx.data.Payload,
	//s.chainId, uint(0), uint(0),
	//})
	//}

	//func rlpHash(x interface{}) (h common.Hash) {
	//hw := sha3.NewKeccak256()
	//rlp.Encode(hw, x)
	//hw.Sum(h[:0])
	//return h
	//}
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
