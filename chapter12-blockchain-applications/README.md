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
