#!/bin/bash
# End-to-end test script for the Todo API with Postgres.
#
# Setup steps before running this script:
#
#   1. Start the Postgres container:
#      docker compose up -d
#
#   2. Wait for it to be healthy:
#      docker compose ps
#
#   3. Run the database migration:
#      docker compose exec -T postgres psql -U todouser -d tododb < migrations/001_create_tasks_table.sql
#
#   4. Start the API server (in a separate terminal):
#      go run ./cmd/todod
#
#   5. Run this script:
#      ./scripts/test_api.sh

BASE_URL="http://localhost:8080"
PASS=0
FAIL=0

# Check if the server is running
if ! curl -s --max-time 2 "$BASE_URL/tasks" > /dev/null 2>&1; then
  echo "ERROR: API server is not running at $BASE_URL"
  echo ""
  echo "Start it in another terminal first:"
  echo "  go run ./cmd/todod"
  exit 1
fi

check() {
  local description="$1"
  local expected_code="$2"
  local actual_code="$3"
  local body="$4"

  if [ "$actual_code" -eq "$expected_code" ]; then
    echo "✓ $description"
    PASS=$((PASS + 1))
  else
    echo "✗ $description (expected $expected_code, got $actual_code)"
    echo "  response: $body"
    FAIL=$((FAIL + 1))
  fi
}

echo "=== Todo API E2E Tests ==="
echo ""

# --- Clean up: delete all existing tasks ---
echo "--- Setup: cleaning existing tasks ---"
existing=$(curl -s "$BASE_URL/tasks")
ids=$(echo "$existing" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
for id in $ids; do
  curl -s -X DELETE "$BASE_URL/tasks/$id" > /dev/null
done
echo ""

# --- Create ---
echo "--- Create ---"

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk and eggs","priority":"high"}')
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "POST /tasks - create task" 201 "$code" "$body"

task1_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{"title":"Read a book","priority":"low"}')
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "POST /tasks - create second task" 201 "$code" "$body"

task2_id=$(echo "$body" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{"description":"no title"}')
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "POST /tasks - missing title returns 400" 400 "$code" "$body"

echo ""

# --- List ---
echo "--- List ---"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/tasks")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "GET /tasks - list all" 200 "$code" "$body"

count=$(echo "$body" | grep -o '"id"' | wc -l)
if [ "$count" -eq 2 ]; then
  echo "  ✓ returned 2 tasks"
else
  echo "  ✗ expected 2 tasks, got $count"
  FAIL=$((FAIL + 1))
fi

echo ""

# --- Toggle Done ---
echo "--- Toggle Done ---"

response=$(curl -s -w "\n%{http_code}" -X PATCH "$BASE_URL/tasks/$task1_id/toggle")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "PATCH /tasks/$task1_id/toggle - toggle done" 200 "$code" "$body"

if echo "$body" | grep -q '"done":true'; then
  echo "  ✓ done is true"
else
  echo "  ✗ expected done=true"
  FAIL=$((FAIL + 1))
fi

echo ""

# --- List with filter ---
echo "--- List with filter ---"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/tasks?filter=done")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "GET /tasks?filter=done" 200 "$code" "$body"

done_count=$(echo "$body" | grep -o '"id"' | wc -l)
if [ "$done_count" -eq 1 ]; then
  echo "  ✓ returned 1 done task"
else
  echo "  ✗ expected 1 done task, got $done_count"
  FAIL=$((FAIL + 1))
fi

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/tasks?filter=pending")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "GET /tasks?filter=pending" 200 "$code" "$body"

pending_count=$(echo "$body" | grep -o '"id"' | wc -l)
if [ "$pending_count" -eq 1 ]; then
  echo "  ✓ returned 1 pending task"
else
  echo "  ✗ expected 1 pending task, got $pending_count"
  FAIL=$((FAIL + 1))
fi

echo ""

# --- Update ---
echo "--- Update ---"

response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/tasks/$task1_id" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy food","priority":"low"}')
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "PUT /tasks/$task1_id - update task" 200 "$code" "$body"

if echo "$body" | grep -q '"title":"Buy food"'; then
  echo "  ✓ title updated"
else
  echo "  ✗ title not updated"
  FAIL=$((FAIL + 1))
fi

response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/tasks/99999" \
  -H "Content-Type: application/json" \
  -d '{"title":"Nope"}')
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "PUT /tasks/99999 - not found returns 404" 404 "$code" "$body"

echo ""

# --- Delete ---
echo "--- Delete ---"

response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/tasks/$task2_id")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "DELETE /tasks/$task2_id - delete task" 200 "$code" "$body"

response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/tasks/99999")
body=$(echo "$response" | head -n -1)
code=$(echo "$response" | tail -n 1)
check "DELETE /tasks/99999 - not found returns 404" 404 "$code" "$body"

echo ""

# --- Summary ---
echo "==========================="
echo "Passed: $PASS"
echo "Failed: $FAIL"
echo "==========================="

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
