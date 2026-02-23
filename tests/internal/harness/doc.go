// Package harness orchestrates environment lifecycle around a test run.
//
// Harness owns:
// - lifecycle ordering (Up/WaitReady/Run/Down)
// - policy (timeouts, fail-fast monitoring, artifact collection)
//
// Harness depends only on tests/internal/provider interfaces.
// Harness must not import concrete providers (e.g., provider/compose).
package harness
