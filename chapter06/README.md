# Chapter 06. Transactions

## Introduction

- Transactions are data structures that encode the transfer of value between participants in the bitcoin system
- Each transaction is a public entry in bitcoin's blockchain

## Transaction Outputs and Inputs

### Transaction Outputs

- the transaction output is the fundamental building block of a bitcoin tx
- **UTXO**: available and spendable outputs, known as unspent transaction outputs
- receiveing bitcoin is notified when the wallet has detected a UTXO that can be spent with one of the keys controlled by that wallet
- A transaction output can have an arbitrary (integer) value **denominated as a multiple of satoshis**
- **If an UTXO is larger than the desired value of a transaction, it must still be consumed in its entirety and change must be generated in the transaction**
- UTXO selection algorithm is implemented by user's wallet application compose payment

#### coinbase tx

- the first transaction in each block
- placed there by the "winning" miner
- creates brand-new bitcoin payable to that miner as a reward for mining

#### Structure

consist of two parts:

| field          | size                                        | description                                                                                                                |
| -------------- | ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| `value`        | 8 bytes (little-endian)                     | an amount of bitcoin, denominated in satoshis, the smallest bitcoin unit                                                   |
| `scriptPubKey` | content size [1-9 bytes (VarInt)] + content | a cryptographic puzzle (a.k.a. locking script, witness script) that determines the conditions required to spend the output |

> **Serialization** is the process of converting the internal representation of a data structure into a format that can be transmitted one byte at a time, also known as a byte stream. Inversely, **deserialization** or transaction parsing is the process of converting from the byte-stream representation of a transaction to a library's internal representation data structure

### Transaction Inputs

- Transaction inputs identify (by reference) which UTXO will be consumed and provide proof of ownership through an unlocking script

#### Structure

| field       | size                                       | description                                                                         |
| ----------- | ------------------------------------------ | ----------------------------------------------------------------------------------- |
| `txid`      | 32 bytes                                   | a hash referencing the transaction that contains the UTXO being spent               |
| `vout`      | 4 bytes                                    | zero-based output index, identifying which UTXO from that transaction is referenced |
| `scriptSig` | content size[1-9 bytes (Varint)] + content | an unlocking script satisfying the spending conditions set in the UTXO              |
| `sequence`  | 4 bytes                                    | a sequence number (detailed later)                                                  |

> because the value of the input is not explicitly stated, we must also use the referenced UTXO in order to calculate the fees that will be paid in this transaction

### Transaction Fees

- **as a security mechanism themselves**: makes it economically infeasible for attackers to flood the network with transactions
- **as an incentive to include (mine) a transaction into the next block**
- **as a disincentive against abuse of the system** by imposing a small cost on every transaction
- **Transaction fees are calculated based on the size of the transaction in kilobytes, not the value of the transaction in bitcoin**
- Transaction fees affect the processing priority
  - a transaction with sufficient fees is likely to be included in the next block mined
  - a transaction with insufficient or no fees might be delayed, processed on a best-effort basis after a few blocks, or not processed at all
- In Bitcoin Core, fee relay policies are set by the `minrelaytxfee` option. The current default `minrelaytxfee` is 0.00001 bitcoin
- Any bitcoin service that creates transactions, including wallets, exchanges, retail applications, etc., must implement dynamic fees by means of
  - a third-party fee estimation service
  - a built-in fee estimation algorithm.

> Static fees are no longer viable on the bitcoin network. Wallets that set static fees will produce a poor user experience as transactions will often get "stuck" and remain unconfirmed. Users who don't understand bitcoin transactions and fees are dismayed by "stuck" transactions because they think they’ve lost their money.

- Fee estimation algorithms calculate the appropriate fee, based on
  - capacity
  - the fees offered by "competing" transactions

### Adding Fees

- implicitly calculated as

  ```
  Fees = Sum(Inputs) – Sum(Outputs)
  ```

