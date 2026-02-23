package provider

type TargetKind string

const (
  TargetService   TargetKind = "service" // logical component (compose service, k8s workload, proc group)
  TargetUnit      TargetKind = "unit"    // runtime instance (container, pod, process)
)

type Target struct {
  Kind 			      TargetKind
  Name            string // stable logical name (e.g. "envoy")
  ID              string // optional backend identifier (e.g. container id, pod uid, pid)
}

type EventType string

const (
  EventReady      EventType = "ready"
  EventHealthy    EventType = "healthy"
  EventUnhealthy  EventType = "unhealthy"
  EventExited     EventType = "exited"
  EventRestarted  EventType = "restarted"
  EventOOMKilled  EventType = "oom_killed"
  EventWarning    EventType = "warning"
)

type Event struct {
  Type            EventType
  Target          Target
  AtUnixMs        int64
  Message         string
  Fields          map[string]string // exitCode, reason, restartCount, etc.
}
