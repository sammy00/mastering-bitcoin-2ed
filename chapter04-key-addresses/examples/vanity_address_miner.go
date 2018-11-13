// +build ignore

package main

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil"

	"github.com/sammy00/base58"

	"github.com/btcsuite/btcd/btcec"
)

func main() {

	const search = "1kid"

	for {
		priv, err := btcec.NewPrivateKey(btcec.S256())
		if nil != err {
			fmt.Println(err)
			break
		}

		// convert to payment address
		//address := base58.CheckEncode(priv.PubKey().SerializeCompressed(), 0x00)
		address := EncodeToAddress(priv.PubKey())
		fmt.Println(address)

		if strings.HasPrefix(strings.ToLower(address), search) {
			fmt.Println("Found vanity address!", address)
			fmt.Printf("Secret: %x\n", priv.Serialize())
			break
		}
	}
}

func EncodeToAddress(pub *btcec.PublicKey) string {
	data := pub.SerializeCompressed()

	return base58.CheckEncode(btcutil.Hash160(data), 0x00)
}
