#!/bin/bash

# Test script to demonstrate AES encryption functionality
echo "Testing AES-GCM Password Encryption..."
echo "======================================"

# Set AES key
export AES_KEY="my-secret-encryption-key-2024"

# Start server in background
echo "Starting server..."
./server &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo -e "\n1. Creating a site with password 'secret123':"
curl -s -X POST http://localhost:8080/sites \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test-site",
    "name": "Test Site",
    "username": "admin",
    "password": "secret123"
  }' | jq .

echo -e "\n2. Retrieving the site (password should be decrypted):"
curl -s http://localhost:8080/sites/test-site | jq .

echo -e "\n3. Checking the stored data file (password should be encrypted):"
echo "Raw data in sites.json:"
cat data/sites.json | jq .

echo -e "\n4. Updating the password to 'newpassword456':"
curl -s -X PUT http://localhost:8080/sites/test-site \
  -H "Content-Type: application/json" \
  -d '{
    "password": "newpassword456"
  }' | jq .

echo -e "\n5. Retrieving updated site:"
curl -s http://localhost:8080/sites/test-site | jq .

echo -e "\n6. Checking updated data file:"
cat data/sites.json | jq .

# Stop server
echo -e "\nStopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo -e "\n======================================"
echo "Encryption test completed!"
echo "Note: Passwords are encrypted in storage but decrypted in API responses"
