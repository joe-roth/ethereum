package addr

type Account struct {
	PrivateKey []byte
}

// NewAddress returns an address given hex encoded private key.
func NewAccount(pk string) Account {

}

func (a Accout) Address() string {

}

func (a Account) Sign(d []byte) []byte {

}
