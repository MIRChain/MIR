## MIR chain

#### Implementation of Quorum/Ethereum with Russian GOST (and Post-Quantum) cryptography under the hood.

### Features

- Different types of crypto signature can be chosen at a new chain initialization: 
  - GOST 34.10 (any 256 bit curve)
  - CyptoProGOST
  - NIST PostQuantum
  - NIST Secp256k1
- Different type of hash function can be chosen at a new chain initialization:
  - SHA3
  - SteebogHash
- Differetnt consensus algoruthms are also avaliable at a new chain initialization:
  - Proof-of-Work
  - Proof-of-Authority 
  - Raft
  - iBFT
  - QBFT
- High transaction throughput at Proof-of-Authority/Raft/iBFT/QBFT consensus
- All of the EVM and Ethereum tools working out of the box

## Building the source

Building `mir` requires both a Go (version 1.18 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make mir
```

or, to build the full suite of utilities:

```shell
make all
```
