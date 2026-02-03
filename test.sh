#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"
API_KEY="${API_KEY:-test-api-key}"

echo -e "${BLUE}=== Testing AI Scam Engagement Backend ===${NC}\n"

# Test 1: Health Check
echo -e "${BLUE}Test 1: Health Check${NC}"
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" ${BASE_URL}/health)
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
BODY=$(echo "$HEALTH_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    echo "Response: $BODY"
else
    echo -e "${RED}✗ Health check failed (HTTP $HTTP_CODE)${NC}"
fi

echo -e "\n---\n"

# Test 2: Initial Scam Detection
echo -e "${BLUE}Test 2: Initial Scam Detection${NC}"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d @examples/request1_initial_scam.json)

echo "Request: Initial scam message with urgency and account threat"
echo "Response:"
echo "$RESPONSE" | jq '.'

STATUS=$(echo "$RESPONSE" | jq -r '.status')
if [ "$STATUS" = "success" ]; then
    echo -e "${GREEN}✓ Status correctly set to 'success'${NC}"
else
    echo -e "${RED}✗ Expected status 'success', got $STATUS${NC}"
fi

echo -e "\n---\n"

# Test 3: Intelligence Extraction
echo -e "${BLUE}Test 3: Intelligence Extraction (UPI + Link)${NC}"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d @examples/request2_intel_extraction.json)

echo "Request: Message with UPI ID and phishing link"
echo "Response:"
echo "$RESPONSE" | jq '.'

STATUS=$(echo "$RESPONSE" | jq -r '.status')
REPLY=$(echo "$RESPONSE" | jq -r '.reply')
if [ "$STATUS" = "success" ] && [ -n "$REPLY" ]; then
    echo -e "${GREEN}✓ Response format correct${NC}"
else
    echo -e "${RED}✗ Response format incorrect${NC}"
fi

echo -e "\n---\n"

# Test 4: Complete Intelligence Collection
echo -e "${BLUE}Test 4: Complete Intelligence Collection${NC}"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d @examples/request3_complete_intel.json)

echo "Request: Message with phone numbers and bank account"
echo "Response:"
echo "$RESPONSE" | jq '.'

STATUS=$(echo "$RESPONSE" | jq -r '.status')
if [ "$STATUS" = "success" ]; then
    echo -e "${GREEN}✓ Session processed successfully${NC}"
    echo -e "${GREEN}✓ Check server logs for callback to GUVI endpoint${NC}"
else
    echo -e "${RED}✗ Expected status 'success', got $STATUS${NC}"
fi

echo -e "\n---\n"

# Test 5: Non-Scam Message
echo -e "${BLUE}Test 5: Non-Scam Message (Low Score)${NC}"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d '{
    "sessionId": "demo-session-002",
    "message": {
      "sender": "scammer",
      "text": "Hello, how are you today?",
      "timestamp": "2026-02-02T10:00:00Z"
    },
    "conversationHistory": [],
    "metadata": {
      "channel": "SMS",
      "language": "English",
      "locale": "IN"
    }
  }')

echo "Request: Benign message"
echo "Response:"
echo "$RESPONSE" | jq '.'

STATUS=$(echo "$RESPONSE" | jq -r '.status')
if [ "$STATUS" = "success" ]; then
    echo -e "${GREEN}✓ Response format correct for non-scam${NC}"
else
    echo -e "${RED}✗ Expected status 'success', got $STATUS${NC}"
fi

echo -e "\n---\n"

# Test 6: API Key Authentication
echo -e "${BLUE}Test 6: API Key Authentication${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST ${BASE_URL}/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: wrong-key" \
  -d @examples/request1_initial_scam.json)
  
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "401" ]; then
    echo -e "${GREEN}✓ API key authentication working${NC}"
else
    echo -e "${RED}✗ Expected 401 for invalid API key, got $HTTP_CODE${NC}"
fi

echo -e "\n${BLUE}=== All Tests Complete ===${NC}"
