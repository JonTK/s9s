#!/bin/bash

echo "=== S9S MOCK MODE TESTING ==="
echo ""
echo "[1] Testing mock mode initialization..."
S9S_ENABLE_MOCK=dev timeout 2 ./s9s --mock 2>&1 | head -20 || true
echo ""

echo "[2] Checking mock status..."
./s9s mock status
echo ""

echo "[3] Testing with debug enabled..."
S9S_ENABLE_MOCK=dev ./s9s --mock --debug 2>&1 | head -5 || true
echo ""

echo "=== Mock Mode Tests Complete ==="
