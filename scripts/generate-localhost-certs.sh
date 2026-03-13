#!/usr/bin/env bash
set -euo pipefail

cert_dir="${1:-certs}"
crt_path="${cert_dir}/localhost.crt"
key_path="${cert_dir}/localhost.key"

mkdir -p "${cert_dir}"

openssl req \
  -x509 \
  -newkey rsa:4096 \
  -sha256 \
  -days 825 \
  -nodes \
  -keyout "${key_path}" \
  -out "${crt_path}" \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1,IP:::1"

chmod 600 "${key_path}"
chmod 644 "${crt_path}"

echo "Generated:"
echo "  ${crt_path}"
echo "  ${key_path}"
