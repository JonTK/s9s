#!/bin/bash

# Test script for mock validation system
S9S_BINARY="/tmp/s9s-mock-test"

echo "🧪 Testing Mock Validation System"
echo "================================="

# Test 1: Mock without environment variable
echo "Test 1: Mock without environment variable"
unset S9S_ENABLE_MOCK
unset ENVIRONMENT
result=$(${S9S_BINARY} --mock 2>&1)
if echo "$result" | grep -q "mock mode disabled"; then
    echo "✅ PASS: Mock correctly blocked without environment variable"
else
    echo "❌ FAIL: Mock should be blocked without environment variable"
fi
echo

# Test 2: Mock with development environment variable
echo "Test 2: Mock with development environment variable"
export S9S_ENABLE_MOCK=development
result=$(${S9S_BINARY} --mock 2>&1)
if echo "$result" | grep -q "mock mode disabled"; then
    echo "❌ FAIL: Mock should be allowed with development environment variable"
else
    echo "✅ PASS: Mock allowed with development environment variable"
fi
echo

# Test 3: Mock in production environment (should warn and prompt)
echo "Test 3: Mock in production environment"
export S9S_ENABLE_MOCK=development
export ENVIRONMENT=production
result=$(${S9S_BINARY} --mock 2>&1)
if echo "$result" | grep -q "WARNING.*production environment"; then
    echo "✅ PASS: Production warning correctly shown"
else
    echo "❌ FAIL: Production warning should be shown"
fi
if echo "$result" | grep -q "cancelled by user"; then
    echo "✅ PASS: Mock correctly cancelled in non-interactive production"
else
    echo "❌ FAIL: Mock should be cancelled in non-interactive production"
fi
echo

# Test 4: Different allowed environment values
echo "Test 4: Different allowed environment values"
unset ENVIRONMENT
for env_val in "testing" "debug" "dev" "local" "true"; do
    export S9S_ENABLE_MOCK=$env_val
    result=$(${S9S_BINARY} --mock 2>&1)
    if echo "$result" | grep -q "mock mode disabled"; then
        echo "❌ FAIL: S9S_ENABLE_MOCK=$env_val should be allowed"
    else
        echo "✅ PASS: S9S_ENABLE_MOCK=$env_val correctly allowed"
    fi
done
echo

# Test 5: Invalid environment values
echo "Test 5: Invalid environment values"
for env_val in "invalid" "false" "0" "production"; do
    export S9S_ENABLE_MOCK=$env_val
    result=$(${S9S_BINARY} --mock 2>&1)
    if echo "$result" | grep -q "mock mode disabled"; then
        echo "✅ PASS: S9S_ENABLE_MOCK=$env_val correctly rejected"
    else
        echo "❌ FAIL: S9S_ENABLE_MOCK=$env_val should be rejected"
    fi
done
echo

# Test 6: Help text mentions environment variable
echo "Test 6: Help text mentions environment variable"
result=$(${S9S_BINARY} --help 2>&1)
if echo "$result" | grep -q "S9S_ENABLE_MOCK"; then
    echo "✅ PASS: Help text mentions S9S_ENABLE_MOCK"
else
    echo "❌ FAIL: Help text should mention S9S_ENABLE_MOCK"
fi
if echo "$result" | grep -q "requires S9S_ENABLE_MOCK"; then
    echo "✅ PASS: --mock flag mentions requirement"
else
    echo "❌ FAIL: --mock flag should mention requirement"
fi
echo

echo "🎯 Mock validation testing complete!"
echo "💡 Remember to set S9S_ENABLE_MOCK=development for local development"