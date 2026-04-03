#!/bin/bash

# Go Knowledge Base - Automatic Index Generator
# This script is a wrapper that calls the Python implementation
# For Windows, use: generate-index.ps1
# 
# Usage: ./generate-index.sh [knowledge_base_path]
# Default path: parent of scripts directory

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
KB_PATH="${1:-$SCRIPT_DIR/..}"

# Check if Python is available
if command -v python3 &> /dev/null; then
    PYTHON_CMD="python3"
elif command -v python &> /dev/null; then
    PYTHON_CMD="python"
else
    echo "[ERROR] Python is required but not installed."
    echo "Please install Python 3.6+ and try again."
    exit 1
fi

# Run the Python script
exec "$PYTHON_CMD" "$SCRIPT_DIR/generate-index.py" "$KB_PATH"
