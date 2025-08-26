#!/bin/bash

# Test script for the HTTP server
# Make sure the server is running on port 8080 before running this script

echo "Testing HTTP Server API endpoints..."
echo "=================================="

# Test health endpoint
echo -e "\n1. Testing health endpoint:"
curl -s http://localhost:8080/health | jq .

# Test create site
echo -e "\n2. Creating a site:"
curl -s -X POST http://localhost:8080/sites \
  -H "Content-Type: application/json" \
  -d '{
    "id": "site1",
    "name": "Example Site 1",
    "username": "admin",
    "password": "password123"
  }' | jq .

# Test create another site
echo -e "\n3. Creating another site:"
curl -s -X POST http://localhost:8080/sites \
  -H "Content-Type: application/json" \
  -d '{
    "id": "site2",
    "name": "Example Site 2",
    "username": "user",
    "password": "secret456"
  }' | jq .

# Test get all sites
echo -e "\n4. Getting all sites:"
curl -s http://localhost:8080/sites | jq .

# Test get specific site
echo -e "\n5. Getting site1:"
curl -s http://localhost:8080/sites/site1 | jq .

# Test update site
echo -e "\n6. Updating site1:"
curl -s -X PUT http://localhost:8080/sites/site1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Example Site 1",
    "username": "newadmin"
  }' | jq .

# Test get updated site
echo -e "\n7. Getting updated site1:"
curl -s http://localhost:8080/sites/site1 | jq .

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

# Test delete site
echo -e "\n12. Deleting site2:"
curl -s -X DELETE http://localhost:8080/sites/site2 | jq .

# Test get all sites after deletion
echo -e "\n13. Getting all sites after deletion:"
curl -s http://localhost:8080/sites | jq .

echo -e "\n=================================="
echo "Testing completed!"
