package integration

import (
	"context"
	"os"
	"testing"

	"github.com/axandce/envoy-llm-control-plane/tests/internal/harness"
	"github.com/axandce/envoy-llm-control-plane/tests/internal/provider/compose"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	h := harness.New(compose.New())

	os.Exit(h.Run(ctx, m))
}
