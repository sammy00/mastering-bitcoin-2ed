package examples_test

import (
	"fmt"
)

// Example 10-1. A script for calculating how much total bitcoin will be issued
func MaxMoney() int {
	const (
		// Original block reward for miners was 50 BTC
		startBlockReward = 50
		// 210000 is around every 4 years with a 10 minute block interval
		rewardInterval = 210000
	)

	// 50 BTC = 50 0000 0000 Satoshis
	var currentReward int = startBlockReward * 1e8
	var total int
	for currentReward > 0 {
		total += rewardInterval * currentReward
		currentReward /= 2
	}

	return total
}

func ExampleMaxMoney() {
	fmt.Println("Total BTC to ever be created:", MaxMoney(), "Satochis")

	// Output:
	// Total BTC to ever be created: 2099999997690000 Satochis
}
