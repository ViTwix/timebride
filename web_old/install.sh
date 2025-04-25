#!/bin/bash

# Install dependencies
npm install

# Create necessary directories if they don't exist
mkdir -p src/components
mkdir -p src/pages
mkdir -p src/services
mkdir -p src/utils
mkdir -p src/hooks
mkdir -p src/contexts
mkdir -p src/types
mkdir -p src/assets

# Create index files
touch src/index.tsx
touch src/react-app-env.d.ts

echo "Frontend setup completed!" 