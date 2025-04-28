import os
import subprocess
import requests
from web3 import Web3
from eth_account import Account
from eth_account.messages import encode_defunct
from datetime import datetime

# === Operator Configuration ===
NTP_SERVER = os.environ.get("NTP_SERVER", "pool.ntp.org")
EIGENDA_API_ENDPOINT = os.environ.get("EIGENDA_API_ENDPOINT", "http://localhost:8080/clock_drift")
PRIVATE_KEY = os.environ.get("OPERATOR_PRIVATE_KEY")  # Must be set in environment

if not PRIVATE_KEY:
    raise ValueError("OPERATOR_PRIVATE_KEY environment variable not set!")

# Derive operator address from private key
acct = Account.from_key(PRIVATE_KEY)
operator_address = acct.address
print(f"Operator address: {operator_address}")

# Function to measure clock offset (sntp command required)
def get_offset():
    try:
        result = subprocess.check_output(['sntp', NTP_SERVER]).decode().strip()
        # Example output: '+0.064426 +/- 0.011359 pool.ntp.org 102.129.185.135'
        offset_str = result.split()[0]
        return float(offset_str)
    except Exception as e:
        print(f"Error fetching NTP offset: {e}")
        return None

# Sign message using operator's private key
def sign_message(private_key, message):
    msg = encode_defunct(text=message)
    signed_message = Account.sign_message(msg, private_key=private_key)
    return signed_message.signature.hex()

# Main routine
def main():
    offset = get_offset()
    if offset is None:
        print("Failed to get NTP offset. Exiting.")
        return

    timestamp = datetime.utcnow().isoformat() + "Z"

    payload = {
        "operator_address": operator_address,
        "timestamp": timestamp,
        "offset_seconds": offset
    }

    message_to_sign = f"{operator_address}|{timestamp}|{offset}"
    signature = sign_message(PRIVATE_KEY, message_to_sign)

    data = {
        "payload": payload,
        "signature": signature
    }

    print(f"Reporting offset: {offset}s at {timestamp}")
    try:
        response = requests.post(EIGENDA_API_ENDPOINT, json=data)
        print(f"[{timestamp}] Offset reported: {offset}s | HTTP {response.status_code} | {response.text}")
    except Exception as e:
        print(f"Error reporting to EigenDA API: {e}")

if __name__ == "__main__":
    main() 