#!/bin/bash

BASE_URL="http://localhost:8888"

echo "=== Testing Calendar API with CORS ==="
echo

echo "1. Testing CORS headers:"
curl -X OPTIONS $BASE_URL/create_event \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -I && echo

echo "2. Creating events:"
curl -X POST $BASE_URL/create_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-1","user_id":"user-123","date":"2025-01-15","title":"Встреча с командой"}' && echo

curl -X POST $BASE_URL/create_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-2","user_id":"user-123","date":"2025-01-15","title":"Презентация проекта"}' && echo

curl -X POST $BASE_URL/create_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-3","user_id":"user-123","date":"2025-01-16","title":"Обучение"}' && echo

echo

echo "3. Events for day 2025-01-15:"
curl "$BASE_URL/events_for_day?user_id=user-123&date=2025-01-15" && echo

echo "4. Events for week (starting 2025-01-15):"
curl "$BASE_URL/events_for_week?user_id=user-123&date=2025-01-15" && echo

echo "5. Events for month (January 2025):"
curl "$BASE_URL/events_for_month?user_id=user-123&date=2025-01-15" && echo

echo

echo "6. Updating event event-1:"
curl -X POST $BASE_URL/update_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-1","user_id":"user-123","date":"2025-01-15","title":"ВСТРЕЧА С КОМАНДОЙ (ОБНОВЛЕННАЯ)"}' && echo

echo

echo "7. Deleting event event-2:"
curl -X POST $BASE_URL/delete_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-2"}' && echo

echo

echo "8. Events after deletion:"
curl "$BASE_URL/events_for_day?user_id=user-123&date=2025-01-15" && echo

echo

echo "9. Error testing:"
echo "   - Empty ID:"
curl -X POST $BASE_URL/create_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"","user_id":"user-123","date":"2025-01-15","title":"Пустой ID"}' && echo

echo "   - Invalid date:"
curl -X POST $BASE_URL/create_event \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"id":"event-err","user_id":"user-123","date":"2025/01/15","title":"Неправильная дата"}' && echo

echo "=== Testing completed ==="