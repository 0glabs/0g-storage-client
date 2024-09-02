# 0G Storage Client

Go implementation for client to interact with storage nodes in 0G Storage network. For more details, please read the [docs](https://docs.0g.ai/0g-doc/docs).

## SDK

Following packages can help applications to integrate with 0g storage network:

- **[core](core)**: provides underlying utilities to build merkle tree for files or iteratable data, and defines data padding standard to interact with [Flow contract](contract/contract.go).
- **[node](node)**: defines RPC client structures to facilitate RPC interactions with 0g storage nodes and 0g key-value (KV) nodes.
- **[kv](kv)**: defines structures to interact with 0g storage kv.
- **[transfer](transfer)** : defines data structures and functions for transferring data between local and 0g storage.
- **[indexer](indexer)**: select storage nodes to upload data from indexer which maintains trusted node list. Besides, allow clients to download files via HTTP GET requests.

## CLI

Run `go build` under the root folder to compile the executable binary. There are several commands to interact with 0g storage node.

**Global flags**

Run `./0g-storage-client --help` to view all available commands along with global flags:
```
Flags:
      --gas-limit uint                Custom gas limit to send transaction
      --gas-price uint                Custom gas price to send transaction
  -h, --help                          help for 0g-storage-client
      --log-color-disabled            Force to disable colorful logs
      --log-level string              Log level (default "info")
      --rpc-retry-count int           Retry count for rpc request (default 5)
      --rpc-retry-interval duration   Retry interval for rpc request (default 5s)
      --rpc-timeout duration          Timeout for single rpc request (default 30s)
      --web3-log-enabled              Enable log for web3 RPC
```

**Generate test file**

To generate a file for test purpose, with a fixed file size or random file size (without `--size` option):

```
./0g-storage-client gen --size <file_size_in_bytes>
```

**Upload file**

```
./0g-storage-client upload --url <blockchain_rpc_endpoint> --contract <0g-storage_contract_address> --key <private_key> --indexer <storage_indexer_endpoint> --file <file_path>
```

The client will submit the data segments to the storage nodes which is determined by the indexer according to their shard configurations.

**Download file**
```
./0g-storage-client download --indexer <storage_indexer_endpoint> --root <file_root_hash> --file <output_file_path>
```

If you want to verify the **merkle proof** of downloaded segment, please specify `--proof` option.

**Write to KV**

By indexer:
```
./0g-storage-client kv-write --url <blockchain_rpc_endpoint> --contract <0g-storage_contract_address> --key <private_key> --indexer <storage_indexer_endpoint> --stream-id <stream_id> --stream-keys <stream_keys> --stream-values <stream_values>
```

`--stream-keys` and `--stream-values` are comma separated string list and their length must be equal.

**Read from KV**

```
./0g-storage-client kv-read --node <kv_node_rpc_endpoint> --stream-id <stream_id> --stream-keys <stream_keys>
```

Please pay attention here `--node` is the url of a KV node.

## Indexer

Indexer service provides RPC to index storages nodes in two ways:

- Trusted nodes: well maintained and provides stable service.
- Discovered nodes: discovered in the whole P2P network.

Please refer to the [RPC API](https://docs.0g.ai/0g-doc/docs/0g-storage/rpc/indexer-api) documentation for more details.

Besides, indexer supports to download files via HTTP GET request:

```
/file?txSeq=7
or
/file?root=0x0376e0d95e483b62d5100968ed17fe1b1d84f0bc5d9bda8000cdfd3f39a59927
```

Note, user could specify `name=foo.log` parameter in GET URL to rename the downloaded file.
