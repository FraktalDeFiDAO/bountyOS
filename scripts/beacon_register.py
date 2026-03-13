#!/usr/bin/env python3
import json
import urllib.request
import ssl
import random
import os
from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey

tmp_dir = "/home/administrator/.gemini/tmp/bountyos"
os.makedirs(tmp_dir, exist_ok=True)

private_key = Ed25519PrivateKey.generate()
public_key = private_key.public_key()

public_hex = public_key.public_bytes_raw().hex()
private_hex = private_key.private_bytes_raw().hex()

agent_id = f"bcn_fraktal_{random.randint(1000,9999)}"

with open(f"{tmp_dir}/beacon_keys.json", "w") as f:
    json.dump({
        "agent_id": agent_id,
        "public_key": public_hex,
        "private_key": private_hex,
        "agent_name": "FraktalDeFiDAO-Worker",
        "model": "qwen-code"
    }, f, indent=2)

ctx = ssl.create_default_context()
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE

registration_data = {
    "agent_id": agent_id,
    "pubkey_hex": public_hex,
    "model_id": "qwen-code",
    "provider": "other",
    "capabilities": ["coding", "research"],
    "name": "FraktalDeFiDAO-Worker"
}

req = urllib.request.Request(
    "https://rustchain.org/beacon/relay/register",
    data=json.dumps(registration_data).encode('utf-8'),
    headers={"Content-Type": "application/json"},
    method="POST"
)

try:
    with urllib.request.urlopen(req, context=ctx, timeout=30) as response:
        result = json.loads(response.read().decode('utf-8'))
        print(f"=== Registration SUCCESS ===")
        print(json.dumps(result, indent=2))
        
        # Save the real agent_id and token from API
        if "agent_id" in result:
            agent_id = result["agent_id"]
            
        with open(f"{tmp_dir}/beacon_keys.json", "w") as f:
            json.dump({
                "agent_id": agent_id,
                "public_key": public_hex,
                "private_key": private_hex,
                "agent_name": "FraktalDeFiDAO-Worker",
                "model": "qwen-code"
            }, f, indent=2)
            
        if "relay_token" in result:
            with open(f"{tmp_dir}/beacon_token.txt", "w") as f:
                f.write(result["relay_token"])
                
except urllib.error.HTTPError as e:
    print(f"Status: {e.code}")
    print(f"Response: {e.read().decode('utf-8')}")
