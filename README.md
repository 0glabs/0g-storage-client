# 0G Storage Client
Go implementation for client to interact with storage nodes in 0G Storage network.

# SDK

Following packages can help applications to integrate with 0g storage network:

**[transfer](transfer)** : defines data structures and functions for transferring data between local and 0g storage.
**[kv]**: defines structures to interact with 0g storage kv.
**[indexer]**: select storage nodes to upload data from indexer which maintains trusted node list.
**[node]**: defines RPC client structures to facilitate RPC interactions with 0g storage nodes and 0g key-value (KV) nodes.

# CLI
Run `go build` under the root folder to compile the executable binary.

**Global options**
```
      --gas-limit uint     Custom gas limit to send transaction
      --gas-price uint     Custom gas price to send transaction
  -h, --help               help for 0g-storage-client
      --log-force-color    Force to output colorful logs
      --log-level string   Log level (default "info")
```

**Generate test file**

To generate a file for test purpose, especially with a fixed file size or random file size (without `--size` option):

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