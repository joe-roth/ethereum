package addr

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	privateKey []byte
}

// NewAddress returns an address given hex encoded private key.
func NewAccount(pk string) (Account, error) {
	priv, err := hex.DecodeString("cb4aab9577130f5c4622f355e5c6c3cad2661518ac968c34e4f14a9fde071bfd")
	if err != nil {
		return Account{}, err
	}

	if len(priv) != 32 {
		return Account{}, fmt.Errorf("private key not 32 bytes")
	}

	return Account{
		privateKey: priv,
	}, nil
}

func (a Account) Address() string {
	// Pubkey is base point times private key...it is x,y point.
	x, y := crypto.S256().ScalarBaseMult(a.privateKey)

	// Create ECDSA public key (64 bytes) as concatenation of x and y points.
	pub := append(x.Bytes(), y.Bytes()...)

	// Keccak-256 hash of pub key (32 bytes)
	// ehtereum uses a special Keccak256 configuration with dsbyte: 0x01
	c := crypto.Keccak256(pub)

	// Take last 20 bytes, and prepend '0x'
	return fmt.Sprintf("0x%x", c[12:])
}

func (a Account) Sign(d []byte) []byte {
	return []byte{}
}
