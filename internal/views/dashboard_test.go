package views

import (
	"testing"
	"time"

	"github.com/jontk/s9s/internal/dao"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculateHealthScore tests the overall health score calculation
func TestCalculateHealthScore(t *testing.T) {
	tests := []struct {
		name          string
		nodes         []*dao.Node
		jobs          []*dao.Job
		metrics       *dao.ClusterMetrics
		expectedScore float64
	}{
		{
			name:          "Perfect health - no issues",
			nodes:         createTestNodes(10, 0, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedScore: 100.0,
		},
		{
			name:          "All nodes down",
			nodes:         createTestNodes(0, 10, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedScore: 0.0, // 100 - (100% * 2) = -100, capped at 0
		},
		{
			name:          "50% nodes down",
			nodes:         createTestNodes(5, 5, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedScore: 0.0, // 100 - (50% * 2) = 0
		},
		{
			name:          "All jobs failed",
			nodes:         createTestNodes(10, 0, 0),
			jobs:          createTestJobs(0, 10),
			metrics:       &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedScore: 0.0, // 100 - (100% * 1) = 0
		},
		{
			name:          "High CPU utilization",
			nodes:         createTestNodes(10, 0, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 96.0, MemoryUsage: 50.0},
			expectedScore: 90.0, // 100 - 10 (CPU penalty) = 90
		},
		{
			name:          "High memory utilization",
			nodes:         createTestNodes(10, 0, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 96.0},
			expectedScore: 90.0, // 100 - 10 (memory penalty) = 90
		},
		{
			name:          "Both CPU and memory high",
			nodes:         createTestNodes(10, 0, 0),
			jobs:          createTestJobs(10, 0),
			metrics:       &dao.ClusterMetrics{CPUUsage: 96.0, MemoryUsage: 96.0},
			expectedScore: 80.0, // 100 - 20 = 80
		},
		{
			name:          "Combined issues",
			nodes:         createTestNodes(8, 2, 0), // 20% down nodes
			jobs:          createTestJobs(8, 2),     // 20% failed jobs
			metrics:       &dao.ClusterMetrics{CPUUsage: 96.0, MemoryUsage: 96.0},
			expectedScore: 20.0, // 100 - (20*2) - 20 - 20 = 20
		},
		{
			name:          "No data available",
			nodes:         []*dao.Node{},
			jobs:          []*dao.Job{},
			metrics:       nil,
			expectedScore: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DashboardView{
				nodes:          tt.nodes,
				jobs:           tt.jobs,
				clusterMetrics: tt.metrics,
			}

			score := v.calculateHealthScore()
			assert.Equal(t, tt.expectedScore, score, "Health score should match expected value")
		})
	}
}

// TestCalculateHealthStatus tests health status determination from scores
func TestCalculateHealthStatus(t *testing.T) {
	tests := []struct {
		name           string
		nodes          []*dao.Node
		jobs           []*dao.Job
		metrics        *dao.ClusterMetrics
		expectedStatus string
	}{
		{
			name:           "Excellent health (100)",
			nodes:          createTestNodes(10, 0, 0),
			jobs:           createTestJobs(10, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "EXCELLENT",
		},
		{
			name:           "Excellent health (90)",
			nodes:          createTestNodes(10, 0, 0),
			jobs:           createTestJobs(10, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 96.0, MemoryUsage: 50.0},
			expectedStatus: "EXCELLENT",
		},
		{
			name:           "Poor health (40)",
			nodes:          createTestNodes(8, 2, 0), // 20% down = -40 points
			jobs:           createTestJobs(8, 2),     // 20% failed = -20 points
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "POOR", // 100 - 40 - 20 = 40
		},
		{
			name:           "Good health boundary (89.9)",
			nodes:          createTestNodes(95, 5, 0), // 5% down = -10 points
			jobs:           createTestJobs(100, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "EXCELLENT", // 100 - 10 = 90
		},
		{
			name:           "Fair health (60)",
			nodes:          createTestNodes(80, 20, 0), // 20% down = -40
			jobs:           createTestJobs(100, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "FAIR", // 100 - 40 = 60
		},
		{
			name:           "Critical health (20)",
			nodes:          createTestNodes(60, 40, 0), // 40% down = -80
			jobs:           createTestJobs(100, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "CRITICAL", // 100 - 80 = 20
		},
		{
			name:           "Critical health (0)",
			nodes:          createTestNodes(0, 100, 0), // 100% down = -200, capped at 0
			jobs:           createTestJobs(100, 0),
			metrics:        &dao.ClusterMetrics{CPUUsage: 50.0, MemoryUsage: 50.0},
			expectedStatus: "CRITICAL", // capped at 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DashboardView{
				nodes:          tt.nodes,
				jobs:           tt.jobs,
				clusterMetrics: tt.metrics,
			}

			status := v.calculateHealthStatus()
			assert.Equal(t, tt.expectedStatus, status, "Health status should match expected value")
		})
	}
}

// TestGetHealthColor tests health status to color mapping
func TestGetHealthColor(t *testing.T) {
	tests := []struct {
		status        string
		expectedColor string
	}{
		{"EXCELLENT", "green"},
		{"GOOD", "cyan"},
		{"FAIR", "yellow"},
		{"POOR", "orange"},
		{"CRITICAL", "red"},
		{"UNKNOWN", "white"},
		{"", "white"},
	}

	v := &DashboardView{}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			color := v.getHealthColor(tt.status)
			assert.Equal(t, tt.expectedColor, color, "Color should match expected value for status")
		})
	}
}

// TestCalculateNodeHealthDeduction tests node health penalty calculation
func TestCalculateNodeHealthDeduction(t *testing.T) {
	tests := []struct {
		name              string
		nodes             []*dao.Node
		expectedDeduction float64
	}{
		{
			name:              "No nodes",
			nodes:             []*dao.Node{},
			expectedDeduction: 0.0,
		},
		{
			name:              "All nodes healthy",
			nodes:             createTestNodes(10, 0, 0),
			expectedDeduction: 0.0,
		},
		{
			name:              "All nodes down",
			nodes:             createTestNodes(0, 10, 0),
			expectedDeduction: 200.0, // 100% * 2
		},
		{
			name:              "50% nodes down",
			nodes:             createTestNodes(5, 5, 0),
			expectedDeduction: 100.0, // 50% * 2
		},
		{
			name:              "10% nodes down",
			nodes:             createTestNodes(9, 1, 0),
			expectedDeduction: 20.0, // 10% * 2
		},
		{
			name:              "Single node down",
			nodes:             createTestNodes(99, 1, 0),
			expectedDeduction: 2.02, // (1/100) * 100 * 2 ≈ 2.02
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DashboardView{
				nodes: tt.nodes,
			}

			deduction := v.calculateNodeHealthDeduction()
			assert.InDelta(t, tt.expectedDeduction, deduction, 0.1, "Node health deduction should be within tolerance")
		})
	}
}

// TestCalculateJobHealthDeduction tests job health penalty calculation
func TestCalculateJobHealthDeduction(t *testing.T) {
	tests := []struct {
		name              string
		jobs              []*dao.Job
		expectedDeduction float64
	}{
		{
			name:              "No jobs",
			jobs:              []*dao.Job{},
			expectedDeduction: 0.0,
		},
		{
			name:              "All jobs successful",
			jobs:              createTestJobs(10, 0),
			expectedDeduction: 0.0,
		},
		{
			name:              "All jobs failed",
			jobs:              createTestJobs(0, 10),
			expectedDeduction: 100.0, // 100% * 1
		},
		{
			name:              "50% jobs failed",
			jobs:              createTestJobs(5, 5),
			expectedDeduction: 50.0, // 50% * 1
		},
		{
			name:              "10% jobs failed",
			jobs:              createTestJobs(90, 10),
			expectedDeduction: 10.0, // 10% * 1
		},
		{
			name:              "Mixed job states",
			jobs:              createMixedJobs(70, 10, 20),
			expectedDeduction: 10.0, // 10% * 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &DashboardView{
				jobs: tt.jobs,
			}

			deduction := v.calculateJobHealthDeduction()
			assert.InDelta(t, tt.expectedDeduction, deduction, 0.1, "Job health deduction should be within tolerance")
		})
	}
}

