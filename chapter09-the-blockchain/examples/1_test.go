package examples_test

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// CreateMerkle demonstrates <Building a merkle tree>
func CreateMerkle(merkle []*chainhash.Hash) *chainhash.Hash {
	if 0 == len(merkle) {
		return nil
	} else if 1 == len(merkle) {
		return merkle[0]
	}

	for len(merkle) > 1 {
		ell := len(merkle)
		if 0 != ell%2 {
			merkle = append(merkle, merkle[ell-1])
			ell++
		}

		newMerkle := make([]*chainhash.Hash, 0, ell/2)
		for i := 0; i < ell; i += 2 {
			newRoot := blockchain.HashMerkleBranches(merkle[i], merkle[i+1])
			newMerkle = append(newMerkle, newRoot)
		}

		merkle = newMerkle

		fmt.Println("Current merkle hash list:")
		for _, hash := range merkle {
			fmt.Println(" ", hash)
		}
		fmt.Println()
	}

	return merkle[0]
}

func NewHashLiteral(literal string) *chainhash.Hash {
	//data, _ := hex.DecodeString(literal)

	//hash, _ := chainhash.NewHash(data)
	hash, _ := chainhash.NewHashFromStr(literal)

	return hash
}

func ExampleCreateMerkle() {
	txHashes := []*chainhash.Hash{
		NewHashLiteral("0000000000000000000000000000000000000000000000000000000000000000"),
		NewHashLiteral("0000000000000000000000000000000000000000000000000000000000000011"),
		NewHashLiteral("0000000000000000000000000000000000000000000000000000000000000022"),
	}

	merkleRoot := CreateMerkle(txHashes)
	fmt.Println("Result:", merkleRoot)

	// Output:
	// Current merkle hash list:
	//   32650049a0418e4380db0af81788635d8b65424d397170b8499cdc28c4d27006
	//   30861db96905c8dc8b99398ca1cd5bd5b84ac3264a4e1b3e65afa1bcee7540c4
	//
	// Current merkle hash list:
	//   d47780c084bad3830bcdaf6eace035e4c6cbf646d103795d22104fb105014ba3
	//
	// Result: d47780c084bad3830bcdaf6eace035e4c6cbf646d103795d22104fb105014ba3
}
