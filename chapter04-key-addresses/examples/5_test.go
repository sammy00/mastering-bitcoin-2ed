package examples_test

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/sammy00/base58"
	"golang.org/x/crypto/ripemd160"

	"github.com/SHDMT/btcec"
)

func pubKeyToAddress(pub []byte) string {
	data1 := sha256.Sum256(pub)

	ripemd := ripemd160.New()
	ripemd.Write(data1[:])

	y := ripemd.Sum(nil)

	return base58.CheckEncode(y, 0x00)
}

//  Key and address generation and formatting
func TestKeyAndAddressGenerate(t *testing.T) {
	// Generate a random private key
	x, err := btcec.NewPrivateKey(btcec.S256())
	if nil != err {
		t.Fatal(err)
	}
	xx := x.Serialize()

	t.Logf("Private Key (hex) is: %x", xx)
	t.Logf("Private Key (decimal) is: %s",
		new(big.Int).SetBytes(xx).String())

	// Convert private key to WIF format
	t.Log("Private Key (WIF) is:", base58.CheckEncode(xx, 0x80))

	// Add suffix "01" to indicate a compressed private key
	xxC := append(xx, 0x01)
	t.Logf("Private Key Compressed (hex) is: %x", xxC)
	// Generate a WIF format from the compressed private key (WIF-compressed)
	t.Log("Private Key (WIF-Compressed) is:",
		base58.CheckEncode(xxC, 0x80))

	Y := x.PubKey()
	t.Logf("Public Key (x,y) coordinate is: (%v,%v)", Y.X, Y.Y)
	// Encode as hex, prefix 0x04
	yy := Y.SerializeUncompressed()
	t.Logf("Public Key (hex) is: %x", yy)

	yyC := Y.SerializeCompressed()
	// Compress public key, adjust prefix depending on whether y is even or odd
	t.Logf("Compressed Public Key (hex) is: %x", yyC)

	// Generate bitcoin address from public key
	t.Log("Bitcoin Address (b58check) is:", pubKeyToAddress(yy))
	// Generate compressed bitcoin address from compressed public key
	t.Log("Compressed Bitcoin Address (b58check) is:", pubKeyToAddress(yyC))
}
