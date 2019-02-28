package examples_test

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
)

func CalcBlockSubsidy() {}

// Example 10-5. Calculating the block reward—Function GetBlockSubsidy
func ExampleCalcBlockSubsidy() {
	const height = 277316

	subsidy := blockchain.CalcBlockSubsidy(height, &chaincfg.MainNetParams)
	fmt.Println(subsidy)

	// Output:
	// 2500000000
}