- you must account for all inputs, if necessary by creating change, or you will end up giving the miners a very big tip
- the **expected** fee is independent of the transaction's bitcoin value

## Transaction Scripts and Script Language

- known as _Script_ is a Forth-like reverse-polish notation stack-based execution language
- When a transaction is validated, the **unlocking script** in each input is executed alongside the corresponding **locking script** to see if it satisfies the spending condition.
- Script is a very simple language that was
  - designed to be limited in scope
  - executable on a range of hardware, perhaps as simple as an embedded device
- Pay-to-Public-Key-Hash script is a common script

### Turing Incompleteness

- **DEFINITION**: scripts have limited complexity and predictable execution times
- there are no loops or complex flow control capabilities other than conditional flow control
- purpose: a limited language prevents the transaction validation mechanism from being used as a vulnerability (e.g. an infinite loop trying to carry out DoS attack against the network)

### Stateless Verification

there is no state prior to execution of the script, or state saved after execution of the script

### Script Construction (Lock + Unlock)

- Locking script `scriptPubKey`: a spending condition placed on an output, specifying the conditions that must be met to spend the output in the future
- Unlocking script `scriptSig`

  - satisfies the conditions placed on an output by a locking script and allows the output to be spent
  - usually contains a digital signature produced by the payer

- validation procedure

  1. copy the unlocking script
  2. retrieve the UTXO referenced by the input
  3. copy the locking script from that UTXO
  4. execute the unlocking and locking script
  5. output `true` if the unlocking script satisfies the locking script conditions, and `false` otherwise

- Only outputs referenced by the inputs of valid transactions is considered as "spent" and removed from the UTXO set

#### The script execution stack

- The scripting language executes the script by processing each item from left to right
- **Numbers (data constants)** are pushed onto the stack
- **Operators**

  - push or pop one or more parameters from the stack
  - act on them
  - might push a result onto the stack

- a simple script goes as

  - locking script

    ```
    3 OP_ADD 5 OP_EQUAL
    ```

  - unlocking script

    ```
    2
    ```

  - validation will
    1. combines both as `2 3 OP_ADD OP_EQUAL`
    2. processing the items from left to right to get the final `true`

#### Result of script running

- Transactions are valid if

  - the top result on the stack is
    - `TRUE` (noted as `{0x01}`)
    - any other nonzero value
  - the stack is empty after script execution

- Transactions are invalid if
  - the top value on the stack is `FALSE` (a zero-length empty value, noted as `{}`)
  - script execution is halted explicitly by an operator, such as `OP_VERIFY`, `OP_RETURN`, or a conditional terminator such as `OP_ENDIF`

#### Separate execution of unlocking and locking scripts

- before 2010, the unlocking and locking scripts were concatenated and executed in sequence
- afterwards, separate execution is implemented to thwart a vulnerability that allowed a malformed unlocking script to push data onto the stack and corrupt the locking script
- implementation
  1.  the unlocking script is executed using the stack execution engine
  2.  if step 1 goes with errors, abort and return `false`
  3.  the main stack (not the alternate stack) is copied and the locking script is executed

### Pay-to-Public-Key-Hash (P2PKH)

- locking script as

  ```
  OP_DUP OP_HASH160 <PKH> OP_EQUALVERIFY OP_CHECKSIG
  ```

  where `PKH` is the hash of public key of the payee, i.e., a bitcoin address without Base58Check encoding

- unlocking script as

  ```
  Sig PK
  ```

## Digital Signatures (ECDSA)

- ECDSA is used by the script functions

  - `OP_CHECKSIG`
  - `OP_CHECKSIGVERIFY`
  - `OP_CHECKMULTISIG`
  - `OP_CHECKMULTISIGVERIFY`

- A digital signature serves three purposes in bitcoin
  - spending is authorized by the owner of the private key
  - the proof of authorization is undeniable (nonrepudiation)
  - the transaction (or specific parts of the transaction) have not and cannot be modified by anyone after it has been signed.
