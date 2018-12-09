package examples_test

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/chaincfg"
)

func ExtractCoinbaseDataFromGensisBlock() {}

/*
Display the genesis block message by Satoshi.
*/
func ExampleExtractCoinbaseDataFromGensisBlock() {
	genesis := chaincfg.MainNetParams.GenesisBlock

	if 1 != len(genesis.Transactions) {
		fmt.Println("too many tx in genesis block")
		return
	}

	coinbaseTx := genesis.Transactions[0]
	if 1 != len(coinbaseTx.TxIn) {
		fmt.Println("too many inputs in coinbase tx")
		return
	}

	coinbaseTxIn := coinbaseTx.TxIn[0]
	data, err := txscript.PushedData(coinbaseTxIn.SignatureScript)
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", data[2])

	// Output:
	// The Times 03/Jan/2009 Chancellor on brink of second bailout for banks
}
