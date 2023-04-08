## MIR chain PRO

#### Implementation of the MIR-protocol with Russian GOST (and Post-Quantum) cryptography under the hood. Based on Quorum/Ethereum.

### Features

- Different types of crypto can be signature chosen at the chain initialization: 
  - GOST 34.10 (any 256 bit curve)
  - CyptoProGOST
  - NIST PostQuantum
  - NIST Secp256k1
- Different type of hash function:
  - SHA3
  - SteebogHash
- Differetnt consensus algoruthms are avaliable:
  - Proof-of-Authority 
  - Raft
  - iBFT
  - QBFT
- High transaction throughput
- GUI for easy nodes management
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
