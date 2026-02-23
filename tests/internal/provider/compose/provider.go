package compose

import (
	"context"
	
	"github.com/axandce/envoy-llm-control-plane/tests/internal/configenv"
)

type Provider struct {
	cfg	Config
	cli	*dockerCLI
}

func New() *Provider {
	return &Provider{
		cfg: configenv.MustLoad[Config](),
		cli: newDockerCLI(),
	}
}

func (p *Provider) Name() string { return "compose" }

func (p *Provider) Up(ctx context.Context) error {
	return p.run(ctx, "up", "-d", "--wait")
}

func (p *Provider) WaitReady(ctx context.Context) error {
	return nil // V0: Up blocks with --wait
}

func (p *Provider) Down(ctx context.Context) error {
	return p.run(ctx, "down", "--remove-orphans")
}

func (p *Provider) run(ctx context.Context, args ...string) error {
	return p.cli.run(ctx, append(p.cfg.ComposeArgs(), args...)...)
}
