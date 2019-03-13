# Chapter 07. Advanced Transactions and Scripting

## Timelocks

- Timelocks are restrictions on transactions or outputs that only allow spending after a point in time
- Timelocks are useful for postdating transactions and locking funds to a date in the future

### Transaction Locktime (nLocktime)

- Transaction locktime is a transaction-level setting (`nLocktime` field in the transaction data structure) that defines the earliest time that a transaction is valid and can be relayed on the network or added to the blockchain
- `0<nLocktime<500,000,000` is interpreted as a block height
- `500,000,000<=nLocktime` is interpreted as a Unix Epoch timestamp (**seconds** since Jan-1-1970)

#### Transaction locktime limitations

- `nLocktime` has the limitation that while it makes it possible to spend some outputs in the future, it does not make it impossible to spend them earlier than that time
- Example: A pays B by means of a tx with `nLocktime` being 3 months
  - Then
    - B cannot transmit the transaction (since it won't be packaged into a mined block) to redeem the funds until 3 months have elapsed
    - B may transmit the transaction after 3 months
  - However
    - A can double spend the same input of the tx with `nLocktime=0`
    - B has no guarantee that A won't do double spending

### Check Lock Time Verify (CLTV)

- `CLTV` is a per-**output** timelock
- CLTV doesn't replace `nLocktime`, but rather restricts specific **UTXO** such that they can only be spent in a future transaction with `nLocktime` set to a greater or equal value
  > **UTXO** indicates the funding tx producing this UTXO has been mined
- The `CLTV` opcode takes one parameter as input, expressed as a number in **the same format** as `nLocktime` (either a block height or Unix epoch time).
- As indicated by the `VERIFY` suffix, `CLTV` is the type of opcode that halts execution of the script if the outcome is `FALSE`
- To spend a UTXO locked by `CLTV`, B constructs a transaction that references the UTXO as an input as follows
  - Uses his signature and public key in the unlocking script of that input
  - Sets the transaction `nLocktime` to be equal or greater to the timelock in the `CHECKLOCKTIMEVERIFY` A set
  - Broadcasts the transaction on the bitcoin network
- After execution, if `CLTV` is satisfied, the time parameter that preceded it remains as the top item on the stack and may need to be dropped, with `DROP`, for correct execution of subsequent script opcodes

### Relative Timelocks

- **WHY**: `nLocktime` and `CLTV` are both **absolute timelocks** in that they specify an absolute point in time
- **HOW**: Relative timelocks doesn't start counting until the funding UTXO is recor‐ ded on the blockchain
  - The transaction-level relative timelock is implemented as a consensus rule on the value of 4-byte `nSequence` field set in every transaction input
  - Script-level relative timelocks are implemented with the `CHECKSEQUENCEVERIFY` (CSV) opcode

### Relative Timelocks with `nSequence`

#### Original meaning of `nSequence`

- `nSequence<(2^32-1)` indicates a non-finalized tx, otherwise a finalized and mined tx
- `nSequence` is customarily set to `(2^32-1)` in transactions that do not utilize timelocks
- For transactions with `nLocktime` or `CHECKLOCKTIMEVERIFY`, the `nSequence` value must be set to less than `(2^32-1)` for the timelock guards to have effect (customarily as `0xFFFF,FFFE`)

#### `nSequence` as a consensus-enforced relative timelock

- If the most significant (bit `(1<<31)`) is not set, it is a flag that means "relative locktime"
- A transaction with `nSequence<2^31` is only valid once the input has aged by the relative timelock amount
- A transaction can include both timelocked inputs (`nSequence<2^31`) and inputs without a relative timelock (`nSequence>=2^31`)
- The `nSequence` value is specified in either blocks or seconds
  - A type-flag set in the 23rd LSB (value `1<<22`) where
    - `1` would cause `nSequence` to be interpreted as a multiple of 512
    - `0` would cause `nSequence` to be interpreted as a number of blocks
  - When interpreting `nSequence` as a relative timelock, only the `16` LSB are considered

### Relative Timelocks with CSV (`CHECKSEQUENCEVERIFY`)

- As with `CLTV`, the value in `CSV` must match the format in the corresponding `nSequence` value
- Relative timelocks with `CSV` are especially useful when several (chained) transactions are created and signed, but not propagated, when they're kept "off-chain"
- Applications
  - Payment channels and state channels
  - Lightning network

### Median-Time-Past

- There is a subtle, but very significant, difference between wall time and consensus time
- The timestamps are set by miners, enabling them to lie about the time in a block so as to earn extra fees by including timelocked transactions that are not yet mature
- **DEFINITION**: Median-Time-Past is calculated by taking the timestamps of the last 11 blocks and finding the median
- Median-Time-Past is consensus time and used for all timelock calculations
- The consensus time calculated by Median-Time-Past is always approximately one hour (6 blocks) behind wall clock time

### Timelock Defense Against Fee Sniping

- Fee Sniping: Miners attempting to rewrite past blocks "snipe" higher-fee transactions from future blocks to maximize their profitability
- To prevent "fee sniping," when Bitcoin Core creates transactions, it uses `nLocktime` to limit them to the "next block," by default
- Official comments go as [Discourage fee sniping](https://github.com/bitcoin/bitcoin/commit/db6047d61b742be07442f891e70570b791c585e3)
