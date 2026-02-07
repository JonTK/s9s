# Error Handling Patterns

This document describes the error handling patterns and conventions used throughout the observability plugin.

## General Principles

### 1. Fail Fast, Fail Safe
- Validate inputs early and return meaningful errors
- Use graceful degradation for non-critical failures
- Never panic in production code (except for unrecoverable states)

### 2. Error Wrapping and Context
```go
// Good: Wrap errors with context
return fmt.Errorf("failed to query Prometheus: %w", err)

// Bad: Lose original error context
return errors.New("query failed")
```

### 3. Structured Error Types
```go
// Use custom error types for specific error categories
type ValidationError struct {
    Field   string
    Value   string
    Reason  string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Reason)
}
```

## Package-Specific Patterns

### Security Package

#### Rate Limiting Errors
```go
// Rate limit exceeded - client should retry with backoff
if !rateLimiter.Allow(clientID) {
    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
    return
}
```

#### Validation Errors
```go
// Validation failures should be logged internally but NOT expose details to clients
// (prevents information disclosure about validation rules/implementation)
if err := validator.ValidateRequest(req); err != nil {
    auditLogger.LogValidationFailure(req, err)
    // Security: Log error internally, return generic message to client
    log.Printf("Validation error (internal): %v", err)
    http.Error(w, "Request validation failed", http.StatusBadRequest)
    return
}
```

**Security Note**: Never expose raw `err.Error()` details in API responses for validation failures. This can leak:
- Internal validation rules and logic
- File paths or system information
- Implementation-specific details
- Patterns/signatures used in security checks

Always log detailed errors internally and return generic messages to clients.

#### Audit Logging Errors
```go
// Audit logging failures should not break request processing
if err := auditLogger.LogEvent(event); err != nil {
    // Log internally but continue processing
    log.Printf("Audit logging failed: %v", err)
}
```

### Prometheus Client Package

#### Connection Errors
```go
// Circuit breaker handles connection failures
client := NewCircuitBreakerClient(baseClient, config)
// Circuit automatically opens on repeated failures
```

#### Query Errors
```go
// Distinguish between client and server errors
result, err := client.Query(ctx, query, time.Now())
if err != nil {
    var promErr *PrometheusError
    if errors.As(err, &promErr) {
        if promErr.Type == "bad_data" {
            return nil, fmt.Errorf("invalid query: %w", err)
        }
        return nil, fmt.Errorf("prometheus server error: %w", err)
    }
    return nil, fmt.Errorf("connection error: %w", err)
}
```

#### Cache Errors
```go
// Cache misses are not errors, cache failures should degrade gracefully
value, err := cache.Get(key)
if err != nil {
    // Cache failure - fetch directly from source
    log.Printf("Cache error (degrading): %v", err)
    return client.Query(ctx, query, time)
}
```

### API Package

#### Request Processing Errors
```go
// Use appropriate HTTP status codes
func (api *ExternalAPI) handleError(w http.ResponseWriter, err error) {
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        api.writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    var timeoutErr *TimeoutError
    if errors.As(err, &timeoutErr) {
        api.writeError(w, http.StatusRequestTimeout, "Request timeout")
        return
    }
    
    // Default to internal server error
    api.writeError(w, http.StatusInternalServerError, "Internal error")
}
```

### Historical Package

#### Data Collection Errors
```go
// Continue collection despite individual metric failures
for _, query := range queries {
    result, err := client.Query(ctx, query, time.Now())
    if err != nil {
        log.Printf("Failed to collect metric %s: %v", query, err)
        continue // Don't stop collection for individual failures
    }
    // Process successful result
}
```

#### Storage Errors
```go
// Implement retry logic for storage operations
const maxRetries = 3
for attempt := 0; attempt < maxRetries; attempt++ {
    if err := storage.Save(data); err != nil {
        if attempt == maxRetries-1 {
            return fmt.Errorf("failed to save after %d attempts: %w", maxRetries, err)
        }
        time.Sleep(time.Duration(attempt+1) * time.Second)
        continue
    }
    break
}
```

## Error Recovery Strategies

### 1. Circuit Breaker Pattern
Used in Prometheus client to handle server failures:
- Open: All requests fail fast
- Half-open: Limited requests to test recovery
- Closed: Normal operation

### 2. Retry with Exponential Backoff
```go
func retryWithBackoff(operation func() error, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        if err := operation(); err != nil {
            if attempt == maxRetries-1 {
                return err
            }
            backoff := time.Duration(1<<attempt) * time.Second
            time.Sleep(backoff)
            continue
        }
        return nil
    }
    return nil
}
```

### 3. Graceful Degradation
```go
// Fallback to cached data when live queries fail
liveData, err := fetchLiveData(query)
if err != nil {
    log.Printf("Live data fetch failed, using cached data: %v", err)
    return fetchCachedData(query)
}
return liveData
```

## Testing Error Conditions

### 1. Error Injection Testing
```go
// Use interfaces to inject test failures
type MockPrometheusClient struct {
    ShouldFail bool
    FailureErr error
}

func (m *MockPrometheusClient) Query(ctx context.Context, query string, ts time.Time) (interface{}, error) {
    if m.ShouldFail {
        return nil, m.FailureErr
    }
    return normalResponse, nil
}
```

### 2. Timeout Testing
```go
// Test timeout handling
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
defer cancel()
_, err := client.Query(ctx, query, time.Now())
assert.Error(t, err)
assert.Contains(t, err.Error(), "context deadline exceeded")
```

### 3. Error Boundary Testing
```go
// Ensure errors don't propagate beyond boundaries
func TestErrorBoundary(t *testing.T) {
    // Create failing dependency
    failingClient := &MockClient{ShouldFail: true}
    
    // Component should handle failure gracefully
    component := NewComponent(failingClient)
    err := component.Process()
    
    // Should get specific error, not panic
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expected error context")
}
```

## Monitoring and Alerting

### Error Metrics Collection
- Error rates by endpoint and error type
- Circuit breaker state changes
- Retry attempt distributions
- Recovery time measurements

### Alert Thresholds
- Error rate > 5% for any endpoint
- Circuit breaker open for > 1 minute
- Audit logging failures
- Validation error spikes