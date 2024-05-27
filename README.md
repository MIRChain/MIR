## **MIR** blockchain

**Inspiration:** _Mir (Russian: Мир, IPA: [ˈmʲir]; lit. 'peace' or 'world') was the first **modular** space station and was assembled in orbit from 1986 to 1996._

### Features

- Different types of crypto signature can be chosen:
  - GOST 34.10 (any 256 bit curve)
  - CyptoProGOST (https://cryptopro.ru/en)
  - NIST Secp256k1
- Different type of hash function can be chosen:
  - SHA3
  - SteebogHash
- Different consensus algorithms are also available:
  - Proof-of-Work
  - Proof-of-Authority
  - Raft
  - iBFT
  - QBFT
- High transaction throughput at Proof-of-Authority/Raft/iBFT/QBFT consensus
- All of the EVM and Ethereum tools working out of the box

Documentation can be found here: [https://mirchain.github.io/MIR/](https://mirchain.github.io/MIR/)

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
