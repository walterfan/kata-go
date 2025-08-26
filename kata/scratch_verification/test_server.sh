#!/bin/bash

# Test script for the HTTP server
# Make sure the server is running on port 8080 before running this script

echo "Testing HTTP Server API endpoints..."
echo "=================================="

# Test health endpoint
echo -e "\n1. Testing health endpoint:"
curl -s http://localhost:8080/health | jq .

# Test create user
echo -e "\n2. Creating a user:"
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user1",
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }' | jq .

# Test create another user
echo -e "\n3. Creating another user:"
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "user2",
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 25
  }' | jq .

# Test get all users
echo -e "\n4. Getting all users:"
curl -s http://localhost:8080/users | jq .

# Test get specific user
echo -e "\n5. Getting user1:"
curl -s http://localhost:8080/users/user1 | jq .

# Test update user
echo -e "\n6. Updating user1:"
curl -s -X PUT http://localhost:8080/users/user1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated Doe",
    "age": 31
  }' | jq .

# Test get updated user
echo -e "\n7. Getting updated user1:"
curl -s http://localhost:8080/users/user1 | jq .

# Test command execution
echo -e "\n8. Executing command 'pwd':"
curl -s -X POST http://localhost:8080/commands \
  -H "Content-Type: application/json" \
  -d '{"command": "pwd"}' | jq .

# Test command execution with date
echo -e "\n9. Executing command 'date':"
curl -s -X POST http://localhost:8080/commands \
  -H "Content-Type: application/json" \
  -d '{"command": "date"}' | jq .

# Test command execution with ls
echo -e "\n10. Executing command 'ls -la':"
curl -s -X POST http://localhost:8080/commands \
  -H "Content-Type: application/json" \
  -d '{"command": "ls -la"}' | jq .

# Test forbidden command
echo -e "\n11. Testing forbidden command (should fail):"
curl -s -X POST http://localhost:8080/commands \
  -H "Content-Type: application/json" \
  -d '{"command": "rm -rf /"}' | jq .

# Test delete user
echo -e "\n12. Deleting user2:"
curl -s -X DELETE http://localhost:8080/users/user2 | jq .

# Test get all users after deletion
echo -e "\n13. Getting all users after deletion:"
curl -s http://localhost:8080/users | jq .

echo -e "\n=================================="
echo "Testing completed!"
