# Phase 4: Advanced Linters & Code Quality Maintenance

This document covers Phase 4 of our linting journey - enabling and configuring advanced linters that require more careful planning and implementation.

## Overview

Phase 4 focuses on advanced code quality linters that catch complex issues like code duplication, cyclomatic complexity, and structural patterns. These linters often have higher barriers to entry and may require architectural decisions or significant refactoring.

**Status**: Partially complete
- ⏳ Phase 4a: revive - DISABLED (package-level violation suppression issue)
- ✅ Phase 4c: gocognit - ENABLED with appropriate thresholds
- ⏳ Phase 4b: dupl - DISABLED (exclude-rules not working)
- ⏳ Phase 4d: cyclop - DISABLED (requires significant refactoring)
- ⏳ Phase 4e: containedctx - DISABLED (architectural pattern conflict)

## Phase 4a: Revive Linter ⏳ DISABLED

**What it does**: Enforces Go idioms and style conventions
- Package documentation requirements
- Exported symbol naming conventions
- Error handling patterns
- Naming consistency

**Implementation**: PR #32 (revive enabled, but disabled in subsequent fix)
- Added 40 `//nolint:revive` directives to type aliases and other intentional violations
- 36 type alias violations for backward compatibility (e.g., `type AuthProvider = Provider`)
- 2 package naming exceptions (internal/errors, plugins/observability/api) - **now disabled**
- 1 empty block in benchmark test
- 1 unused-parameter suppression

**Current Status**: DISABLED - Known Issue

**Problem**:
Two package-level revive violations cannot be suppressed:
1. `internal/errors/errors_test.go`: "avoid package names that conflict with Go standard library package names"
2. `plugins/observability/api/external.go`: "avoid meaningless package names"

**Root Cause**:
- Go doesn't allow `//nolint` comments on package statements (the comment must be immediately before the construct)
- golangci-lint `exclude-rules` do not match package-level violations properly
- The `.revive.toml` configuration file var-naming rule arguments don't support package exemptions

**Attempted Solutions**:
1. Adding nolint directives on package line - Go parsing error
2. Adding nolint directives on preceding line - Revive still reports violation
3. exclude-rules in golangci.yml - Pattern matching fails for package-level violations
4. .revive.toml var-naming arguments - Does not support package name exemptions

**Workaround**:
- Revive linter disabled in `.golangci.yml`
- 40 nolint directives for type-alias violations remain valid and documented
- Package naming violations are minor and acceptable for legitimate uses

**Example (Type Aliases - Still Valid)**:
```go
//nolint:revive // type alias for backward compatibility
type AuthProvider = Provider
```

**TODO (Phase 5)**:
- Re-enable revive after finding proper solution for package-level violations
- Options to investigate:
  1. Configure revive to accept "api" and "errors" as valid package names
  2. Document as known limitation with accepted violations
  3. Contribute fix to revive or golangci-lint for better package-level handling

## Phase 4c: Gocognit Linter ✅ ENABLED

**What it does**: Measures cognitive complexity - how difficult code is to understand
- Different from cyclomatic complexity
- Considers nested conditions, loops, recursion
- Helps identify code that's hard to maintain

**Configuration**:
```yaml
gocognit:
  min-complexity: 50    # Allows moderate complexity
  skip-tests: true      # Don't check test files
```

**Status**: Enabled with 0 violations
- Threshold of 50 is appropriate for current codebase
- No code needs refactoring to reduce complexity
- Configuration is well-tuned

**When to use**: Run analysis periodically to identify complex functions before they become problematic

## Phase 4b: Dupl Linter ⏳ DISABLED

**What it does**: Detects duplicated code blocks (threshold: 150 lines)
- Identifies opportunities for abstraction
- Helps maintain DRY principle

**Current Status**: 16 violations
- 2 violations in `internal/dao/types_test.go` (similar test tables)
- 4 violations in `internal/layouts/layout_manager.go` (similar layout configs)
- 10 violations in view files (`internal/views/*`) with similar modal dialogs

**Problem**: Exclude-rules path patterns not working properly
- Attempted multiple regex patterns: `.*_test\.go$`, `.*/views/.*`, etc.
- golangci-lint limitation with dupl linter's exclude handling

**Why violations are acceptable**:
- UI views have similar patterns by design (jobs, nodes, partitions, users, etc.)
- Refactoring UI code for DRY would reduce readability
- Test code duplication is acceptable for test clarity

**Future Work**:
- Option 1: Add `//nolint:dupl` directives to duplicate blocks (16 locations)
- Option 2: Investigate alternative exclude-rule syntax
- Option 3: Live with disabled linter (acceptable given intentional duplication)

### Dupl Violations by Category

**Test Files (2 violations)**:
- `internal/dao/types_test.go:146` vs `:212` - Similar test case tables

**Layout Files (4 violations)**:
- `internal/layouts/layout_manager.go` - Standard vs Monitoring layouts (similar structure)

