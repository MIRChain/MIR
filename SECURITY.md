# Security Policy

## Supported Versions

Please see [Releases](https://github.com/pavelkrolevets/MIR-pro/releases). We recommend using the [most recently released version](https://github.com/pavelkrolevets/MIR-pro/releases/latest).

## Audit reports

Audit reports are published in the `docs` folder: https://github.com/pavelkrolevets/MIR-pro/tree/master/docs/audits 

| Scope | Date | Report Link |
| ------- | ------- | ----------- |
| `geth` | 20170425 | [pdf](https://github.com/pavelkrolevets/MIR-pro/blob/master/docs/audits/2017-04-25_Geth-audit_Truesec.pdf) |
| `clef` | 20180914 | [pdf](https://github.com/pavelkrolevets/MIR-pro/blob/master/docs/audits/2018-09-14_Clef-audit_NCC.pdf) |

## Reporting a Vulnerability

**Please do not file a public ticket** mentioning the vulnerability.

To find out how to disclose a vulnerability in Ethereum visit [https://bounty.ethereum.org](https://bounty.ethereum.org) or email bounty@ethereum.org. Please read the [disclosure page](https://github.com/pavelkrolevets/MIR-pro/security/advisories?state=published) for more information about publically disclosed security vulnerabilities.

Use the built-in `geth version-check` feature to check whether the software is affected by any known vulnerability. This command will fetch the latest [`vulnerabilities.json`](https://geth.ethereum.org/docs/vulnerabilities/vulnerabilities.json) file which contains known security vulnerabilities concerning `geth`, and cross-check the data against its own version number.

The following key may be used to communicate sensitive information to developers.

Fingerprint: `AE96 ED96 9E47 9B00 84F3 E17F E88D 3334 FA5F 6A0A`

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

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
=r6KK
-----END PGP PUBLIC KEY BLOCK-----
```
