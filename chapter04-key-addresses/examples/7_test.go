package examples_test

import (
	"testing"

	"github.com/SHDMT/btcec"
)

// Example 4-7. A script demonstrating elliptic curve math used for bitcoin keys
func TestEllipticCurveMath(t *testing.T) {
	curve := btcec.S256()

	priv, err := btcec.NewPrivateKey(curve)
	if nil != err {
		t.Fatal(err)
	}

	x, y := curve.ScalarMult(curve.Gx, curve.Gy, priv.D.Bytes())

	pub := priv.PubKey()
	if 0 != pub.X.Cmp(x) || 0 != pub.Y.Cmp(y) {
		t.Fatalf("mismatched public key: got (%s,%s), expect (%s,%s)",
			x.String(), y.String(), pub.X.String(), pub.Y.String())
	}
}
