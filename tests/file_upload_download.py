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
from test_framework.test_framework import TestFramework


class FileUploadDownloadTest(TestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 2

    def run_test(self):
        data_size = [
            2,
            255,
            256,
            257,
            1023,
            1024,
            1025,
            256 * 1023,
            256 * 1024,
            256 * 1025,
            256 * 2048,
            256 * 16385,
            256 * 1024 * 64,
        ]

        for i, v in enumerate(data_size):
            self.__test_upload_download_file(
                v, i + 1, False if v >= 256 * 1024 * 64 else True
            )

    def __test_upload_download_file(self, size, submission_index, rand_data=True):
        node_idx = random.randint(0, self.num_nodes - 1)
        self.log.info("node index: %d, file size: %d", node_idx, size)

        file_to_upload = tempfile.NamedTemporaryFile(dir=self.root_dir, delete=False)
        data = random.randbytes(size) if rand_data else b"\x10" * size

        file_to_upload.write(data)
        file_to_upload.close()

        root = self._upload_file_use_cli(
            self.blockchain_nodes[0].rpc_url,
            self.contract.address(),
            GENESIS_ACCOUNT.key,
            self.nodes[node_idx].rpc_url,
            None,
            file_to_upload,
        )

        self.log.info("root: %s", root)
        wait_until(lambda: self.contract.num_submissions() == submission_index)

        client = self.nodes[node_idx]
        wait_until(lambda: client.zgs_get_file_info(root) is not None)
        wait_until(lambda: client.zgs_get_file_info(root)["finalized"])
        
        self._download_file_use_cli(self.nodes[node_idx].rpc_url, None, root)

if __name__ == "__main__":
    FileUploadDownloadTest().main()
