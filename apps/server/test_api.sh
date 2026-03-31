#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"
COOKIES="cookies.txt"
EMAIL="testuser_$(date +%s)@test.com"
PASSWORD="Password123!"

echo "====================================="
echo "Testing Coderz.space API endpoints"
echo "====================================="

rm -f $COOKIES

echo -e "\n1. Health Check [GET /api/health]"
curl -s -X GET "http://localhost:8080/api/health" | grep "ok" > /dev/null && echo "✅ OK" || echo "❌ FAILED"

echo -e "\n2. Signup [POST /api/v1/auth/signup]"
SIGNUP_RES=$(curl -s -X POST "$BASE_URL/auth/signup" \
  -H "Content-Type: application/json" \
  -H "X-Requested-With: XMLHttpRequest" \
  -c $COOKIES \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"name\":\"Test User\"}")

echo "$SIGNUP_RES" | grep "true" > /dev/null && echo "✅ OK" || echo "❌ FAILED: $SIGNUP_RES"

echo -e "\n3. Login [POST /api/v1/auth/login]"
LOGIN_RES=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -H "X-Requested-With: XMLHttpRequest" \
  -c $COOKIES -b $COOKIES \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

echo "$LOGIN_RES" | grep "true" > /dev/null && echo "✅ OK" || echo "❌ FAILED: $LOGIN_RES"

echo -e "\n4. Get Current User [GET /api/v1/auth/me]"
ME_RES=$(curl -s -X GET "$BASE_URL/auth/me" \
  -b $COOKIES \
  -H "X-Requested-With: XMLHttpRequest")

echo "$ME_RES" | grep "true" > /dev/null && echo "✅ OK" || echo "❌ FAILED: $ME_RES"

USER_ID=$(echo "$ME_RES" | grep -o '"id":"[^"]*' | head -n 1 | cut -d'"' -f4)
if [ -z "$USER_ID" ]; then
    echo "❌ Failed to extract user ID. Stopping tests."
    exit 1
fi
echo "Extracted User ID: $USER_ID"

echo -e "\n5. Request Mentee Role [POST /api/v1/mentorship/request-role]"
ROLE_RES=$(curl -s -X POST "$BASE_URL/mentorship/request-role" \
  -b $COOKIES \
  -H "Content-Type: application/json" \
  -d "{\"role\":\"mentee\"}")

echo "$ROLE_RES" | grep "true" > /dev/null && echo "✅ OK" || echo "❌ FAILED: $ROLE_RES"

echo -e "\n6. List Mentorship Requests [GET /api/v1/mentorship/requests]"
REQ_LIST_RES=$(curl -s -X GET "$BASE_URL/mentorship/requests" \
  -b $COOKIES)
echo "Response: $REQ_LIST_RES"

echo -e "\n7. Get Mentee Profile [GET /api/v1/mentees/:id/profile]"
PROFILE_RES=$(curl -s -X GET "$BASE_URL/mentees/$USER_ID/profile" \
  -b $COOKIES)
echo "Response: $PROFILE_RES"

echo -e "\n8. Get Leaderboard [GET /api/v1/leaderboard]"
LEADERBOARD_RES=$(curl -s -X GET "$BASE_URL/leaderboard" \
  -b $COOKIES)
echo "Response: $LEADERBOARD_RES"

echo -e "\n9. Get Mentee Tasks [GET /api/v1/mentees/:id/questions]"
TASKS_RES=$(curl -s -X GET "$BASE_URL/mentees/$USER_ID/questions" \
  -b $COOKIES)
echo "Response: $TASKS_RES"

echo -e "\n====================================="
echo "API Test Complete"
echo "====================================="
