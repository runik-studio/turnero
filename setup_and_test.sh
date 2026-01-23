#!/bin/bash

echo "Installing dependencies..."
go mod tidy

echo "Generating docs..."
./update_docs.sh

echo "Starting server in background..."
export MOCK_AUTH=true
go run cmd/api/main.go &
PID=$!
sleep 5

echo "Running tests..."

echo "Testing POST /auth/login"
curl -X POST -H "Authorization: Bearer mock-token" -H "Content-Type: application/json" -d '{"email": "test@example.com", "password": "password123"}' http://localhost:8080/auth/login
echo "\n"
echo "Testing POST /api/services"
curl -X POST -H "Authorization: Bearer mock-token" -H "Content-Type: application/json" -d '{"title": "test_title", "description": "test_description", "icon_url": "test_icon_url", "duration_minutes": 10}' http://localhost:8080/api/services
echo "\n"
echo "Testing POST /api/providers"
curl -X POST -H "Authorization: Bearer mock-token" -H "Content-Type: application/json" -d '{"address": "test_address", "avatar_url": "test_avatar_url", "full_name": "test_full_name", "establishment_name": "test_establishment_name"}' http://localhost:8080/api/providers
echo "\n"
echo "Testing POST /api/appointments"
curl -X POST -H "Authorization: Bearer mock-token" -H "Content-Type: application/json" -d '{"scheduled_at": "2023-01-01T00:00:00Z", "status": "test_status", "notes": "test_notes", "user": "test_user", "service": "test_service", "provider": "test_provider"}' http://localhost:8080/api/appointments
echo "\n"
echo "Testing POST /api/users"
curl -X POST -H "Authorization: Bearer mock-token" -H "Content-Type: application/json" -d '{"picture": "test_picture", "role_id": "test_role_id", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z", "uid": "test_uid", "email": "test_email", "name": "test_name"}' http://localhost:8080/api/users
echo "\n"

echo "Killing server (PID: $PID)..."
kill $PID
