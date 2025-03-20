from enum import Enum, unique

from utility.utils import is_windows_platform
from utility.build_binary import __build_from_github

ZGS_BINARY = "zgs_node.exe" if is_windows_platform() else "zgs_node"
KV_BINARY = "zgs_kv.exe" if is_windows_platform() else "zgs_kv"

ZG_GIT_REV = "7bc25a060fab9c17bc9942b6747cd07a668d3042"  # v0.1.0


@unique
class BuildBinaryResult(Enum):
    AlreadyExists = 0
    Installed = 1
    NotInstalled = 2


def build_zgs(dir: str) -> BuildBinaryResult:
    return __build_from_github(
        dir=dir,
        binary_name=ZGS_BINARY,
        github_url="https://github.com/0glabs/0g-storage-node.git",
        build_cmd="cargo build --release",
        compiled_relative_path=["target", "release"],
    )


def build_kv(dir: str) -> BuildBinaryResult:
    return __build_from_github(
        dir=dir,
        binary_name=KV_BINARY,
        github_url="https://github.com/0glabs/0g-storage-kv.git",
        build_cmd="git submodule update --init --recursive && cargo build --release",
        compiled_relative_path=["target", "release"],
    )
