package provider

import "context"

// Provider is the minimal lifecycle contract required by the harness.
type Provider interface {
  Name() string

  Up(ctx context.Context) error
  WaitReady(ctx context.Context) error
  Down(ctx context.Context) error
}

// --- OPTIONAL CAPABILITIES ---

// StatusProvider is an optional capability: providers that can return a snapshot implement this.
// In v0 we keep return type as `any` to avoid premature normalization.
type StatusProvider interface {
  Status(ctx context.Context) (any, error)
}

// Watcher is an optional capability: providers that can stream environment events
// (unhealthy/exited/restarted/oomkilled) implement this.
// Harness uses this for fail-fast monitoring.
type Watcher interface {
  Watch(ctx context.Context) (<-chan Event, error)
}

// LogProvider is an optional capability: providers that can provide logs implement this.
type LogProvider interface {
  Logs(ctx context.Context, req LogsRequest) (LogStream, error)
}

type LogsRequest struct {
  Target *Target // nil => all
  Follow bool
}

type LogStream interface {
  Close() error
  // TODO (v1): change to io.Reader once we have a concrete shape.
  Reader() any
}

