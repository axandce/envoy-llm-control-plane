package integration

import (
	"context"
	"testing"
	"time"

	"github.com/axandce/envoy-llm-control-plane/tests/internal/harness"
)


func TestPrintConfig(t *testing.T) {
	h, err := harness.LoadHarness();
	if err != nil { t.Fatalf("LoadHarness: %v", err)}
	
	h.Show()
}

func TestHarnessUpDown(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	h, err := harness.LoadHarness()
	if err != nil {
		t.Fatalf("LoadHarness: %v", err)
	}

	h.WithRunner(harness.ExecRunner{})

	t.Log("--- Attempting Up....")

	if _, err := h.Up(ctx); err != nil {
		t.Fatalf("Up: %v", err)
	}
	t.Log("--- Up happened.  Sleeping...")

	t.Log("--- zzz ")
	// Just wait.
	time.Sleep(5 * time.Second)

	t.Log("--- Awake!  Attempting Logs...")

	if logs, err := h.Logs(ctx); err == nil {
		t.Logf("compose logs:\n%s", logs)
	}

	t.Log("--- Logs happened.  Attempting Down...")

	if _, err := h.Down(ctx); err != nil {
		t.Fatalf("Down: %v", err)
	}

	t.Log("--- Down happened.")
}