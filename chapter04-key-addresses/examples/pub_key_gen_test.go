package examples_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/sammy00/secp256k1/koblitz"
)

func TestGeneratingAPublicKey(t *testing.T) {

	type Expect struct {
		x *big.Int
		y *big.Int
	}

	var expect Expect
	expect.x, _ = new(big.Int).SetString("F028892BAD7ED57D2FB57BF33081D5CFCF6F9ED3D3D7F159C2E2FFF579DC341A", 16)
	expect.y, _ = new(big.Int).SetString("07CF33DA18BD734C600B96A72BBC4749D5141C90EC8AC328AE52DDFE2E505BDB", 16)

	k, err := hex.DecodeString("1E99423A4ED27608A15A2616A2B0E9E52CED330AC530EDCC32C8FFC6A526AEDD")
	if nil != err {
		t.Fatal(err)
	}

	curve := koblitz.S256()
	x, y := curve.ScalarBaseMult(k)

	if 0 != expect.x.Cmp(x) {
		t.Fatalf("invalid x: got %s, expect %s", x.String(), expect.x.String())
	}

	if 0 != expect.y.Cmp(y) {
		t.Fatalf("invalid y: got %s, expect %s", y.String(), expect.y.String())
	}
}
