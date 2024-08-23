#!/bin/bash

# Start the Go project

go run ./Backend/service/upload/main.go &
go run ./Backend/service/transfer/main.go &
./Frontend/node_modules/.bin/vite

# Wait for a few seconds to ensure the server is up and running
sleep 5

# Check if xdg-open is available and try to open the browser
if command -v xdg-open > /dev/null; then
    xdg-open http://localhost:8081/user/login
else
    echo "Please open a web browser and go to http://localhost:8081/user/login"
fi

