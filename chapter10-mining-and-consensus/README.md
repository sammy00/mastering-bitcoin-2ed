# Chapter 10. Mining and Consensus

## Introduction

- Mining is the mechanism that underpins the decentralized **clearinghouse** validating and clearing tx
- Mining is the invention that makes bitcoin special, a decentralized security mechanism that is the basis for P2P digital cash
- Mining secures the bitcoin system and enables the emergence of network-wide consensus without a central authority
- The reward of newly minted coins and transaction fees is an incentive scheme
  - aligns the actions of miners with the security of the network
  - simultaneously implementing the monetary supply
- Mining rewards consist of
  - new coins created with each new block
  - transaction fees from all the transactions included in the block
- Rewards is earned by means of compete to solve a difficult mathematical problem based on a cryptographic hash algorithm.
- The solution to the problem, called the **Proof-of-Work** (a.k.a., PoW), is included in the new block and acts as proof that the miner expended significant computing effort
- The competition to find the PoW to earn the reward and the right to record transactions on the blockchain is the basis for **bitcoin's security model**
- The process is called mining because the reward (new coin generation) is designed to simulate diminishing returns, just like mining for precious metals
- Today, the major incoming of bitcoin mining comes from the newly minted bitcoin (over >99.5%), which will beaten by tx fees gradually

### Bitcoin Economics and Currency Creation

> In practice, a miner may intentionally mine a block taking less than the full reward

The total amount of bitcoin can be calculated as [Example 10-1](examples/1_test.go)

- The finite and diminishing issuance creates a fixed monetary supply that resists inflation.

> **Deflationary Money**
>
> - Deflation is the phenomenon of appreciation of value due to a mismatch in supply and demand that drives up the value (and exchange rate) of a currency
> - The opposite of inflation, price deflation, means that the money has more purchasing power over time
> - In a fiat currency with the possibility of unlimited printing, it is very difficult to enter a deflationary spiral unless there is
>   - a complete collapse in demand, and
>   - an unwillingness to print money
> - Deflation in bitcoin is not caused by a collapse in demand, but by a predictably constrained supply

## Decentralized Consensus

- All traditional payment systems depend on a trust model that has a central authority providing a clearinghouse service, basi‐ cally verifying and clearing all transactions
- Satoshi Nakamoto's main invention is the decentralized mechanism for **emergent consensus**: consensus is an emergent artifact of the asynchronous interaction of thousands of independent nodes, all following simple rules
  - All the properties of bitcoin, including currency, transactions, payments, and the security model that does not depend on central authority or trust, derive from this invention
- Bitcoin's decentralized consensus emerges from the interplay of four processes that occur independently on nodes across the network
  - Independent verification of each transaction, by every full node, based on a comprehensive list of criteria
  - Independent aggregation of those transactions into new blocks by mining nodes, coupled with demonstrated computation through a PoW
  - Independent verification of the new blocks by every node and assembly into a chain
  - Independent selection, by every node, of the chain with the most cumulative computation demonstrated through PoW

## Independent Verification of Transactions

- Before being forwarded to other nodes, every received tx will be first verified against a long checklist as

  - The transaction's syntax and data structure must be correct.
  - Neither lists of inputs or outputs are empty.
  - The transaction size in bytes is less than `MAX_BLOCK_SIZE`.
  - Each output value, as well as the total, must be within the allowed range of values (less than 21m coins, more than the dust threshold)
  - None of the inputs have `hash=0`, `N=–1` (coinbase transactions should not be relayed)
  - `nLocktime` is equal to `INT_MAX`, or `nLocktime` and `nSequence` values are satisfied according to `MedianTimePast`
  - The transaction size in bytes is greater than or equal to 100
  - The number of signature operations (`SIGOPS`) contained in the transaction is less than the signature operation limit
  - The unlocking script (`scriptSig`) can only push numbers on the stack, and the locking script (`scriptPubkey`) must match `isStandard` forms (this rejects "non‐standard" transactions)
  - A matching transaction in the pool, or in a block in the main branch, must exist
  - For each input, if the referenced output exists in any other transaction in the pool, the transaction must be rejected
  - For each input, look in the main branch and the transaction pool to find the referenced output transaction. If the output transaction is missing for any input this will be an orphan transaction. Add to the orphan transactions pool, if a matching transaction is not already in the pool
  - For each input, if the referenced output transaction is a coinbase output, it must have at least `COINBASE_MATURITY` (100) confirmations
  - For each input, the referenced output must exist and cannot already be spent
  - Using the referenced output transactions to get input values, check that each input value, as well as the sum, are in the allowed range of values (less than 21m coins, more than 0)
  - Reject if the sum of input values is less than sum of output values
  - Reject if transaction fee would be too low (`minRelayTxFee`) to get into an empty block
  - The unlocking scripts for each input must validate against the corresponding output locking scripts

  > The conditions change over time, to address new types of denial-of-service attacks or sometimes to relax the rules so as to include more types of transactions

- The unconfirmed and valid tx is stored in _transaction pool_ (a.k.a, _memory pool_, _mempool_) before being progagated to others

## Mining Nodes

- 2 types of miners
  - full nodes
  - non-full nodes joining mining pools
- The new mined block is
  - a checkered flag marking the end of the race
  - the starting pistol in the race for the next block

## Aggregating Transactions into Blocks

