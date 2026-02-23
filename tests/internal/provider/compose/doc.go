// Package compose implements provider contracts using Docker Compose.
//
// V0 semantic choice:
//   - Up() uses `docker compose up --wait` and therefore blocks until readiness.
//   - WaitReady() is a documented no-op in V0 for this provider.
//
// Later hardening may split Up() to `up -d` and move readiness into WaitReady().
package compose
