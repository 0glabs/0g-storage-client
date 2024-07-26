#!/usr/bin/env python3

import os

from utility.run_all import run_all

if __name__ == "__main__":
    run_all(
        test_dir = os.path.dirname(__file__),
        slow_tests={},
        long_manual_tests={},
    )