- Validated txs residing in **mempool** waits until they can included/mined into a block
- Every mining round,
  - a winning miner would
    1. construct a _candidate block_
    2. struggle to find a PoW w.r.t the candidate block
    3. broadcast the mined block with a valid PoW to others
  - a losing miner would
    1. construct a _candidate block_
    2. struggle to find a PoW w.r.t the candidate block
    3. stop upon receiveing a mined block of the same height as that of the working candidate block
    4. remove the confirmed tx included in the mined block from mempool

> A winning block goes as [Example 10-3](examples/3_test.go)

### The Coinbase Transaction

- **DEFINITION**: the first transaction in any mined block (as [Example 10-4](examples/4_test.go))
- Coinbase tx does not consume (spend) UTXO as inputs. Instead, it has only one input, called the **coinbase**, which creates bitcoin from nothing

### Coinbase Reward and Fees

- `Fee = Sum(Inputs) - Sum(Outputs)`)
- Coinbase reward starts at 50 BTC/block and reduced by half every 210000 blocks (demo as [Example 10-5](examples/5_test.go)

### Structure of the Coinbase Transaction

- Table 10-1. The structure of a "normal" transaction input

| Size               | Field                 | Description                                                      |
| ------------------ | :-------------------- | :--------------------------------------------------------------- |
| 32 bytes           | Transaction Hash      | Pointer to the transaction containing the UTXO to be spent       |
| 4 bytes            | Output Index          | The index number of the UTXO to be spent, first one is 0         |
| 1–9 bytes (VarInt) | Unlocking-Script Size | Unlocking-Script length in bytes, to follow                      |
| Variable           | Unlocking-Script      | A script that fulfills the conditions of the UTXO locking script |
| 4 bytes            | Sequence Number       | Currently disabled Tx-replacement feature, set to `0xFFFFFFFF`   |

- Table 10-2. The structure of a coinbase transaction input

| Size               | Field              | Description                                                                                     |
| ------------------ | ------------------ | ----------------------------------------------------------------------------------------------- |
| 32 bytes           | Transaction Hash   | All bits are zero: Not a transaction hash reference                                             |
| 4 bytes            | Output Index       | All bits are ones: `0xFFFFFFFF`                                                                 |
| 1–9 bytes (VarInt) | Coinbase Data Size | Length of the coinbase data, from 2 to 100 bytes                                                |
| Variable           | Coinbase Data      | Arbitrary data used for extra nonce and mining tags. In v2 blocks; must begin with block height |
| 4 bytes            | Sequence Number    | Set to `0xFFFFFFFF`                                                                             |

### Coinbase Data

- Except for the first few bytes, the rest of the coinbase data can be used by miners in any way they want; it is arbitrary data
- Before BIP-34 is activated, the first few bytes of the coinbase used to be arbitrary
- As per BIP-34, version-2 blocks (blocks with the version field set to 2) must contain the block height index as a script `push` operation in the beginning of the coinbase field
- The ASCII-encoded string `/P2SH/` indicates that the mining node that mined this block supports the **P2SH** improvement defined in BIP-16

An demo to extract the coinbase data from the genesis block goes as [Example 10-6](examples/6_test.go)

## Constructing the Block Header

- The structure of the block header

  | Size     | Field               | Description                                                                      |
  | -------- | ------------------- | -------------------------------------------------------------------------------- |
  | 4 bytes  | Version             | A version number to track software/protocol upgrades                             |
  | 32 bytes | Previous Block Hash | A reference to the hash of the previous (parent) block in the chain              |
  | 32 bytes | Merkle Root         | A hash of the root of the merkle tree of this block's transactions               |
  | 4 bytes  | Timestamp           | The approximate creation time of this block (seconds from Unix Epoch)            |
  | 4 bytes  | Target              | The Proof-of-Work algorithm target (encoded as mantissa-exponent) for this block |
  | 4 bytes  | Nonce               | A counter (initialized as 0) used for the Proof-of-Work algorithm                |

## Mining the Block

- Mining is the process of hashing the block header repeatedly, changing one parameter, until the resulting hash matches a specific target (i.e., when parsed as number, the hash should be smaller than the target)

### Proof-of-Work Algorithm

- The key characteristic of a cryptographic hash algorithm is collision-resistant
- 3 demo of using SHA256
  - [Example 10-8. SHA256 example](examples/8_test.go)
  - [Example 10-9. SHA256 script for generating many hashes by iterating on a nonce](examples/9_test.go)
    - **The nonce is used to vary the output of a cryptographic function**
  - [Example 10-11. Simplified Proof-of-Work implementation](examples/11_test.go)
    - The PoW is a hash less than the target
    - The target and difficulty are inversely related

### Target Representation

### Retargeting to Adjust Difficulty

## Successfully Mining the Block

## Validating a New Block

## Assembling and Selecting Chains of Blocks

### Blockchain Forks

## Mining and the Hashing Race

### The Extra Nonce Solution

### Mining Pools

## Consensus Attacks

## Changing the Consensus Rules

### Hard Forks

### Hard Forks: Software, Network, Mining, and Chain

### Diverging Miners and Difficulty

### Contentious Hard Forks

### Soft Forks

### Criticisms of Soft Forks

## Soft Fork Signaling with Block Version

### BIP-34 Signaling and Activation

### BIP-9 Signaling and Activation

### Consensus Software Developmen
