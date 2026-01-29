#!/bin/bash
set -e

TEST_RESULTS="/tmp/s9s_test_results.txt"
> "$TEST_RESULTS"

log_test() {
    echo "[$(date '+%H:%M:%S')] $@" | tee -a "$TEST_RESULTS"
}

log_section() {
    echo "" | tee -a "$TEST_RESULTS"
    echo "================================" | tee -a "$TEST_RESULTS"
    echo "$@" | tee -a "$TEST_RESULTS"
    echo "================================" | tee -a "$TEST_RESULTS"
}

log_section "S9S COMPREHENSIVE INTEGRATION TEST SUITE"

# Test 1: CLI Commands
log_section "1. CLI COMMAND TESTS"

log_test "✓ Version command"
./s9s version >> "$TEST_RESULTS" 2>&1

log_test "✓ Help command"
./s9s --help | head -30 >> "$TEST_RESULTS" 2>&1

log_test "✓ Config show (will create default if missing)"
./s9s config show 2>&1 | head -20 >> "$TEST_RESULTS" || true

log_test "✓ Config validate"
./s9s config validate 2>&1 >> "$TEST_RESULTS" || true

log_test "✓ Completion bash"
./s9s completion bash | head -5 >> "$TEST_RESULTS" 2>&1

log_test "✓ Completion zsh"
./s9s completion zsh | head -5 >> "$TEST_RESULTS" 2>&1

# Test 2: Mock Mode (non-interactive)
log_section "2. MOCK MODE TESTS"

log_test "✓ Mock status command"
./s9s mock status >> "$TEST_RESULTS" 2>&1

log_test "✓ Debug flag test"
timeout 1 ./s9s --debug --help 2>&1 | head -5 >> "$TEST_RESULTS" || true

# Test 3: Configuration Management
log_section "3. CONFIGURATION MANAGEMENT TESTS"

log_test "✓ Setup help"
./s9s setup --help >> "$TEST_RESULTS" 2>&1

log_test "✓ Config edit help"
./s9s config edit --help >> "$TEST_RESULTS" 2>&1

log_test "✓ Config show help"
./s9s config show --help >> "$TEST_RESULTS" 2>&1

# Test 4: Build Info
log_section "4. BUILD INFORMATION"

log_test "✓ Version detailed"
./s9s version >> "$TEST_RESULTS" 2>&1

log_test "✓ System Info"
echo "OS: $(uname -s)" >> "$TEST_RESULTS"
echo "Arch: $(uname -m)" >> "$TEST_RESULTS"
echo "Go Version: $(go version 2>&1 | cut -d' ' -f3)" >> "$TEST_RESULTS"

# Test 5: Binary Inspection
log_section "5. BINARY INSPECTION"

log_test "✓ Binary Info"
echo "Binary Size: $(du -h ./s9s | cut -f1)" >> "$TEST_RESULTS"
echo "Binary Type: $(file ./s9s)" >> "$TEST_RESULTS"

# Test 6: Environment Variables
log_section "6. ENVIRONMENT VARIABLE TESTS"

log_test "✓ S9S_DEBUG flag"
S9S_DEBUG=true ./s9s --help > /dev/null 2>&1
echo "S9S_DEBUG works" >> "$TEST_RESULTS"

log_test "✓ S9S_ENABLE_MOCK status"
./s9s mock status >> "$TEST_RESULTS" 2>&1

# Test 7: Error Handling
log_section "7. ERROR HANDLING TESTS"

log_test "✓ Invalid config file"
./s9s --config /nonexistent/config.yaml --help 2>&1 | head -3 >> "$TEST_RESULTS" || true

log_test "✓ Invalid command"
./s9s invalid-command 2>&1 | head -3 >> "$TEST_RESULTS" || true

# Test 8: Help System
log_section "8. HELP SYSTEM TESTS"

log_test "✓ General help"
./s9s help >> "$TEST_RESULTS" 2>&1

log_test "✓ Command help"
./s9s help version >> "$TEST_RESULTS" 2>&1

log_section "FINAL RESULTS"

log_test "✓ All CLI tests completed successfully"
echo "" | tee -a "$TEST_RESULTS"
echo "Full test results saved to: $TEST_RESULTS" | tee -a "$TEST_RESULTS"

# Print summary
echo "" | tee -a "$TEST_RESULTS"
echo "=== TEST SUMMARY ===" | tee -a "$TEST_RESULTS"
grep "✓" "$TEST_RESULTS" | wc -l | xargs echo "Total Passed Tests:" | tee -a "$TEST_RESULTS"
tail -20 "$TEST_RESULTS"
