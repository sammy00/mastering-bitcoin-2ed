package examples_test

import (
	"crypto/sha256"
	"fmt"
)

func SHA256() {}

// Example 10-8. SHA256 example
func ExampleSHA256() {
	fmt.Printf("%x\n", sha256.Sum256([]byte("I am Satoshi Nakamoto")))

	// Output:
	// 5d7c7ba21cbbcd75d14800b100252d5b428e5b1213d27c385bc141ca6b47989e
}
