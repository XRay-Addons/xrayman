#!/bin/sh

# Generate an X25519 key pair using Xray and replace placeholders
# with the generated PrivateKey and PublicKey in the specified config files.
#
# Usage:
#   generate-x25519.sh <xray> <config1:config2:...> <private-placeholder> <public-placeholder>

set -eu

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <xray> <config1:config2:...> <private-placeholder> <public-placeholder>" >&2
    exit 1
fi

XRAY="$1"
CONFIGS="$2"
PRIVATE_PLACEHOLDER="$3"
PUBLIC_PLACEHOLDER="$4"

echo "Generate X25519 key pair..."
echo "Xray: $XRAY"

if [ ! -x "$XRAY" ]; then
    echo "Error: Xray executable not found or not executable: $XRAY" >&2
    exit 1
fi

echo "Running: $XRAY x25519"

if ! OUTPUT="$("$XRAY" x25519 2>&1)"; then
    echo "Error: xray x25519 failed." >&2
    echo "Xray output:" >&2
    printf '%s\n' "$OUTPUT" >&2
    exit 1
fi

echo "Xray output:"
printf '%s\n' "$OUTPUT"

PRIVATE_KEY=$(printf '%s\n' "$OUTPUT" |
    awk -F': ' '$1 == "PrivateKey" { print $2 }')

PUBLIC_KEY=$(printf '%s\n' "$OUTPUT" |
    awk -F': ' '$1 == "Password (PublicKey)" { print $2 }')

if [ -z "$PRIVATE_KEY" ]; then
    echo "Error: failed to extract PrivateKey from xray output." >&2
    exit 1
fi

if [ -z "$PUBLIC_KEY" ]; then
    echo "Error: failed to extract PublicKey from xray output." >&2
    exit 1
fi

echo "PrivateKey: $PRIVATE_KEY"
echo "PublicKey:  $PUBLIC_KEY"

# Replace placeholders in all config files.
OLDIFS=$IFS
IFS=:

for config in $CONFIGS; do
    IFS=$OLDIFS

    if [ ! -f "$config" ]; then
        echo "Error: config file not found: $config" >&2
        exit 1
    fi

    echo "Updating: $config"

    PRIVATE_PLACEHOLDER="$PRIVATE_PLACEHOLDER" \
    PUBLIC_PLACEHOLDER="$PUBLIC_PLACEHOLDER" \
    PRIVATE_KEY="$PRIVATE_KEY" \
    PUBLIC_KEY="$PUBLIC_KEY" \
    perl -0pi -e '
        s/\Q$ENV{PRIVATE_PLACEHOLDER}\E/$ENV{PRIVATE_KEY}/g;
        s/\Q$ENV{PUBLIC_PLACEHOLDER}\E/$ENV{PUBLIC_KEY}/g;
    ' "$config"

    echo "Updated: $config"

    IFS=:
done

IFS=$OLDIFS

echo "X25519 key generation completed successfully."

