# Appendix D. Segregated Witness (segwit)

- **A witness** is more broadly any solution that can satisfy the conditions imposed on an UTXO and unlock that UTXO for spending
- Before segwit's introduction, the witness data was embedded in the transaction as part of each input.
- The term segwit simply means separating the signature or unlocking script of a specific output
- Implementation: move the witness data from the `scriptSig` (unlocking script) field of a transaction into a separate `witness` data structure that accompanies a transaction
- Relevant BIPs

  |                                                                       BIP | Description                                                      |
  | ------------------------------------------------------------------------: | :--------------------------------------------------------------- |
  | [BIP-141](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki) | The main definition of Segregated Witness.                       |
  | [BIP-143](https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki) | Transaction Signature Verification for Version 0 Witness Program |
  | [BIP-144](https://github.com/bitcoin/bips/blob/master/bip-0144.mediawiki) | Peer Services—New network messages and serialization formats     |
  | [BIP-145](https://github.com/bitcoin/bips/blob/master/bip-0145.mediawiki) | getblocktemplate Updates for Segregated Witness (for mining)     |

## Why Segregated Witness?

- Transaction Malleability
  - By moving the witness outside the transaction, the transaction hash used as an identifier no longer includes the witness data
  - With Segregated Witness, transaction hashes become immutable by anyone other than the creator of the transaction
    > Without segwit, attackers may introduce changes to the `scriptSig` of a tx to change it's `txid`  
    > Transaction malleability isn't a problem for most Bitcoin transactions which are designed to be added to the block chain immediately. But it does become a problem when the output from a transaction is spent before that transaction is added to the block chain.
- Script Versioning
  - The addition of a script version number allows the scripting language to be upgraded in a soft-fork way
- Network and Storage Scaling
  - Nodes can prune the witness data after validating the signatures, or ignore it altogether when doing SPV
- Signature Verification Optimization
  - Data-hashing computations increased in `O(n^2)` with respect to the number of signature operations without segwit
  - With segwit, the algorithm is changed to reduce the complexity to `O(n)` (**HOW??**)
- Offline Signing Improvement

  - Segregated Witness signatures incorporate the value (amount) referenced by each input in the hash that is signed, allowing later tx verification without previous txs

## How Segregated Witness Works

- Segregated Witness is also a change to how individual UTXO are spent and therefore is a per-output feature
- A Segregated Witness UTXO specifies a locking script that can be satisfied with witness data outside of the input (segregated)

## Soft Fork (Backward Compatibility)

- To an old wallet or node, a Segregated Witness output looks like an output that anyone can spend

## Segregated Witness Output and Transaction Examples

### Pay-to-Witness-Public-Key-Hash (P2WPKH)

#### For payment

Given a P2PKH output script as

```
DUP HASH160 ab68025513c3dbd2f7b92a94e0581f5d50f654e7 EQUALVERIFY CHECKSIG
```

Its corresponding P2WPKH `scriptPubKey` goes as

```
0 ab68025513c3dbd2f7b92a94e0581f5d50f654e7
```

To a newer, segwit-aware client

- The first number (`0`) is interpreted as a version number (the witness version)
- The second part (20 bytes) is the equivalent of a locking script known as a witness program
  > The 20-byte witness program is simply the hash of the public key, as in a P2PKH script

#### For spending

Unlocking tx to spend the above P2PKH output goes as

```
[...]
"Vin": [
"txid": "0627052b6f28912f2703066a912ea577f2ce4da4caa5a5fbd8a57286c345c2f2",
"vout": 0,
"scriptSig": "<Bob's scriptSig>",
]
[...]
```

To spend the Segwit output, the tx has an empty `scriptSig` and includes a Segregated Witness `witness`, outside the transaction itself

```
[...]
"Vin": [
"txid": "0627052b6f28912f2703066a912ea577f2ce4da4caa5a5fbd8a57286c345c2f2",
"vout": 0,
"scriptSig": "",
]
[...]
"witness": "<Bob's witness data>"
[...]
```

### Wallet construction of P2WPKH

- It is extremely important to note that **P2WPKH should only be created by the payee (recipient)** and not converted by the sender from a known public key, P2PKH script, or address
- P2WPKH outputs MUST be constructed from the hash of a compressed public key

### Pay-to-Witness-Script-Hash (P2WSH)

#### Traditional P2SH

A simple P2SH script with a redeem script defining a 2-of-3 multisig goes as

```
HASH160 54c557e07dde5bb6cb791c7a540e0a4796f5e97e EQUAL
```

whose unlocking script goes as

```
<SigA> <SigB> <2 PubA PubB PubC PubD PubE 5 CHECKMULTISIG>
```

#### P2WSH

- The Segwit program consists of two values pushed to the stack

  - a witness version (`0`)
  - the **32-byte** SHA256 hash of the redeem script

  > While P2SH uses the 20-byte RIPEMD160(SHA256(script)) hash, the P2WSH witness program uses a 32-byte SHA256(script) hash. This difference in the selection of the hashing algorithm is deliberate and is used to differentiate between the two types of witness programs (P2WPKH and P2WSH) by the length of the hash and to provide stronger security to P2WSH (128 bits versus 80 bits of P2SH)

- And the corresponding `scriptPubKey` goes as

  ```
  0 9592d601848d04b172905e0ddb0adde59f1590f1e553ffc81ddc4b0ed927dd73
  ```

- The unlocking tx goes as

  ```json
  [...]
  "Vin": [
  "txid": "abcdef12345...",
  "vout": 0,
    "scriptSig": "",
  ]
  [...]
  "witness": "<SigA> <SigB> <2 PubA PubB PubC PubD PubE 5 CHECKMULTISIG>"
  [...]
  ```

### Differentiating between P2WPKH and P2WSH

- The critical difference between them is the length of the hash
  - The public key hash in P2WPKH is 20 bytes
  - The script hash in P2WSH is 32 bytes

## Upgrading to Segregated Witness

- **HOW**

  1. Wallets must create special segwit type outputs
  2. These outputs can be spent by wallets that know how to construct Segwit transactions

  > For P2WPKH and P2WSH payment types, both the sender and the recipient wallets need to be upgraded to be able to use segwit. Furthermore, the sender's wallet needs to know that the recipient's wallet is segwit-aware

- 2 important scenarios
  - Ability of a sender's wallet that is not segwit-aware to make a payment to a recipient's wallet that can process segwit transactions
  - Ability of a sender's wallet that is segwit-aware to recognize and distinguish between recipients that are segwit-aware and ones that are not, by their addresses

### Embedding Segregated Witness inside P2SH

- **WHEN**
  - Payer isn't segwit-aware
  - Payee is segwit-aware and want segwit payment
- **HOW**: Embed the witness scripts P2WPKH/P2WSH as a P2SH address, which is known as P2SH(P2WPKH)/P2SH(P2WSH)

### P2SH(P2WPKH): Pay-to-Witness-Public-Key-Hash inside Pay-to-Script-Hash

Given a P2WPKH witness program as

```
0 ab68025513c3dbd2f7b92a94e0581f5d50f654e7
```

Embed it into a P2SH script as

```
HASH160 RIPEMD160(SHA256(ab68025513c3dbd2f7b92a94e0581f5d50f654e7)) EQUAL
```

Then the P2SH is encoded as an address to receive payment

### P2SH(P2WSH): Pay-to-Witness-Script-Hash inside Pay-to-Script-Hash

Given a P2WSH witnessage program as

```
0 9592d601848d04b172905e0ddb0adde59f1590f1e553ffc81ddc4b0ed927dd73
```

Embed it into a P2SH script as

```
HASH160 RIPEMD160(SHA256(9592d601848d04b172905e0ddb0adde59f1590f1e553ffc81ddc4b0ed927dd73)) EQUAL
```

Then the P2SH is encoded as an address to receive payment

### Segregated Witness addresses

- There have been a number of proposals for a Segregated Witness address scheme, but none have been actively pursued

### Transaction identifiers

- One of the greatest benefits of Segregated Witness is that it eliminates third-party transaction malleability
- Before segwit, transactions could have their signatures subtly modified by third parties, changing their transaction ID (hash) without changing any fundamental properties (inputs, outputs, amounts).
- With the introduction of Segregated Witness, transactions have two identifiers, `txid` and `wtxid`
- **Only the `txid` of a segwit transaction can be considered immutable by third parties and only if all the inputs of the transaction are segwit inputs**

## Segregated Witness' New Signing Algorithm

- Segwit modifies the semantics of the four signature verification functions
  - `CHECKSIG`
  - `CHECKSIGVERIFY`
  - `CHECKMULTISIG`
  - `CHECKMULTISIGVERIFY`
- Signatures in bitcoin transactions are applied on a **commitment hash**, which is calculated from the transaction data, locking specific parts of the data indicating the signer's commitment to those values
- The new algorithm (specified in BIP-143) achieves 2 important goals
  - The number of hash operations increases by a much more gradual `O(n)` to the number of signature operations
    - Reducing the opportunity to create denial-of-service attacks with overly complex transactions
  - The commitment hash now also includes the value (amounts) of each input as part of the commitment

## Economic Incentives for Segregated Witness

- Fees are intended to align the needs of bitcoin users with the burden their transactions impose on the network, through a market-based price discovery mechanism
- From the perspective of full nodes and miners, some parts of a transaction carry much higher costs
- Every transaction added to the bitcoin network affects the consumption of four resources on nodes
  - **Disk space** about storing txs
  - **CPU** about validating txs
  - **Bandwidth** about transmitting txs
  - **Memory** (most expensive) about maintaining the UTXO set
- The most expensive part of a transaction are the newly created outputs, as they are added to the in-memory UTXO set
- Signatures introduce the least burden because witness data are only validated once and then never used again
- If fees are calculated on transaction size, without discriminating between these two types of data, then the market incentives of fees are not aligned with the actual costs imposed by a transaction
- If the fees are overwhelmingly motivating wallets to use as few inputs as possible in transactions, this can lead to UTXO picking and change address strategies that inadvertently bloat the UTXO set
- **Net-new-UTXO**: `#(txout)-#(txin)`

  - An important metric, as it tells us what impact a transaction will have on the most expensive network-wide resource, the in-memory UTXO set

- An example of

  - `A` is a 3-input/2-output tx locked by a 2-of-3 multisig script
  - `B` is a 2-input/3-output tx locked by a 2-of-3 multisig script
  - Assumption: tx fee is 30 satoshi/byte, and a 75% discount on witness data
  - Without Segregated Witness
    - Transaction A fee: 25,710 satoshi
    - Transaction B fee: 18,990 satoshi
  - With Segregated Witness
    - Transaction A fee: 8,130 satoshi
    - Transaction B fee: 12,045 satoshi

- Segwit two main effects on the fees paid by bitcoin users
  - Reduces the overall cost of transactions by discounting witness data and increasing the capacity of the bitcoin blockchain
  - Segwit's discount on witness data corrects a misalignment of incentives that may have inadvertently created more bloat in the UTXO set
