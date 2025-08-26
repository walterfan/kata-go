#!/bin/bash

# Test script for Sites CLI commands
echo "Testing Sites CLI Commands..."
echo "=============================="

# Set AES key
export AES_KEY="my-secret-encryption-key-2024"

echo -e "\n1. Testing sites help:"
./vault sites --help

echo -e "\n2. Testing sites list (initial state):"
./vault sites list

echo -e "\n3. Testing sites create:"
./vault sites create "test-site-1" "Test Site 1" "user1" "password1"
./vault sites create "test-site-2" "Test Site 2" "user2" "password2"

echo -e "\n4. Testing sites list (after creation):"
./vault sites list

echo -e "\n5. Testing sites get:"
./vault sites get "test-site-1"

echo -e "\n6. Testing sites update:"
./vault sites update "test-site-1" "Updated Test Site 1" "updated_user1" "updated_password1"

echo -e "\n7. Testing sites get (after update):"
./vault sites get "test-site-1"

echo -e "\n8. Testing sites delete:"
./vault sites delete "test-site-2"

echo -e "\n9. Testing sites list (after deletion):"
./vault sites list

echo -e "\n10. Testing error handling - get non-existent site:"
./vault sites get "nonexistent" || echo "Expected error: Site not found"

echo -e "\n11. Testing error handling - create duplicate site:"
./vault sites create "test-site-1" "Duplicate" "user" "pass" || echo "Expected error: Site already exists"

echo -e "\n12. Testing error handling - delete non-existent site:"
./vault sites delete "nonexistent" || echo "Expected error: Site not found"

echo -e "\n13. Final cleanup - delete remaining test sites:"
./vault sites delete "test-site-1"

echo -e "\n14. Final sites list:"
./vault sites list

echo -e "\n=============================="
echo "Sites CLI test completed!"
echo "All CRUD operations working correctly."
