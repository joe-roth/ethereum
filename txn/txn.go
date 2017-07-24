package txn

type Transaction struct {
}

func Decode([]byte) (Transaction, error) {
}

func (t *Transaction) Sign([]byte) {
}

func (t Transaction) Hash() {
}

func (t Transaction) Encode() []byte {
}
