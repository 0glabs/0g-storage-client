#!/usr/bin/env python3

import base64
import random
import tempfile

from config.node_config import GENESIS_ACCOUNT
from utility.submission import ENTRY_SIZE, bytes_to_entries
from utility.utils import (
    assert_equal,
    wait_until,
)
from utility.kv import to_stream_id
from test_framework.test_framework import TestFramework


class KVTest(TestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 1

    def run_test(self):
        self.setup_kv_node(0, [to_stream_id(0)])
        self._kv_write_use_cli(
            self.blockchain_nodes[0].rpc_url,
            self.contract.address(),
            GENESIS_ACCOUNT.key,
            self.nodes[0].rpc_url,
            None,
            to_stream_id(0),
            "0,1,2,3,4,5,6,7,8,9,10",
            "0,1,2,3,4,5,6,7,8,9,10"
        )
        wait_until(lambda: self.kv_nodes[0].kv_get_trasanction_result(0) == "Commit")
        res = self._kv_read_use_cli(
            self.kv_nodes[0].rpc_url,
            to_stream_id(0),
            "0,1,2,3,4,5,6,7,8,9,10,11"
        )
        for i in range(11):
            assert_equal(res[str(i)], str(i))
        assert_equal(res["11"], "")

if __name__ == "__main__":
    KVTest().main()
