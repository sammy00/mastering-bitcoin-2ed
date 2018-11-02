package examples_test

import (
	"testing"

	"github.com/sammy00/bip38/nonec"

	"github.com/sammy00/bip38/encoding"
)

func TestBIP38(t *testing.T) {
	wif := "5J3mBbAH58CpQ3Y5RNJpUKPE62SQ5tfcvU2JpbnkeyhfsYB1Jcn"
	passphrase := "MyTestPassphrase"

	priv, err := encoding.WIFToPrivateKey(wif)
	if nil != err {
		t.Fatal(err)
	}

	got, err := nonec.Encrypt(priv, passphrase, len(priv) == 33)
	if nil != err {
		t.Fatal(err)
	}

	const expect = "6PRTHL6mWa48xSopbU1cKrVjpKbBZxcLRRCdctLJ3z5yxE87MobKoXdTsJ"
	if got != expect {
		t.Fatalf("invalid encrypted: got %s, expect %s", got, expect)
	}
}
