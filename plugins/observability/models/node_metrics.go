package models

import (
	"fmt"
	"strings"
	"time"
)

// NodeMetrics represents metrics for a compute node
type NodeMetrics struct {
	NodeName      string             `json:"node_name"`
	NodeState     string             `json:"node_state"` // SLURM state
	LastUpdate    time.Time          `json:"last_update"`
	Resources     ResourceMetrics    `json:"resources"`
	JobCount      int                `json:"job_count"`
	Labels        map[string]string  `json:"labels"`         // Prometheus labels
	CustomMetrics map[string]float64 `json:"custom_metrics"` // Additional metrics
}

// NodeMetricsCollector collects and manages node metrics
type NodeMetricsCollector struct {
	nodes     map[string]*NodeMetrics
	nodeLabel string // Prometheus label for node identification
	// TODO(lint): Review unused code - field prometheusNode is unused
	// prometheusNode string // How nodes are identified in Prometheus
}

// NewNodeMetricsCollector creates a new node metrics collector
func NewNodeMetricsCollector(nodeLabel string) *NodeMetricsCollector {
	return &NodeMetricsCollector{
		nodes:     make(map[string]*NodeMetrics),
		nodeLabel: nodeLabel,
	}
}

// UpdateFromPrometheus updates node metrics from Prometheus data
func (nmc *NodeMetricsCollector) UpdateFromPrometheus(nodeName string, metrics map[string]*TimeSeries) {
	// Add logging to debug the issue
	if len(metrics) == 0 {
		return
	}

	node, exists := nmc.nodes[nodeName]
	if !exists {
		node = &NodeMetrics{
			NodeName:      nodeName,
			NodeState:     "up", // Default to up when metrics are available
			Labels:        make(map[string]string),
			CustomMetrics: make(map[string]float64),
		}
		nmc.nodes[nodeName] = node
	}

	// Update resource metrics
	node.Resources = nmc.extractResourceMetrics(metrics)
	node.LastUpdate = time.Now()

	// Extract labels from first metric
	for _, ts := range metrics {
		if len(ts.Labels) > 0 {
			node.Labels = ts.Labels
			break
		}
	}
}

// extractResourceMetrics extracts ResourceMetrics from Prometheus time series
func (nmc *NodeMetricsCollector) extractResourceMetrics(metrics map[string]*TimeSeries) ResourceMetrics {
	rm := ResourceMetrics{
		Timestamp: time.Now(),
	}

	nmc.extractCPUMetrics(&rm, metrics)
	nmc.extractMemoryMetrics(&rm, metrics)
	nmc.extractDiskMetrics(&rm, metrics)
	nmc.extractNetworkMetrics(&rm, metrics)

	return rm
}

// extractCPUMetrics extracts CPU-related metrics from time series data
func (nmc *NodeMetricsCollector) extractCPUMetrics(rm *ResourceMetrics, metrics map[string]*TimeSeries) {
	if val := getMetricValue(metrics, "node_cpu_usage"); val != nil {
		rm.CPU.Usage = *val
	}
	if val := getMetricValue(metrics, "node_cpu_cores"); val != nil {
		rm.CPU.Cores = int(*val)
	}
	if val := getMetricValue(metrics, "node_load_1m"); val != nil {
		rm.CPU.Load1m = *val
	}
	if val := getMetricValue(metrics, "node_load_5m"); val != nil {
		rm.CPU.Load5m = *val
	}
	if val := getMetricValue(metrics, "node_load_15m"); val != nil {
		rm.CPU.Load15m = *val
	}
}

// extractMemoryMetrics extracts memory-related metrics from time series data
func (nmc *NodeMetricsCollector) extractMemoryMetrics(rm *ResourceMetrics, metrics map[string]*TimeSeries) {
	if val := getMetricValue(metrics, "node_memory_total"); val != nil {
		rm.Memory.Total = uint64(*val)
	}
	if val := getMetricValue(metrics, "node_memory_available"); val != nil {
		rm.Memory.Available = uint64(*val)
	}
	if rm.Memory.Total > 0 && rm.Memory.Available > 0 {
		rm.Memory.Used = rm.Memory.Total - rm.Memory.Available
		rm.Memory.Usage = float64(rm.Memory.Used) / float64(rm.Memory.Total) * 100
	}
}

