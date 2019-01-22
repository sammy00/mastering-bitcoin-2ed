# Chapter12. Blockchain Applications

## Payment Channels and State Channels

- **Payment channels** are a trustless mechanism for exchanging bitcoin transactions between two parties, outside of the bitcoin blockchain. These transactions, which would be valid if settled on the bitcoin blockchain, are held off-chain instead, acting as promissory notes for eventual batch settlement.
- **State channels** are virtual constructs represented by the exchange of state between two parties, outside of the blockchain.
- **channel** is used to represent the relationship and shared state between two parties, outside of the blockchain
- a payment channel is just a series of transactions
- 2 types of payment channels: **unidirectional** and **bidirectional**

### Relationship

- Payment channels are part of the broader concept of a state channel, which represents an off-chain alteration of state, secured by eventual settlement in a blockchain.
- A payment channel is a state channel where the state being altered is the balance of a virtual currency.

### State Channels - Basic Concepts and Terminology

#### funding transaction (a.k.a. anchor transaction)

is composed by **two parties** and then transmitted to the network and **mined** to lock a shared state on the blockchain

#### commitment tx

- signed by both parties to alter the initial state
- are valid transactions that could be submitted for settlement by either party, but instead are held off-chain by each party pending the channel closure
- all previous states are invalidated by the most up-to-date commitment transaction, which is always the only one that can be redeemed.

#### settlement tx

the channel can be closed

- either cooperatively, by submitting a final settlement transaction to the blockchain
- unilaterally, by either party submitting the last commitment transaction to the blockchain

> - only the funding and settlement transactions need to be submitted for mining on the blockchain
> - any intermediate commitment transactions is held offchain

### Simple Payment Channel Example

**assumption**: no cheater and unidirectional

Emma buying the video streaming service (1 millibit/s) from Fabian using a payment channel.

1. establish the channel with a tx (of total amount `a`) of Emma to a 2-of-2 multisignature address, with each of them holding one of the keys
2. for `t` second of video, Emma updates the channel balance in form a new commitment tx
   - crediting `t` millibits to Fabian's address
   - refunding `a-t` millibits back to Emma's address

> every commitment tx are sourced from the same 2-of-2 output from the funding tx

### Making Trustless Channels

- problems of the 'simple payment channel'

  - Once constructed, funds will be lost if one of the parties disconnects before there is at least one commitment transaction signed by both parties.
  - the payer can cheat by broadcasting a prior commitment that is in her favor.

- two solutions
  - timelocks
  - asymmetric revocable commitments

### timelocks-based solution

- tool: the tx-level **timelocks** `nLockTime`

- for payer

  1. signs the **funding tx** privately without transmitting
  2. constructs the **refund tx** timelocked to future and requests payee's signature
  3. once getting the fully signed **refund tx**, transmitting the **funding tx** for mining

- refund tx

  - acts as the first commitment transaction
  - its timelock establishes the upper bound for the channel's life

- newer commitment tx should have shorter timelocks than the older ones and the refund tx, precluding these older ones from being redeemed to execute **double-spend attack**

- settlement

  - cooperative way: either party takes the most recent commitment transaction and builds a settlement transaction that is identical in every way except that it omits the timelock.
  - unilaterally: submit the most recent commitment tx to mine on chain and then wait for refunding until the timelock expires

- **2 disadvantages**
  - the lifetime of the channel is limited by the refund tx and intermediate commitment txs
  - #(tx) is limited due to the monotonically decreasing timelocks enforced by intermediate commitment txs

### asymmetric revocable commitments

- The only way to cancel a transaction is by double-spending its inputs with another transaction before it is mined
- tx can be constructed to be undesirable to use: gives each party a revocation key that can be used to punish the other party if they try to cheat

> the funding tx should exceed the channel capacity to cover the tx fees

#### implementation

Two parties as A and B

- every commitment tx pays the payee immediately whereas forcing the payer to wait for a short timelock to expire
- every commitment tx of A looks like
  - this tx can be constructed by A only since nobody else knows `PK_A`
  - the revocation public key `RPK_A` is created and kept privately by A
  - on state transition, `RPK_A` would be revealed to B

```
Input: 2-of-2 funding output, signed by A

Output 0 <5 bitcoin>:
  PK_A CHECKSIG

Output 1:
  IF
    # Revocation penalty output
    RPK_A
  ELSE
    <1000 blocks> CHECKSEQUENCEVERIFY DROP PK_A
  ENDIF
  CHECKSIG
```

- symmetrically for B's side

```
Input: 2-of-2 funding output, signed by B

Output 0 <5 bitcoin>:
  PK_A CHECKSIG

Output 1:
  IF
    # Revocation penalty output
    <Revocation Public Key created by B>
  ELSE
    <1000 blocks> CHECKSEQUENCEVERIFY DROP PK_B
  ENDIF
  CHECKSIG
```

#### bilateral revocation protocol

in each round, as the channel state is advanced, the two parties

- exchange new commitments
- exchange revocation keys for the previous commitment
- sign each other's commitment transactions

suppose A is to pay B

- for B, he would benefit, i.e. the upcoming state would grant him more balance, so he wouldn't broadcast the previous commitment tx
- for A, her ability to cheat by broadcasting previous commitment has been revoked, since if she does that, A could redeem the exact tx outputs with the full signature (signed by `PK_A` and `RPK_B`)

> the revocation doesn't happen automatically. The payee should watch the blockchain, and execute the revocation protocol within the delay specified in the payer's tx (i.e., 1000 blocks in our cases) which would weaken the revocability

- Asymmetric revocable commitments with relative time locks (CSV) are a much better way to implement payment channels and a very significant innovation in this technology

### Hash Time Lock Contracts (HTLC)

- **definition**: a special type of smart contract that allows the participants to commit funds to **a redeemable secret**, with **an expiration time**
- **hash** part: the intended recipent create a secret `R`, whose hash `H=Hash(R)` can be included in an output's locking script, and this output can be redeemed by everyone knows `R`
- **time lock** part: payer will be refunded in case of no secret be revealed before the expiration of the time lock (achieved by `CHECKLOCKTIMEVERIFY`)
- a naive example HTLC script

```
IF
  # Payment if you have the secret R
  HASH160 <H> EQUALVERIFY
ELSE
  # Refund after timeout.
  <locktime> CHECKLOCKTIMEVERIFY DROP
  <Payee Pubic Key> CHECKSIG
ENDIF
```

a practical solution would be adding a CHECKSIG operator and a public key in the
first clause restricts redemption of the hash to a named recipient, who must also
know the secret R.

## Routed Payment Channels (Lightning Network)

**motivation**

- payer isn't connected to payee by a payment channel
- a new payment channel would require a funding tx locked in blockchain

### Basic Lightning Network Example

![lightning network](./images/lightning-network.png)

**notices**

- payment is decremented hop by hop
- time locks is alose decremented hop by hop

### Lightning Network Transport and Routing

- a long-term public key that they use as an identifier and to authenticate each other.
- every payment requires payer to construct a path through the network by connecting payment channels with **sufficient capacity**.
- the routing path is only known to payer. All other participants in the payment route see only the adjacent nodes
  - ensuring privacy of payments
  - making it very difficult to apply surveillance, censorship, or blacklists
- routing the path is always fixed at 20 hops and padded with random data which can signal the payee to stop secretly

### Lightning Network Benefits

- privacy
- fungibility
- speed
- granularity
- capacity
- truestless operation
