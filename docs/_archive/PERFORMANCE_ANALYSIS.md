# Performance Analysis Guide

s9s includes comprehensive performance analysis tools to help you monitor and optimize your SLURM cluster operations and the s9s application itself.

## Overview

The performance analysis system consists of:
- **Real-time Profiler**: Monitors CPU, memory, goroutines, and operation latency
- **Performance Optimizer**: Provides recommendations and automatic tuning
- **Dashboard View**: Visual interface for monitoring performance metrics
- **Export Capabilities**: Export performance reports in multiple formats

## Accessing Performance Analysis

### Via UI
1. Launch s9s: `s9s` or `s9s --mock`
2. Press `p` or navigate to Performance view
3. Use keyboard shortcuts for different analysis modes

### Via Command Line
```bash
# Enable performance monitoring
s9s --profile

# Enable with specific interval
s9s --profile --profile-interval 5s

# Export performance report
s9s --export-performance /path/to/report.html
```

## Performance Metrics

### System Metrics

| Metric | Description | Unit | Optimal Range |
|--------|-------------|------|---------------|
| **CPU Usage** | Application CPU utilization | % | < 80% |
| **Memory Usage** | RAM consumption | MB/GB | < 1GB |
| **Goroutines** | Active goroutines | Count | < 1000 |
| **Heap Size** | Heap memory allocation | MB | Stable growth |
| **GC Pause** | Garbage collection pause time | ms | < 10ms |

### Operation Metrics

| Operation | Tracked Metrics | Purpose |
|-----------|----------------|---------|
| **Job Refresh** | Latency, success rate, throughput | Monitor SLURM API performance |
| **Node Updates** | Response time, error rate | Track cluster responsiveness |
| **UI Rendering** | Frame rate, render time | Ensure smooth user experience |
| **Data Export** | Export time, file size | Optimize export operations |

## Performance Dashboard

### Main Dashboard
The performance view provides:
- **Real-time Graphs**: CPU, memory, and network usage over time
- **Operation Timeline**: Recent operations with timing information
- **Alert Panel**: Active performance warnings and recommendations
- **Resource Utilization**: Detailed breakdown of system resources

### Key Performance Indicators (KPIs)
- **Response Time**: Average SLURM API response time
- **Refresh Rate**: Jobs/nodes updated per second
- **Error Rate**: Failed operations percentage
- **Memory Efficiency**: Memory usage per displayed item

## Performance Optimization

### Automatic Optimization

s9s includes an intelligent optimizer that:

1. **Monitors Performance**: Continuously tracks key metrics
2. **Identifies Issues**: Detects performance bottlenecks automatically
3. **Provides Recommendations**: Suggests specific improvements
4. **Auto-Fixes**: Applies safe optimizations automatically (when enabled)

### Optimization Levels

#### Light Optimization (Default)
- Adjusts refresh intervals based on activity
- Enables basic caching for static data
- Optimizes UI rendering for better responsiveness

#### Moderate Optimization
- All Light optimizations plus:
- Advanced caching strategies
- Background data prefetching
- Memory management tuning

#### Aggressive Optimization
- All previous optimizations plus:
- Reduces data resolution for better performance
- Implements more aggressive caching
- May sacrifice some features for speed

### Manual Optimizations

#### Configuration Tuning
```yaml
# ~/.s9s/config.yaml
performance:
  # Refresh intervals
  job_refresh_interval: 30s
  node_refresh_interval: 60s
  
  # Caching
  cache_ttl: 5m
  max_cache_size: 100MB
  
  # UI optimizations
  max_table_rows: 1000
  enable_virtualization: true
  
  # Network optimizations
  concurrent_requests: 5
  request_timeout: 30s
```

#### Environment Variables
```bash
# Performance tuning
export S9S_MAX_MEMORY=512MB
export S9S_CACHE_SIZE=50MB
export S9S_REFRESH_INTERVAL=30s

# Profiling
export S9S_PROFILE=true
export S9S_PROFILE_PATH=/tmp/s9s-profile
```

## Troubleshooting Performance Issues

### Common Issues and Solutions

#### 1. High Memory Usage
**Symptoms**: Memory consumption > 1GB, slow response
**Causes**: Large datasets, memory leaks, inefficient caching
**Solutions**:
- Reduce data retention period
- Enable memory optimization
- Restart application periodically

```bash
# Check memory usage
s9s --debug 2>&1 | grep "Memory:"

# Enable memory optimization
echo "performance.optimization_level: aggressive" >> ~/.s9s/config.yaml
```

#### 2. Slow SLURM API Responses
**Symptoms**: Long delays when refreshing data
**Causes**: Network latency, SLURM server overload, inefficient queries
**Solutions**:
- Increase request timeout
- Reduce refresh frequency
- Use connection pooling

