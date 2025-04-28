Operators run script that reports time offset to EigenDA 
1. Check clock drift (NTP synchronization).
2. Sign the collected data (node_id, timestamp, offset) locally with their Ethereum operator wallet private key.
3. Send the signed message to EigenDA disperser’s API endpoint.
4. EigenDA disperser server’s endpoint
    1. validates signatures on incoming reports using the operator's registered wallet addresses (retrieved from the on-chain registry).
    2. record metrics 
    3. alert if the offset is greater than the configured threshold


Acceptance

1. Define an acceptance bound for latency/drift; currently it is set to reservation interval for reservation payments, and much more strict for global on-demand payments (no cusioning when processing on-demand at the edge of global on-demand ratelimit interval)
    1. if the receive time for a given period from majority of nodes are within $\delta$, then clock synchrony is acceptable
2. Dispatcher makes statistical analysis over multiple measurements (median, mean, variance, …)
3. DataAPI exposes a new endpoint `operators/clock-sync` to report clock sync across all nodes



Idea 4: Script reporting time sync to EigenDA

```protobuf

// The payload submitted by an operator reporting their clock drift
message ClockDriftPayload {
  string operator_address = 1;  // Ethereum address of the operator
  string timestamp = 2;         // ISO 8601 UTC timestamp (RFC3339 format)
  double offset_seconds = 3;    // Clock drift in seconds
}

// The full request sent to the server, including the payload and the signature
message ClockDriftRequest {
  ClockDriftPayload payload = 1; // Payload containing drift data
  string signature = 2;          // Ethereum signed message (hex encoded)
}

```

(Script in Python)

Prerequisites:

```bash
pip install web3 requests
```

Script

```bash
import subprocess
import requests
from web3 import Web3
from eth_account import Account
from eth_account.messages import encode_defunct
from datetime import datetime

# === Operator Configuration ===
NTP_SERVER = "pool.ntp.org"
EIGENDA_API_ENDPOINT = "disperser address new api /clock_drift"
PRIVATE_KEY = "0xYourOperatorPrivateKey"  # Keep securely, e.g., env variable

# Derive operator address from private key
acct = Account.from_key(PRIVATE_KEY)
operator_address = acct.address

# Function to measure clock offset (sntp command required)
def get_offset():
    try:
        result = subprocess.check_output(['sntp', NTP_SERVER]).decode().split()[0]
        return float(result)
    except Exception as e:
        print(f"Error fetching NTP offset: {e}")
        return None

# Sign message using operator's private key
def sign_message(private_key, message):
    msg = encode_defunct(text=message)
    signed_message = Account.sign_message(msg, private_key=private_key)
    return signed_message.signature.hex()

# Make this a routine to run periodically for a couple minutes
offset = get_offset()
if offset is None:
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

response = requests.post(EIGENDA_API_ENDPOINT, json=data)
print(f"[{timestamp}] Offset reported: {offset}s | HTTP {response.status_code} | {response.text}")

```

Corresponding EigenDA disperser Server

```go

// Handler function for clock drift reporting
func ClockDriftHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var driftReq ClockDriftRequest
	if err := json.Unmarshal(body, &driftReq); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	payload := driftReq.Payload
	signature := driftReq.Signature

	// Verify operator is registered from OperatorState
	

	// Verify signature

	// Verify timestamp freshness (within a reservation interval)
	reportTime, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
		return
	}
	if time.Since(reportTime).Abs() > RESERVATION_INTERVAL {
		http.Error(w, "Timestamp expired", http.StatusBadRequest)
		return
	}

	log.Printf("[Clock Drift] Operator: %s, Offset: %f sec at %s\n",
		payload.OperatorAddress, payload.OffsetSeconds, payload.Timestamp)
	// Report to grafana, collect metrics

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

```