- each transaction input is signed independently

### How Digital Signatures Work

two parts

- **signing**: creating a signature, using a private key (the signing key), from a message (the transaction)

  ```
  Sig = Sign(Hash(m), d)
  ```

  where:

  - `d` is the signing private key
  - `m` is the transaction (or parts of it)
  - `Hash` is the hashing function
  - `Sign` is the signing algorithm
  - `Sig` is the resulting signature, and serialized using DER encoding

- **verification**: verify the signature, given also the message and a public key

### Signature Hash Types (`SIGHASH`)

- The signature implies a commitment by the signer to specific transaction data
- `SIGHASH` flag indicates which part of a transaction's data is included in the hash signed by the private key, and is a single byte that is appended to the signature
- Many of the `SIGHASH` flag types only make sense if you think of multiple participants collaborating outside the bitcoin network and updating a partially signed transaction

#### 3 basic flags

| `SIGHASH` flag | Value | Description                                                                                            |
| -------------- | ----- | ------------------------------------------------------------------------------------------------------ |
| `ALL`          | 0x01  | Signature applies to all inputs and outputs                                                            |
| `NONE`         | 0x02  | Signature applies to all inputs, none of the outputs                                                   |
| `SINGLE`       | 0x03  | Signature applies to all inputs but only the one output with the same index number as the signed input |

#### `ANYONECANPAY` modifier

- When `ANYONECANPAY` is set, **ONLY ONE INPUT IS SIGNED**, leaving the rest (and their sequence numbers) open for modification
- `SIGHASH` types with modifiers and their meanings

  | `SIGHASH` flag         | Value | Description                                                              |
  | ---------------------- | ----- | ------------------------------------------------------------------------ |
  | `ALL\|ANYONECANPAY`    | 0x81  | Signature applies to one inputs and all outputs                          |
  | `NONE\|ANYONECANPAY`   | 0x82  | Signature applies to one inputs, none of the outputs                     |
  | `SINGLE\|ANYONECANPAY` | 0x83  | Signature applies to one input and the output with the same index number |

#### Detailed Implementation

1. A copy of the transaction is made
2. Certain fields within are truncated (set to zero length and emptied)
3. The resulting transaction is serialized
4. The `SIGHASH` flag is added to the end of the serialized transaction
5. The result is hashed
6. The hash itself is the "message" that is signed

#### Use cases

| Flag                 | Use case                                                    |
| -------------------- | ----------------------------------------------------------- |
| `ALL\|ANYONECANPAY`  | Crowdfunding by specifying a single output of "goal" amount |
| `NONE`               | Bearer check/blank check with payee unspecified             |
| `NONE\|ANYONECANPAY` | Dust collector                                              |

## ECDSA Math

given

- `d`: signing private key
- `m`: the tx data
- `p`: the prime order of the elliptic curve
- `G`: the elliptic curve generator point
- `Q`: public key corresponds to `d`

during signing

1. generates an ephemeral (temporary) private public key pair `(k, P=kG)`
2. let `R` be the x coordinate of `P`
3. estimates `S=k^-1 (Hash(m) + d*R) mod p`
4. construct signature as `(R, S)`

during verification

1. estimates `P = S^-1 Hash(m)*G + S^-1 * R * Q`
2. output `true` if the x coordinate of `P` equals to `R`

## The Importance of Randomness in Signatures

- Reuse of the same value for k in a signature algorithm leads to exposure of the private key
- the industry best practice is to not generate `k` with a random-number generator seeded with entropy, but instead to use a deterministic-random process seeded with the transaction data itself

## Bitcoin Addresses, Balances, and Other Abstractions

- the information presented to users through wallet applications, blockchain explorers, and other bitcoin user interfaces is often composed of higher-level abstractions
- these abstractions are derived by
  - searching many different transactions
  - inspecting their content
  - manipulating the data contained within them
