#!/usr/bin/env python3

import base64
import random
import tempfile
from web3 import Web3, HTTPProvider

from config.node_config import GENESIS_ACCOUNT
from utility.submission import ENTRY_SIZE, bytes_to_entries
from utility.utils import (
    assert_equal,
    wait_until,
)
from test_framework.test_framework import TestFramework


class SkipTxTest(TestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 2

    def run_test(self):
        node_idx = random.randint(0, self.num_nodes - 1)

        file_to_upload = tempfile.NamedTemporaryFile(dir=self.root_dir, delete=False)
        data = random.randbytes(256 * 2048) 
        file_to_upload.write(data)
        file_to_upload.close()
        w3 = Web3(HTTPProvider(self.blockchain_nodes[0].rpc_url))

        nonce = w3.eth.get_transaction_count(GENESIS_ACCOUNT.address)
        # first submission
        root = self._upload_file_use_cli(
            self.blockchain_nodes[0].rpc_url,
            self.contract.address(),
            GENESIS_ACCOUNT.key,
            self.nodes[node_idx].rpc_url,
            None,
            file_to_upload,
        )

        self.log.info("root: %s", root)
        wait_until(lambda: self.contract.num_submissions() == 1)

        client = self.nodes[node_idx]
        wait_until(lambda: client.zgs_get_file_info(root) is not None)
        wait_until(lambda: client.zgs_get_file_info(root)["finalized"])
        
        self._download_file_use_cli(self.nodes[node_idx].rpc_url, None, root)
        assert_equal(w3.eth.get_transaction_count(GENESIS_ACCOUNT.address), nonce + 1)

        nonce = w3.eth.get_transaction_count(GENESIS_ACCOUNT.address)
        # second submission
        root = self._upload_file_use_cli(
            self.blockchain_nodes[0].rpc_url,
            self.contract.address(),
            GENESIS_ACCOUNT.key,
            self.nodes[node_idx].rpc_url,
            None,
            file_to_upload,
            True
        )
        wait_until(lambda: self.contract.num_submissions() == 1)
        assert_equal(w3.eth.get_transaction_count(GENESIS_ACCOUNT.address), nonce)

if __name__ == "__main__":#
    SkipTxTest().main()
