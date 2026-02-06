package harness

import (
	"encoding/json"
  "errors"
  "fmt"
  "path/filepath"
  "reflect"
  "time"
  "github.com/caarlos0/env/v11"
)

func LoadConfig() (Config, error) {
  return env.ParseAsWithOptions[Config](configOptions())
}

func (c Config) Show() {
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("Loaded config:")
	fmt.Println(string(b))
  fmt.Println("Getting ComposePath...")
  fmt.Println(string(c.ComposePath()))
}

func (c Config) ComposePath() Path {
  if filepath.IsAbs(c.ComposeFile) { 
    return Path(c.ComposeFile) 
  }

  return c.RootDir.Join(c.ComposeFile)
}

type Config struct {
  Project		    string        `envDefault:"DefaultProject"`
  RootDir       Path          `envDefault:"../.."`
  ComposeFile   string        `envDefault:"docker-compose.yml"`
  ReadyTimeout 	time.Duration `envDefault:"120s"`
  PollEvery     time.Duration `envDefault:"500ms"`
}

func configOptions() env.Options {
  return env.Options{
    Prefix:                 "HARNESS_",
    UseFieldNameByDefault:  true,
    
    OnSet:                  logOnSet,
    // Environment:            DebugEnvironment(),

    FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[Path](): parsePath,
		},
  }
}

type Path string

func (p Path) Join(elem ...string) Path {
	parts := append([]string{string(p)}, elem...)
	return Path(filepath.Join(parts...))
}

func parsePath(v string) (any, error) {
  if v == "" { 
    return nil, errors.New("path is empty")
  }

  abs, err := filepath.Abs(filepath.Clean(v))
  if err != nil { 
    return nil, err
  }

  return Path(abs), nil
}

func logOnSet(key string, v any, isDefault bool) {
  fmt.Printf("env %s=%v (isDefault=%v)\n", key, v, isDefault)
}

func DebugEnvironment() map[string]string {
  return map[string]string{
    "HARNESS_PROJECT":        "TEST_PROJECT_NAME",
    "HARNESS_ROOT_DIR":       "../..",

    "HARNESS_READY_TIMEOUT":  "170s",
    "HARNESS_POLL_EVERY":     "600ms",
  }
}
