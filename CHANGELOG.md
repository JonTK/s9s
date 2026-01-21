# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Centralized version management in `internal/version` package
- GoReleaser configuration for automated releases
- CHANGELOG.md following Keep a Changelog format
- GitHub Actions release workflow
- Version-aware build targets in Makefile

## [0.1.0] - 2026-01-21

### Added
- **Core TUI Application**: Terminal-based user interface for SLURM cluster management
  - Vim-like navigation and keybindings
  - Real-time job and node monitoring
  - Interactive job management (submit, cancel, hold, release)
  - Multiple view modes (jobs, nodes, partitions, accounts, QOS, reservations, users)
  - SSH terminal integration for direct node access
  - Plugin system for extensibility

- **Configuration Management**:
  - Multi-cluster support with context switching
  - Configuration wizard (`s9s setup`) with auto-discovery
  - YAML-based configuration
  - Environment variable support
  - Mock mode for testing and development

- **Export Capabilities**:
  - CSV, JSON, and Markdown export formats
  - Job output streaming and filtering
  - Performance metrics export

- **Observability Plugin**:
  - Prometheus integration for metrics collection
  - Historical data collection and analysis
  - Performance dashboards
  - Resource efficiency analysis
  - Security audit logging
  - Rate limiting and circuit breaking

- **Testing & Quality**:
  - Comprehensive test suite
  - Mock SLURM client for testing
  - Integration tests
  - CI/CD with GitHub Actions

- **Code Quality**:
  - golangci-lint v2 integration
  - All errcheck, ineffassign, and govet issues resolved
  - Race condition fixes with proper synchronization
  - Thread-safe components

### Changed
- Upgraded to Go 1.24
- Updated slurm-client to upstream version v0.0.0-20260120203936
- Improved CI reliability with non-blocking security scans

### Fixed
- Race conditions in StatusBar, PerformanceDashboard, and MetricCache
- Deadlock in HistoricalDataCollector Stop() method
- Test isolation issues in key manager tests
- Config parser now supports both flat dotted keys and nested maps

### Security
- Security audit logging in observability plugin
- Rate limiting for API requests
- Circuit breaker pattern for fault tolerance
- Secrets management with encryption

## [Initial Development] - 2026-01-17 to 2026-01-20

### Added
- Initial project structure
- SLURM client integration
- Basic TUI components
- Plugin architecture
- Authentication and authorization
- SSH integration
- Export functionality

---

**Note**: Versions prior to 0.1.0 were in active development and did not follow semantic versioning.

[Unreleased]: https://github.com/jontk/s9s/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/jontk/s9s/releases/tag/v0.1.0
