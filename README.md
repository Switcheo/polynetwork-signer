<h1> polynetwork signer </h1>

Signing helper for hydrogen

## Build From Source

### Prerequisites

- [Golang](https://golang.org/doc/install) version 1.14 or later

### Build

```shell
go build -o hydrogen_polynetwork_signer main.go
```

After building the source code successfully, you should see the executable program `hydrogen_polynetwork_signer`.

## Setup

Before you can use the signer you will need to create a wallet file of PolyNetwork. After creation, you need to register
it as a Relayer to Poly net and get consensus nodes approving your registration. You can then send transactions to Poly
network and start relaying.

Before running, you need to have the configuration file `config.json` at the root

```
rpcEndPoints:
  polynetwork: 'http://seed1.poly.network:20336'
broadcaster:
  polynetwork:
    fileName: './wallet.dat'
    password: '<yourpassword>'
```

## Run

Example:


```shell
./hydrogen_polynetwork_signer <chain_id> <tx_data> <height> <proof>
```

```shell
./hydrogen_polynetwork_signer 2 0x3e 12344 0x3333
```
