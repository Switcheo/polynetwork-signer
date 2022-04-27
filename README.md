<h1> polynetwork signer </h1>

Signing helper for hydrogen to help with getting the correct signed transaction to be broadcasted to poly network.

Created as a stopgap solution because the relayer code was all written in golang and its non-trivial to port them all into js in a short amount of time  

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

Before running, you need to have the configuration file `config.yml`
the path to the file should be set in a .env file in the root
`CONFIG_FILE_PATH=<your_file_path>`

```
rpcEndPoints:
  polynetwork: 'http://seed1.poly.network:20336'
broadcaster:
  polynetwork:
    fileName: './wallet.dat'
    password: '<yourpassword>'
```

## Run

Examples:

```shell
# For a cross chain tx from evm
./hydrogen_polynetwork_signer create_crosschain_tx <chain_id> <tx_data> <height> <proof>

# For a cross chain tx from carbon (will automatically get header and proof from provided rpc url)
./hydrogen_polynetwork_signer create_crosschain_tx_tendermint <rpc_url> <chain_id> <height> <ccmKeyHash>
```

```shell
./hydrogen_polynetwork_signer create_crosschain_tx 2 3e 12344 3333FF

./hydrogen_polynetwork_signer create_crosschain_tx_tendermint http://52.76.86.86:26657 2 12344 12333
```
