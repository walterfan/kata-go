#!/bin/bash

# Test script for the weather tool

echo "Testing weather tool..."

# Start the tool server (this would normally be done by the main program)
go run main.go --help > /dev/null 2>&1 &

# Wait a moment for the server to start
sleep 2

# Test the weather API
curl -X POST http://localhost:8080/tool/get_weather \
  -H "Content-Type: application/json" \
  -d '{"location": "Beijing"}'

echo ""
echo "Weather tool test completed!" 