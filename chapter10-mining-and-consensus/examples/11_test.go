package examples_test

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"
	"time"
)

const maxNonce = math.MaxUint32 + 1

// TODO: parallelize the jobs to speed up running
func ProofOfWork(header string, difficultyBits uint) (string, int) {
	target := new(big.Int).SetInt64(1)
	target = target.Lsh(target, 256-difficultyBits)

	for nonce := 0; nonce < maxNonce; nonce++ {
		//hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", header, nonce)))
		hash := sha256.Sum256([]byte(header + strconv.Itoa(nonce)))

		if got := new(big.Int).SetBytes(hash[:]); got.Cmp(target) < 0 {
			fmt.Println("Success with nonce", nonce)
			fmt.Printf("Hash is %x\n", hash)
			return fmt.Sprintf("%x", hash), nonce
		}
	}

	fmt.Printf("Failed after %d (maxNonce) tries\n", maxNonce)
	return "", maxNonce
}

// This would pretty long to run in the desktop computer,
// so the range of difficulty bits for demo has been shrinked as
// [0,24)
func TestProofOfWork(t *testing.T) {
	var (
		nonce int
		hashS string
	)

	for difficultyBits := uint(0); difficultyBits < 24; difficultyBits++ {
		difficulty := 1 << difficultyBits
		fmt.Printf("Difficulty: %d (%d bits)\n", difficulty, difficultyBits)

		fmt.Println("Starting search...")

		// checkpoint the current time
		start := time.Now()

		// make a new block which includes the hash from the previous block
		// we fake a block of transactions - just a string
		newBlock := "test block with transactions" + hashS

		// find a valid nonce for the new block
		hashS, nonce = ProofOfWork(newBlock, difficultyBits)

		// checkpoint how long it took to find a result
		elapsed := time.Since(start)
		fmt.Printf("Elapsed Time: %d.%d seconds\n", elapsed.Truncate(time.Second),
			elapsed.Truncate(time.Microsecond*100))

		if elapsed > 0 {
			fmt.Printf("Hashing Power: %d hashes per second\n",
				time.Duration(nonce+1)*time.Second/elapsed)
		}
	}
}
