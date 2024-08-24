#!/bin/bash

# Start the Go project

go run ./Backend/service/upload/main.go &
go run ./Backend/service/transfer/main.go &
./Frontend/node_modules/.bin/vite

# Wait for a few seconds to ensure the server is up and running


# Check if xdg-open is available and try to open the browser


