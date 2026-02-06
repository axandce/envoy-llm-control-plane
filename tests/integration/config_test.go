package integration

import (
	"testing"

	"github.com/axandce/envoy-llm-control-plane/tests/internal/harness"
)


func TestPrintConfig(t *testing.T) {
	cfg, err := harness.LoadConfig();
	if err != nil { t.Fatalf("LoadConfig: %v", err)}
	
	cfg.Show()
}
