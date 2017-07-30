package accnt

import "testing"

func TestAddress(t *testing.T) {
	var tests = []struct {
		pk, address string
	}{
		{"cb4aab9577130f5c4622f355e5c6c3cad2661518ac968c34e4f14a9fde071bfd", "0x9f872283587d655cba2b13f313511aea353903d9"},
	}

	for _, test := range tests {
		acct, err := NewAccount(test.pk)
		if err != nil {
			t.Fatal(err)
		}
		if ad := acct.Address(); ad != test.address {
			t.Fatalf("Expected: %s, received: %s", test.address, ad)
		}
	}
}