```yaml
clusters:
  production:
    url: https://slurm.example.com
    timeout: 60s
    max_connections: 10
    retry_attempts: 3
```

#### 3. UI Lag and Stuttering
**Symptoms**: Slow navigation, delayed key responses
**Causes**: Large tables, complex rendering, background processing
**Solutions**:
- Enable table virtualization
- Reduce displayed columns
- Lower refresh rates

```yaml
preferences:
  jobs:
    columns: ["ID", "Name", "User", "State"]  # Limit columns
    max_rows: 500  # Limit displayed rows
```

#### 4. High CPU Usage
**Symptoms**: CPU usage > 80%, fan noise, battery drain
**Causes**: Frequent refreshes, complex operations, background tasks
**Solutions**:
- Increase refresh intervals
- Disable auto-refresh when inactive
- Use light optimization mode

### Performance Monitoring

#### Real-time Monitoring
```bash
# Monitor performance in separate terminal
s9s --debug 2>&1 | grep -E "(CPU|Memory|Latency):"

# Enable detailed profiling
s9s --profile --profile-cpu --profile-memory
```

#### Performance Reports
Generate detailed performance reports:

```bash
# Generate comprehensive report
s9s --export-performance report.html

# JSON format for analysis
s9s --export-performance report.json --format json

# CSV for spreadsheet analysis
s9s --export-performance report.csv --format csv
```

## Performance Benchmarking

### Built-in Benchmarks
s9s includes benchmarks for common operations:

```bash
# Run all benchmarks
s9s benchmark

# Specific benchmarks
s9s benchmark --jobs-refresh
s9s benchmark --table-rendering
s9s benchmark --export-operations

# Compare with baseline
s9s benchmark --baseline /path/to/baseline.json
```

### Custom Benchmarks
Create custom benchmarks for your environment:

```bash
# Benchmark with your data
s9s benchmark --cluster production --duration 5m

# Test different configurations
s9s benchmark --config aggressive.yaml
s9s benchmark --config conservative.yaml
```

## Performance Best Practices

### For Large Clusters (1000+ nodes/jobs)
1. **Increase Cache TTL**: Set cache_ttl to 10-15 minutes
2. **Limit Data Display**: Show only essential columns and recent jobs
3. **Use Filters**: Pre-filter data at the API level
4. **Enable Pagination**: Process data in smaller chunks

### For Slow Networks
1. **Increase Timeouts**: Set generous timeout values
2. **Reduce Refresh Frequency**: Use 60s+ refresh intervals
3. **Enable Compression**: If supported by SLURM API
4. **Use Local Caching**: Cache data locally for longer periods

### For Resource-Constrained Systems
1. **Limit Memory Usage**: Set max_memory limits
2. **Reduce Goroutines**: Lower concurrent_requests
3. **Disable Heavy Features**: Turn off real-time streaming
4. **Use Text Mode**: Disable fancy graphics and animations

## Integration with External Tools

### Prometheus Integration
Export metrics to Prometheus for long-term monitoring:

```yaml
integrations:
  prometheus:
    enabled: true
    endpoint: http://prometheus:9090/metrics
    push_interval: 30s
```

### Grafana Dashboards
Import pre-built Grafana dashboards:
- s9s Application Performance
- SLURM Cluster Metrics
- User Activity Monitoring

### Log Analysis
Integrate with log analysis tools:
```bash
# Export logs in structured format
s9s --log-format json --log-file /var/log/s9s.json

# Forward to ELK stack
s9s --log-output elasticsearch://elk:9200/s9s-logs
```

## API Performance Monitoring

### SLURM REST API Metrics
Monitor SLURM API performance:
- Request/response times
- Error rates by endpoint
- Throughput metrics
- Authentication overhead

### Connection Health
Track connection quality:
- DNS resolution time
- TCP connection time
- TLS handshake duration
- Keep-alive effectiveness

## Advanced Performance Analysis

### Memory Profiling
```bash
# Generate memory profile
s9s --profile-memory --profile-output mem.prof

# Analyze with go tool
go tool pprof mem.prof
```

### CPU Profiling
```bash
# Generate CPU profile
s9s --profile-cpu --profile-duration 30s

# Analyze hotspots
go tool pprof -http=:8080 cpu.prof
```

### Goroutine Analysis
```bash
# Check for goroutine leaks
s9s --debug 2>&1 | grep "Goroutines:"

# Generate goroutine dump
kill -QUIT <s9s-pid>  # Creates goroutine dump
```

---

For more information on performance optimization, see:
- [Configuration Guide](CONFIGURATION.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [API Documentation](API.md)