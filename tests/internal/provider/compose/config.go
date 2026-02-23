package compose

import (
  "path/filepath"
  
  "github.com/axandce/envoy-llm-control-plane/tests/internal/provider"
)

type Config struct {
  // Project     string          `env:"required"`
  // Project     string          `envRequired:"true"`
  Project     string          `envDefault:"envoy-llm-control-plane"`
  ProjectDir  provider.Path   `envDefault:"../.."`
  ComposeFile string          `envDefault:"docker-compose.yml"`
}

func (c Config) ComposePath() provider.Path {
  if filepath.IsAbs(c.ComposeFile) { 
    return provider.Path(c.ComposeFile) 
  }

  return c.ProjectDir.Join(c.ComposeFile)
}

func (c Config) ComposeArgs() []string {
	return []string{
		"compose",
		"--project-directory", c.ProjectDir.String(),
		"-p", c.Project,
		"-f", c.ComposePath().String(),
	}
}
