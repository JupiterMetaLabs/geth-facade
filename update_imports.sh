#!/bin/bash

# Script to update import paths after restructuring

echo "Updating import paths..."

# Update Services files
echo "Updating Services files..."

# Update handlers.go
sed -i '' 's/package rpc/package Services/g' Services/handlers.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/backend|github.com/jupitermetalabs/geth-facade/Types|g' Services/handlers.go
sed -i '' 's/backend\./Types\./g' Services/handlers.go

# Update http_server.go
sed -i '' 's/package rpc/package Services/g' Services/http_server.go

# Update ws_server.go
sed -i '' 's/package rpc/package Services/g' Services/ws_server.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/backend|github.com/jupitermetalabs/geth-facade/Types|g' Services/ws_server.go
sed -i '' 's/backend\./Types\./g' Services/ws_server.go

# Update memory.go
sed -i '' 's/package memorybackend/package Services/g' Services/memory.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/backend|github.com/jupitermetalabs/geth-facade/Types|g' Services/memory.go
sed -i '' 's/backend\./Types\./g' Services/memory.go

# Update facade.go
sed -i '' 's/package jmdtgethfacade/package Services/g' Services/facade.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/backend|github.com/jupitermetalabs/geth-facade/Types|g' Services/facade.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/rpc|github.com/jupitermetalabs/geth-facade/Services|g' Services/facade.go
sed -i '' 's/backend\./Types\./g' Services/facade.go
sed -i '' 's/rpc\./Services\./g' Services/facade.go

# Update doc.go
sed -i '' 's/package jmdtgethfacade/package Services/g' Services/doc.go

# Update Types files
echo "Updating Types files..."

# Update backend.go
sed -i '' 's/package backend/package Types/g' Types/backend.go

# Update types.go
sed -i '' 's/package rpc/package Types/g' Types/types.go

# Update main.go
echo "Updating main.go..."
sed -i '' 's|github.com/jupitermetalabs/geth-facade/pkg/jmdtgethfacade|github.com/jupitermetalabs/geth-facade/Services|g' main.go
sed -i '' 's|github.com/jupitermetalabs/geth-facade/pkg/memorybackend|github.com/jupitermetalabs/geth-facade/Services|g' main.go
sed -i '' 's/jmdtgethfacade\./Services\./g' main.go
sed -i '' 's/memorybackend\./Services\./g' main.go

echo "Import paths updated successfully!"
