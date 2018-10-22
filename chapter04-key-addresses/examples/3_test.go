package examples_test

// Example 4-3. Creating a Base58Check-encoded bitcoin address
// from a private key

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"golang.org/x/crypto/ripemd160"

	"github.com/SHDMT/btcec"
	"github.com/sammy00/base58"
)

func toAddress(pub []byte) []byte {
	data1 := sha256.Sum256(pub)

	ripemd := ripemd160.New()
	ripemd.Write(data1[:])

	return ripemd.Sum(nil)
}

func TestBase58Check(t *testing.T) {
	expect := struct {
		pub  []byte
		addr string
	}{}
	expect.pub, _ = hex.DecodeString("0202a406624211f2abbdc68da3df929f938c3399dd79fac1b51b0e4ad1d26a47aa")
	expect.addr = "1PRTTaJesdNovgne6Ehcdu1fpEdX7913CK"

	k, _ := hex.DecodeString("038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776")

	_, pub := btcec.PrivKeyFromBytes(btcec.S256(), k)

	data := pub.SerializeCompressed()
	if !bytes.Equal(data, expect.pub) {
		t.Fatalf("unexpected public key: got %x, expect %x", data, expect.pub)
	}

	if addr := base58.CheckEncode(toAddress(data), 0); addr != expect.addr {
		t.Fatalf("wrong address: got %s, expect %s", addr, expect.addr)
	}
}