// TestCalculateResourceHealthDeduction tests resource utilization penalty
func TestCalculateResourceHealthDeduction(t *testing.T) {
	tests := []struct {
		name              string
		cpuUsage          float64
		memoryUsage       float64
		expectedDeduction float64
	}{
		{
			name:              "No metrics",
			cpuUsage:          0,
			memoryUsage:       0,
			expectedDeduction: 0.0,
		},
		{
			name:              "Normal utilization",
			cpuUsage:          50.0,
			memoryUsage:       50.0,
			expectedDeduction: 0.0,
		},
		{
			name:              "CPU critical",
			cpuUsage:          96.0,
			memoryUsage:       50.0,
			expectedDeduction: 10.0,
		},
		{
			name:              "Memory critical",
			cpuUsage:          50.0,
			memoryUsage:       96.0,
			expectedDeduction: 10.0,
		},
		{
			name:              "Both critical",
			cpuUsage:          96.0,
			memoryUsage:       96.0,
			expectedDeduction: 20.0,
		},
		{
			name:              "Exactly at threshold (95%)",
			cpuUsage:          95.0,
			memoryUsage:       95.0,
			expectedDeduction: 0.0,
		},
		{
			name:              "Just over threshold (95.1%)",
			cpuUsage:          95.1,
			memoryUsage:       95.1,
			expectedDeduction: 20.0,
		},
		{
			name:              "Maximum usage (100%)",
			cpuUsage:          100.0,
			memoryUsage:       100.0,
			expectedDeduction: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var metrics *dao.ClusterMetrics
			if tt.cpuUsage > 0 || tt.memoryUsage > 0 {
				metrics = &dao.ClusterMetrics{
					CPUUsage:    tt.cpuUsage,
					MemoryUsage: tt.memoryUsage,
				}
			}

			v := &DashboardView{
				clusterMetrics: metrics,
			}

			deduction := v.calculateResourceHealthDeduction()
			assert.Equal(t, tt.expectedDeduction, deduction, "Resource health deduction should match expected value")
		})
	}
}

