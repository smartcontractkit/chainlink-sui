#!/usr/bin/env bash
set -euo pipefail

# Generate new ed25519 keypair file (<ADDRESS>.key)
sui keytool generate ed25519 >/dev/null

# Find the most recently created key file
F=$(ls -t *.key | head -1)

# Extract private key (hex, 32 bytes, with 0x prefix)
HEX_PRIV=$(base64 -d "$F" | tail -c +2 | xxd -p -c 64)

# Get the address/public key associated with this new key
ADDR=$(basename "$F" .key)

echo "Address: $ADDR"
echo "Private key (hex): 0x$HEX_PRIV"

rm -f "$F"
