package harness

import (
	"context"
  "encoding/json"
  "errors"
	"fmt"
  "path/filepath"
  "reflect"
  "time"
  "github.com/caarlos0/env/v11"
)

type Harness struct {
  Project		    string        `envDefault:"DefaultProject"`
  RootDir       Path          `envDefault:"../.."`
  ComposeFile   string        `envDefault:"docker-compose.yml"`
  ReadyTimeout 	time.Duration `envDefault:"120s"`
  PollEvery     time.Duration `envDefault:"500ms"`

  runner        Runner         `env:"-"`
}

func LoadHarness() (Harness, error) {
  return env.ParseAsWithOptions[Harness](harnessOptions())
}

func logOnSet(key string, v any, isDefault bool) {
  fmt.Printf("env %s=%v (isDefault=%v)\n", key, v, isDefault)
}

func harnessOptions() env.Options {
  return env.Options{
    Prefix:                 "HARNESS_",
    UseFieldNameByDefault:  true,
    
    OnSet:                  logOnSet,

    FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[Path](): parsePath,
		},
  }
}

func (h *Harness) WithRunner(runner Runner) *Harness {
  h.runner = runner
  return h
}

func (h Harness) ComposePath() Path {
  if filepath.IsAbs(h.ComposeFile) { 
    return Path(h.ComposeFile) 
  }

  return h.RootDir.Join(h.ComposeFile)
}


func (h Harness) baseComposeArgs() []string {
	return []string{
    "compose",
    "-p", h.Project,
    "--project-directory", h.RootDir.String(),
    "-f", h.ComposePath().String(),
  }
}

func (h Harness) runDockerCompose(ctx context.Context, args ...string) (string, error) {
  if h.runner == nil { return "", errors.New("Runner not set")}
  return h.runner.Run(ctx, "docker", append(h.baseComposeArgs(), args...)...)
}

func (h Harness) Up(ctx context.Context) (string, error) {
  return h.runDockerCompose(ctx, "up", "-d")
}

func (h Harness) Down(ctx context.Context) (string, error) {
  return h.runDockerCompose(ctx, "down", "-v")
}

func (h Harness) Logs(ctx context.Context) (string, error) {
  return h.runDockerCompose(ctx, "logs")
}

func (h Harness) Show() {
	b, _ := json.MarshalIndent(h, "", "  ")

	fmt.Println("Loaded config:")
	fmt.Println(string(b))

  fmt.Println("Getting ComposePath...")
  fmt.Println(string(h.ComposePath()))
}