// extractDiskMetrics extracts disk I/O metrics from time series data
func (nmc *NodeMetricsCollector) extractDiskMetrics(rm *ResourceMetrics, metrics map[string]*TimeSeries) {
	if val := getMetricValue(metrics, "node_disk_read_bytes"); val != nil {
		rm.Disk.ReadBytesPerSec = *val
	}
	if val := getMetricValue(metrics, "node_disk_write_bytes"); val != nil {
		rm.Disk.WriteBytesPerSec = *val
	}
}

// extractNetworkMetrics extracts network metrics from time series data
func (nmc *NodeMetricsCollector) extractNetworkMetrics(rm *ResourceMetrics, metrics map[string]*TimeSeries) {
	if val := getMetricValue(metrics, "node_network_receive_bytes"); val != nil {
		rm.Network.ReceiveBytesPerSec = *val
	}
	if val := getMetricValue(metrics, "node_network_transmit_bytes"); val != nil {
		rm.Network.TransmitBytesPerSec = *val
	}
}

// getMetricValue safely retrieves a metric value from the time series map
func getMetricValue(metrics map[string]*TimeSeries, key string) *float64 {
	if ts, ok := metrics[key]; ok && ts != nil {
		if latest := ts.Latest(); latest != nil {
			return &latest.Value
		}
	}
	return nil
}

// GetNode returns metrics for a specific node
func (nmc *NodeMetricsCollector) GetNode(nodeName string) (*NodeMetrics, bool) {
	node, exists := nmc.nodes[nodeName]
	return node, exists
}

// GetAllNodes returns all node metrics
func (nmc *NodeMetricsCollector) GetAllNodes() map[string]*NodeMetrics {
	return nmc.nodes
}

// UpdateNodeState updates the SLURM state for a node
func (nmc *NodeMetricsCollector) UpdateNodeState(nodeName, state string, jobCount int) {
	node, exists := nmc.nodes[nodeName]
	if !exists {
		node = &NodeMetrics{
			NodeName:      nodeName,
			Labels:        make(map[string]string),
			CustomMetrics: make(map[string]float64),
		}
		nmc.nodes[nodeName] = node
	}

	node.NodeState = state
	node.JobCount = jobCount
}

// MapSLURMToPrometheus maps SLURM node names to Prometheus labels
func (nmc *NodeMetricsCollector) MapSLURMToPrometheus(slurmName string) string {
	// This is a simple implementation - in reality, this might need
	// more complex mapping based on your infrastructure

	// Common patterns:
	// 1. SLURM uses short names, Prometheus uses FQDNs
	// 2. Different naming conventions

	// For now, return as-is but this should be configurable
	return slurmName
}

// GetNodesSummary returns a summary of all nodes by state
func (nmc *NodeMetricsCollector) GetNodesSummary() map[string]int {
	summary := make(map[string]int)

	for _, node := range nmc.nodes {
		state := node.NodeState
		if state == "" {
			state = "unknown"
		}
		summary[state]++
	}

	return summary
}

