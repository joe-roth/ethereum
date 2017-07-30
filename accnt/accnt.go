package accnt

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type Private struct {
	*ecdsa.PrivateKey
	Public
}

type Public ecdsa.PublicKey

// Compressed signature allows for the recovery of a public key given the hash
// of a message.
type Signature struct {
	R *big.Int
	S *big.Int
	V bool // false == 0, true == 1
}

// NewAddress returns an address given hex encoded private key.
func NewAccount(pk string) (Private, error) {
	priv, err := hex.DecodeString(pk)
	if err != nil {
		return Private{}, err
	}

	privk, err := crypto.ToECDSA(priv)
	if err != nil {
		return Private{}, nil
	}

	return Private{
		PrivateKey: privk,
		Public:     Public(privk.PublicKey),
	}, nil
}

func (p Private) Sign(in []byte) (Signature, error) {
	// The produced signature is in the [R || S || V] format where V is 0 or 1.
	sig, err := crypto.Sign(in, p.PrivateKey)
	if err != nil {
		return Signature{}, err
	}

	return Signature{
		R: new(big.Int).SetBytes(sig[:32]),
		S: new(big.Int).SetBytes(sig[32:64]),
		V: func(b byte) bool {
			return b == byte(1)
		}(sig[64]),
	}, nil
}

func Recover(in []byte, s Signature) (Public, error) {
	sig := append(s.R.Bytes(), s.S.Bytes()...)
	sig = append(sig, func(b bool) byte {
		if b {
			return 1
		}
		return 0
	}(s.V))

	pubk, err := crypto.SigToPub(in, sig)
	if err != nil {
		return Public{}, err
	}

	return Public(*pubk), nil
}

func (p Public) Address() string {
	// Create ECDSA public key (64 bytes) as concatenation of x and y points.
	pub := append(p.X.Bytes(), p.Y.Bytes()...)

	// Keccak-256 hash of pub key (32 bytes)
	// ehtereum uses a special Keccak256 configuration with dsbyte: 0x01
	c := crypto.Keccak256(pub)

	// Take last 20 bytes, and prepend '0x'
	return fmt.Sprintf("0x%x", c[12:])
}
