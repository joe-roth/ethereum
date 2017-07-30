package accnt

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	privateKey []byte
}

type Private struct {
	*ecdsa.PrivateKey
	Public
}

type Public ecdsa.PublicKey

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

func (p Public) Address() string {
	// Create ECDSA public key (64 bytes) as concatenation of x and y points.
	pub := append(p.X.Bytes(), p.Y.Bytes()...)

	// Keccak-256 hash of pub key (32 bytes)
	// ehtereum uses a special Keccak256 configuration with dsbyte: 0x01
	c := crypto.Keccak256(pub)

	// Take last 20 bytes, and prepend '0x'
	return fmt.Sprintf("0x%x", c[12:])
}