// GetAggregateMetrics returns aggregate metrics across all nodes
func (nmc *NodeMetricsCollector) GetAggregateMetrics() *AggregateNodeMetrics {
	agg := &AggregateNodeMetrics{
		Timestamp: time.Now(),
	}

	nodeCount := 0

	for _, node := range nmc.nodes {
		if node.NodeState == "down" || node.NodeState == "drain" {
			continue
		}

		nodeCount++

		// Sum up resources
		agg.TotalCPUCores += node.Resources.CPU.Cores
		agg.TotalMemory += node.Resources.Memory.Total
		agg.UsedMemory += node.Resources.Memory.Used

		// Track utilization
		agg.TotalCPUUsage += node.Resources.CPU.Usage
		agg.TotalLoadAverage += node.Resources.CPU.Load1m

		// Sum I/O
		agg.TotalDiskRead += node.Resources.Disk.ReadBytesPerSec
		agg.TotalDiskWrite += node.Resources.Disk.WriteBytesPerSec
		agg.TotalNetworkRx += node.Resources.Network.ReceiveBytesPerSec
		agg.TotalNetworkTx += node.Resources.Network.TransmitBytesPerSec

		// Count jobs
		agg.TotalJobs += node.JobCount
	}

	// Calculate averages
	if nodeCount > 0 {
		agg.ActiveNodes = nodeCount
		agg.AverageCPUUsage = agg.TotalCPUUsage / float64(nodeCount)
		agg.AverageLoadPerCore = agg.TotalLoadAverage / float64(agg.TotalCPUCores)
		agg.MemoryUsagePercent = float64(agg.UsedMemory) / float64(agg.TotalMemory) * 100
	}

	return agg
}

// AggregateNodeMetrics represents cluster-wide aggregate metrics
type AggregateNodeMetrics struct {
	Timestamp          time.Time `json:"timestamp"`
	ActiveNodes        int       `json:"active_nodes"`
	TotalCPUCores      int       `json:"total_cpu_cores"`
	TotalMemory        uint64    `json:"total_memory"`
	UsedMemory         uint64    `json:"used_memory"`
	MemoryUsagePercent float64   `json:"memory_usage_percent"`
	TotalCPUUsage      float64   `json:"total_cpu_usage"`
	AverageCPUUsage    float64   `json:"average_cpu_usage"`
	TotalLoadAverage   float64   `json:"total_load_average"`
	AverageLoadPerCore float64   `json:"average_load_per_core"`
	TotalDiskRead      float64   `json:"total_disk_read"`
	TotalDiskWrite     float64   `json:"total_disk_write"`
	TotalNetworkRx     float64   `json:"total_network_rx"`
	TotalNetworkTx     float64   `json:"total_network_tx"`
	TotalJobs          int       `json:"total_jobs"`
}

// FormatNodeMetrics formats node metrics for display
func FormatNodeMetrics(node *NodeMetrics) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Node: %s", node.NodeName))
	parts = append(parts, fmt.Sprintf("State: %s", node.NodeState))
	parts = append(parts, fmt.Sprintf("Jobs: %d", node.JobCount))

	if node.Resources.CPU.Usage > 0 {
		parts = append(parts, fmt.Sprintf("CPU: %.1f%% (%d cores)",
			node.Resources.CPU.Usage, node.Resources.CPU.Cores))
	}

	if node.Resources.Memory.Total > 0 {
		parts = append(parts, fmt.Sprintf("Memory: %s / %s (%.1f%%)",
			FormatValue(float64(node.Resources.Memory.Used), "bytes"),
			FormatValue(float64(node.Resources.Memory.Total), "bytes"),
			node.Resources.Memory.Usage))
	}

	if node.Resources.CPU.Load1m > 0 {
		parts = append(parts, fmt.Sprintf("Load: %.2f, %.2f, %.2f",
			node.Resources.CPU.Load1m,
			node.Resources.CPU.Load5m,
			node.Resources.CPU.Load15m))
	}

	return strings.Join(parts, " | ")
}

// GetHealthStatus returns the health status of a node based on metrics
func (n *NodeMetrics) GetHealthStatus() string {
	if n.NodeState == "down" || n.NodeState == "drain" {
		return "unhealthy"
	}

	// Check various thresholds
	if n.Resources.CPU.Usage > 95 {
		return "critical"
	}
	if n.Resources.Memory.Usage > 95 {
		return "critical"
	}
	if n.Resources.CPU.Load1m > float64(n.Resources.CPU.Cores)*2 {
		return "warning"
	}
	if n.Resources.CPU.Usage > 80 || n.Resources.Memory.Usage > 80 {
		return "warning"
	}

	return "healthy"
}
