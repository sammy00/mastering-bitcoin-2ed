package examples_test

import (
	"math/big"
	"testing"
)

// confirm if a point is on the elliptic curve
func TestOnCurve(t *testing.T) {
	p, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)

	x, _ := new(big.Int).SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)

	y, _ := new(big.Int).SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)

	// z=x^3
	z := new(big.Int).Mul(x, x)
	z.Mul(z, x)

	// z+7
	z.Add(z, big.NewInt(7))

	// z+7-y^2
	z.Sub(z, y.Mul(y, y))

	// (x^3+7-y^2)%p==0
	if u := z.Mod(z, p); 0 != big.NewInt(0).Cmp(u) {
		t.Fatalf("got %s, expect 0", u.String())
	}
}