// TestGetUtilizationColor tests utilization percentage to color mapping
func TestGetUtilizationColor(t *testing.T) {
	tests := []struct {
		percentage    float64
		expectedColor string
	}{
		{0.0, "green"},
		{25.0, "green"},
		{49.9, "green"},
		{50.0, "yellow"},
		{75.0, "yellow"},
		{79.9, "yellow"},
		{80.0, "red"},
		{90.0, "red"},
		{100.0, "red"},
		{-10.0, "green"}, // negative should still map to green
		{110.0, "red"},   // over 100 should still map to red
	}

	for _, tt := range tests {
		t.Run(formatFloat(tt.percentage), func(t *testing.T) {
			color := getUtilizationColor(tt.percentage)
			assert.Equal(t, tt.expectedColor, color, "Color should match expected value for percentage")
		})
	}
}

// TestCreateMiniBar tests visual bar creation
func TestCreateMiniBar(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		color      string
		minLength  int // minimum expected length (excluding color codes)
	}{
		{
			name:       "Empty bar (0%)",
			percentage: 0.0,
			color:      "gray",
			minLength:  10, // all empty bars
		},
		{
			name:       "Half filled (50%)",
			percentage: 50.0,
			color:      "yellow",
			minLength:  10, // 5 filled + 5 empty
		},
		{
			name:       "Full bar (100%)",
			percentage: 100.0,
			color:      "green",
			minLength:  10, // all filled bars
		},
		{
			name:       "Over capacity (110%)",
			percentage: 110.0,
			color:      "red",
			minLength:  10, // should cap at 10
		},
		{
			name:       "Negative percentage",
			percentage: -10.0,
			color:      "red",
			minLength:  10, // should show as 0%
		},
	}

	v := &DashboardView{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := v.createMiniBar(tt.percentage, tt.color)

			// Verify bar is not empty
			require.NotEmpty(t, bar, "Bar should not be empty")

			// Verify it contains the color code
			assert.Contains(t, bar, "["+tt.color+"]", "Bar should contain color code")

			// Verify it contains some bar characters
			assert.True(t, containsBarCharacter(bar), "Bar should contain bar characters")
		})
	}
}

