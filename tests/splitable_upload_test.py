#!/usr/bin/env python3

import random
import tempfile
import os

from config.node_config import GENESIS_ACCOUNT
from utility.utils import (
    wait_until,
)
from client_test_framework.test_framework import ClientTestFramework

def files_are_equal(file1, file2):
    if os.path.getsize(file1) != os.path.getsize(file2):
        return False

    with open(file1, 'rb') as f1, open(file2, 'rb') as f2:
        while True:
            chunk1 = f1.read(4096)
            chunk2 = f2.read(4096)
            
            if chunk1 != chunk2:
                return False
            
            if not chunk1:  
                break

    return True 


class FileUploadDownloadTest(ClientTestFramework):
    def setup_params(self):
        self.num_blockchain_nodes = 1
        self.num_nodes = 4
        self.zgs_node_configs[0] = {
            "db_max_num_sectors": 2 ** 30,
            "shard_position": "0/4"
        }
        self.zgs_node_configs[1] = {
            "db_max_num_sectors": 2 ** 30,
            "shard_position": "1/4"
        }
        self.zgs_node_configs[2] = {
            "db_max_num_sectors": 2 ** 30,
            "shard_position": "2/4"
        }
        self.zgs_node_configs[3] = {
            "db_max_num_sectors": 2 ** 30,
            "shard_position": "3/4"
        }

    def run_test(self):
        size = 50 * 1024 * 1024 # 50M
        file_to_upload = tempfile.NamedTemporaryFile(dir=self.root_dir, delete=False)
        data = random.randbytes(size)

        file_to_upload.write(data)
        file_to_upload.close()

        roots = self._upload_file_use_cli(
            self.blockchain_nodes[0].rpc_url,
            GENESIS_ACCOUNT.key,
            ','.join([x.rpc_url for x in self.nodes]),
            None,
            file_to_upload,
            fragment_size=1024*1024*3, # 3M, will aligned to 4M
        )

        self.log.info("roots: %s", roots)
        wait_until(lambda: self.contract.num_submissions() == 13)
        
        file_to_download = os.path.join(self.root_dir, "downloaded")
        self._download_file_use_cli(','.join([x.rpc_url for x in self.nodes]), None, roots=roots, with_proof=True, file_to_download=file_to_download, remove=False)
        assert(files_are_equal(file_to_upload.name, file_to_download))

if __name__ == "__main__":
    FileUploadDownloadTest().main()
