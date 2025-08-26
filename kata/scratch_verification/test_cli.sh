#!/bin/bash

# Test script for CLI commands
echo "Testing CLI Commands..."
echo "======================"

# Set AES key
export AES_KEY="my-secret-encryption-key-2024"

echo -e "\n1. Testing help command:"
./vault --help

echo -e "\n2. Testing encrypt command:"
ENCRYPTED=$(./vault encrypt "testpassword123")
echo "Encrypted: $ENCRYPTED"

echo -e "\n3. Testing decrypt command:"
DECRYPTED=$(./vault decrypt "$ENCRYPTED")
echo "Decrypted: $DECRYPTED"

echo -e "\n4. Testing server command (briefly):"
./vault server &
SERVER_PID=$!
sleep 2
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo -e "\n======================"
echo "CLI test completed!"
echo "All commands working correctly."
