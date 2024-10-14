from utility.utils import PortMin, MAX_NODES

def kv_rpc_port(n):
    return PortMin.n + 5 * MAX_NODES + n

def indexer_port(n):
    return PortMin.n + 6 * MAX_NODES + n