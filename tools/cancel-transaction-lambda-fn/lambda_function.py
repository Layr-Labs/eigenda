import logging
import os
import requests

from lambda_helper import (
    assemble_tx,
    get_params,
    get_tx_params,
    calc_eth_address,
    get_kms_public_key
)

LOG_LEVEL = os.getenv("LOG_LEVEL", "WARNING")
LOG_FORMAT = "%(levelname)s:%(lineno)s:%(message)s"
handler = logging.StreamHandler()

_logger = logging.getLogger()
_logger.setLevel(LOG_LEVEL)


def get_dynamic_gas_fees():
    """
    Dynamically fetch gas fee parameters from the ETH RPC.
    This function calls 'eth_gasPrice' to get the base fee.
    The max priority fee is read from the environment variable 'PRIORITY_FEE_GWEI'
    (submitted in gwei) and converted to wei. The max fee per gas is then calculated
    as the sum of the base fee and the max priority fee.
    """
    eth_rpc = os.getenv("ETH_RPC")
    if not eth_rpc:
        raise ValueError("ETH_RPC environment variable not set")
    
    # Fetch the current base gas price from the ETH RPC
    headers = {"Content-Type": "application/json"}
    data = {
        "jsonrpc": "2.0",
        "method": "eth_gasPrice",
        "params": [],
        "id": 1,
    }
    response = requests.post(eth_rpc, json=data, headers=headers)
    result = response.json()
    gas_price_hex = result.get("result")
    if not gas_price_hex:
        raise Exception("Failed to fetch gas price from RPC")
    
    base_fee = int(gas_price_hex, 16)
    
    # Get the priority fee from the environment variable in gwei and convert to wei
    priority_fee_gwei = os.getenv("PRIORITY_FEE_GWEI")
    if priority_fee_gwei is None:
        raise ValueError("PRIORITY_FEE_GWEI environment variable not set")
    try:
        priority_fee_int = int(priority_fee_gwei)
    except ValueError:
        raise ValueError("PRIORITY_FEE_GWEI environment variable must be an integer")
    
    if priority_fee_int < 0:
        raise ValueError("PRIORITY_FEE_GWEI environment variable must be greater than or equal to 0")
    
    max_priority_fee = priority_fee_int * 10**9
    
    max_fee = base_fee + max_priority_fee
    return max_fee, max_priority_fee


def lambda_handler(event, context):
    _logger.debug("incoming event: {}".format(event))

    # Load common transaction parameters
    try:
        params = get_params()
    except Exception as e:
        raise e

    # Get the nonce to be used (from environment)
    nonce_env = os.getenv("NONCE")
    if nonce_env is None:
        raise ValueError("NONCE environment variable not set")
    try:
        nonce = int(nonce_env)
    except ValueError:
        raise ValueError("NONCE environment variable must be an integer")
    
    if nonce < 0:
        raise ValueError("NONCE environment variable must be greater than or equal to 0")

    # Get the KMS key ID from the environment
    key_id = os.getenv("KMS_KEY_ID")
    if not key_id:
        raise ValueError("KMS_KEY_ID environment variable not set")

    # Dynamically fetch gas fees from the ETH RPC
    max_fee_per_gas, max_priority_fee_per_gas = get_dynamic_gas_fees()

    # Retrieve the public key from KMS and calculate the Ethereum address
    pub_key = get_kms_public_key(key_id)
    eth_checksum_addr = calc_eth_address(pub_key)

    # For a cancellation transaction, send 0 ETH to your own address
    dst_address = eth_checksum_addr
    amount = 0

    # Set the chain ID to Ethereum mainnet (chain id 1)
    chainid = 1

    # Set the transaction type (EIP-1559 transaction type 2)
    tx_type = 2

    # Build the transaction parameters
    tx_params = get_tx_params(
        dst_address=dst_address,
        amount=amount,
        nonce=nonce,
        chainid=chainid,
        type=tx_type,
        max_fee_per_gas=max_fee_per_gas,
        max_priority_fee_per_gas=max_priority_fee_per_gas
    )

    # Assemble and sign the Ethereum transaction offline
    raw_tx_signed_hash, raw_tx_signed_payload = assemble_tx(
        tx_params=tx_params,
        params=params,
        eth_checksum_addr=eth_checksum_addr,
        chainid=chainid
    )

    # Send the signed transaction using the ETH RPC endpoint
    eth_rpc = os.getenv("ETH_RPC")
    if not eth_rpc:
        raise ValueError("ETH_RPC environment variable not set")
    headers = {"Content-Type": "application/json"}
    data = {
        "jsonrpc": "2.0",
        "method": "eth_sendRawTransaction",
        "params": [raw_tx_signed_payload],
        "id": 1,
    }
    response = requests.post(eth_rpc, json=data, headers=headers)
    rpc_response = response.json()

    return {
        "signed_tx_hash": raw_tx_signed_hash,
        "rpc_response": rpc_response,
    }