**View Files (10 violations)**:
- Users view modal ↔ Accounts view modal (similar detail display)
- Users filter ↔ Accounts filter (similar UI patterns)
- Job modal ↔ Node modal ↔ Partition modal (similar detail displays)
- QoS modal ↔ Reservation modal (similar detail displays)

## Phase 4d: Cyclop Linter ⏳ DISABLED

**What it does**: Measures cyclomatic complexity - number of decision paths through code
- Higher = harder to test and maintain
- Different from cognitive complexity

**Current Status**: 136 pre-existing violations
- Would require substantial refactoring
- Not worth enabling without architectural review

**Examples of violations**:
- Complex CLI command handlers with many flags
- UI view update logic with many conditional branches
- Configuration parsing with nested conditions
- Test suites with many test cases

**Why disabled**:
- Enabling would immediately show 136 violations
- Would require refactoring functions to reduce complexity
- Complex functions often indicate needed functionality, not bad design
- Not worth the effort without clear benefits

**Configuration needed** (if enabling):
```yaml
cyclop:
  max-complexity: 15-25   # Conservative threshold
  skip-tests: true         # Don't check tests
```

**Future**: Consider for Phase 5 with targeted refactoring strategy

## Phase 4e: Containedctx Linter ⏳ DISABLED

**What it does**: Prevents storing `context.Context` in struct fields
- Recommends: pass context as function parameter
- Detects: `type MyStruct struct { ctx context.Context }`

**Current Status**: 22 pre-existing instances throughout codebase
- `internal/app/app.go` - Application context
- `internal/ssh/session_manager.go` - Session management
- `internal/streaming/types.go` - Streaming coordination
- `plugins/observability/*` - Plugin context management

**Why disabled**:
- This is our intentional architectural pattern
- Context stored in struct fields for state management
- Refactoring would be complex and error-prone
- Current pattern works well and is consistent

**Trade-offs**:
- ✅ Consistent state management throughout app
- ✅ Easier to thread context through multiple layers
- ✅ Clear ownership of context lifecycle
- ❌ Differs from linter recommendation
- ❌ Makes context harder to trace through code

**Future**: Only enable if architecture is significantly refactored to use parameter passing

## Phase 4 Summary

### Linter Status Table

| Linter | Status | Violations | Reason |
|--------|--------|-----------|--------|
| revive | ⏳ Disabled | 2 (package-level - unfixable) | Package-level violations can't be suppressed |
| gocognit | ✅ Enabled | 0 (threshold: 50) | Well-tuned for codebase |
| nolintlint | ✅ Enabled | 0 | Validates linting directives |
| noctx | ✅ Enabled | 0 | Fixed in PR #31 |
| dupl | ⏳ Disabled | 16 (intentional UI duplication) | Exclude-rules not working |
| cyclop | ⏳ Disabled | 136 (requires refactoring) | Not worth enabling now |
| containedctx | ⏳ Disabled | 22 (architectural pattern) | Consistent pattern choice |
| gosec | ⚠️ Configured | 87 (with exclusions) | Phase 2 security audit |

### Total Enabled: 3 linters (revive disabled)
### Total Pre-existing Violations (if all enabled): 460+ (mostly intentional)
### Current Clean Violations: 0 issues ✅
### Known Limitations: 2 package-level revive violations (cannot be suppressed)

## Recommendations for Phase 5+

### High Priority (Critical)
1. **Phase 4a (revive) - URGENT**: Investigate proper solution for package-level violations
   - Research revive issue tracker for similar problems
   - Consider contributing fix to revive or golangci-lint
   - Temporary workaround: Accept 2 package-level violations as known limitations
   - May need to rename packages or disable rule entirely

2. **Phase 4b (dupl)**: Decide between nolint directives or exclude-rules investigation
3. **Phase 4c (gocognit)**: Monitor cognitive complexity, add review to PR process

### Medium Priority
4. **Phase 4d (cyclop)**: Plan refactoring strategy for complex functions
5. **Phase 4e (containedctx)**: Document architectural decision for context management

### Low Priority
6. **Modernization**: Address efaceany, mapsloop, slicescontains suggestions
7. **Performance**: Enable and configure whitespace linting (wsl_v5)

## Implementation Details

See `.golangci.yml` for current configuration and exclude rules.

For backward-compatible type aliases, see:
- `internal/auth/interface.go` - AuthProvider, AuthConfig, AuthManager
- `internal/config/schema.go` - ConfigField, ConfigSchema, ConfigGroup, ConfigTemplate, ConfigValidationResult, ConfigValidator
- `internal/ssh/ssh_client.go` - SSHClient, SSHConfig
- `plugins/observability/*` - Plugin-specific type aliases

## Related Documents

- [LINTING.md](LINTING.md) - Overall linting standards
- [CI_CD_SETUP.md](CI_CD_SETUP.md) - CI/CD linting enforcement
- [PRE_COMMIT_SETUP.md](PRE_COMMIT_SETUP.md) - Pre-commit hook setup