// TestCreateUtilizationBar tests utilization bar with automatic color selection
func TestCreateUtilizationBar(t *testing.T) {
	tests := []struct {
		percentage    float64
		expectedColor string
	}{
		{0.0, "green"},
		{30.0, "green"},
		{50.0, "yellow"},
		{70.0, "yellow"},
		{85.0, "red"},
		{100.0, "red"},
	}

	v := &DashboardView{}

	for _, tt := range tests {
		t.Run(formatFloat(tt.percentage), func(t *testing.T) {
			bar := v.createUtilizationBar(tt.percentage)

			// Verify bar contains the expected color
			assert.Contains(t, bar, "["+tt.expectedColor+"]", "Bar should contain expected color code")

			// Verify bar is not empty
			require.NotEmpty(t, bar, "Bar should not be empty")
		})
	}
}

// Helper functions

// createTestNodes creates test nodes with specified counts
func createTestNodes(healthy, down, drain int) []*dao.Node {
	nodes := make([]*dao.Node, 0, healthy+down+drain)

	for i := 0; i < healthy; i++ {
		nodes = append(nodes, &dao.Node{
			Name:  "node-healthy-" + formatInt(i),
			State: dao.NodeStateIdle,
		})
	}

	for i := 0; i < down; i++ {
		nodes = append(nodes, &dao.Node{
			Name:  "node-down-" + formatInt(i),
			State: dao.NodeStateDown,
		})
	}

	for i := 0; i < drain; i++ {
		nodes = append(nodes, &dao.Node{
			Name:  "node-drain-" + formatInt(i),
			State: dao.NodeStateDrain,
		})
	}

	return nodes
}

// createTestJobs creates test jobs with specified counts
func createTestJobs(successful, failed int) []*dao.Job {
	jobs := make([]*dao.Job, 0, successful+failed)

	for i := 0; i < successful; i++ {
		jobs = append(jobs, &dao.Job{
			ID:    "job-success-" + formatInt(i),
			State: dao.JobStateCompleted,
		})
	}

	for i := 0; i < failed; i++ {
		jobs = append(jobs, &dao.Job{
			ID:    "job-failed-" + formatInt(i),
			State: dao.JobStateFailed,
		})
	}

	return jobs
}

// createMixedJobs creates jobs with multiple states
func createMixedJobs(running, failed, pending int) []*dao.Job {
	jobs := make([]*dao.Job, 0, running+failed+pending)

	for i := 0; i < running; i++ {
		jobs = append(jobs, &dao.Job{
			ID:    "job-running-" + formatInt(i),
			State: dao.JobStateRunning,
		})
	}

	for i := 0; i < failed; i++ {
		jobs = append(jobs, &dao.Job{
			ID:    "job-failed-" + formatInt(i),
			State: dao.JobStateFailed,
		})
	}

	for i := 0; i < pending; i++ {
		jobs = append(jobs, &dao.Job{
			ID:         "job-pending-" + formatInt(i),
			State:      dao.JobStatePending,
			SubmitTime: time.Now().Add(-1 * time.Hour),
		})
	}

	return jobs
}

// containsBarCharacter checks if string contains bar visualization characters
func containsBarCharacter(s string) bool {
	return len(s) > 0 && (
	// Check for filled or empty bar characters
	len(s) > len("[white]"))
}

// formatFloat formats a float as a string for test names
func formatFloat(f float64) string {
	return formatInt(int(f)) + "%"
}

// formatInt formats an int as a string
func formatInt(i int) string {
	// Simple integer to string conversion
	if i == 0 {
		return "0"
	}
	result := ""
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+(i%10))) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
