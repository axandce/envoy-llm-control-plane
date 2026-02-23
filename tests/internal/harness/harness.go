package harness

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/axandce/envoy-llm-control-plane/tests/internal/configenv"
  "github.com/axandce/envoy-llm-control-plane/tests/internal/provider"
)

type Harness struct {
	cfg Config
	p   provider.Provider
}

func New(p provider.Provider) *Harness {
	return &Harness{
		cfg:	configenv.MustLoad[Config](),
		p:		p,
	}
}

func (h *Harness) Run(ctx context.Context, m *testing.M) int {
	if h.p == nil {
		fmt.Fprintln(os.Stderr, "[harness] no provider configured")
		return 1
	}

	readyCtx, cancel := context.WithTimeout(ctx, h.cfg.ReadyTimeout.Duration())
	defer cancel()

	defer func() {
		downCtx, cancelDown := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelDown()

		if err := h.p.Down(downCtx); err != nil {
			fmt.Fprintf(os.Stderr, "[harness] down failed (%s): %v\n", h.p.Name(), err)
		}
	}()

	if err := h.p.Up(readyCtx); err != nil {
		fmt.Fprintf(os.Stderr, "[harness] up failed (%s): %v\n", h.p.Name(), err)
		return 1
	}

	if err := h.p.WaitReady(readyCtx); err != nil {
		fmt.Fprintf(os.Stderr, "[harness] wait-ready failed (%s): %v\n", h.p.Name(), err)
		return 1
	}

	return m.Run()
}
