package compose

import (
	"context"
	
	"github.com/axandce/envoy-llm-control-plane/tests/internal/cli"
)

type dockerCLI struct{}

func newDockerCLI() *dockerCLI { return &dockerCLI{} }

func (d *dockerCLI ) run(ctx context.Context, args ...string) error {
	return cli.Run(ctx, "docker", args...)
}
