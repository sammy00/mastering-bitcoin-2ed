package examples_test

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
)

func CalcBlockSubsidy() {}

func ExampleCalcBlockSubsidy() {
	const height = 277316

	subsidy := blockchain.CalcBlockSubsidy(height, &chaincfg.MainNetParams)
	fmt.Println(subsidy)

	// Output:
	// 2500000000
}
