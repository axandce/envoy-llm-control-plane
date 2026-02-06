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

type Config struct {
  Project		    string        `envDefault:"DefaultProject"`
  RootDir       Path          `envDefault:"../.."`

  ReadyTimeout 	time.Duration `envDefault:"120s"`
  PollEvery     time.Duration `envDefault:"500ms"`

  // DownVolumes 	bool          `envDefault:"false"`
}

type Path string

func parsePath(v string) (any, error) {
  if v == "" { return nil, errors.New("path is empty") }

  // abs, err := filepath.Abs(filepath.Clean(v))
  abs, err := filepath.Abs(v)
  if err != nil { return nil, err }

  fmt.Println("WHOEIFHOIEHFOWIEHFOIWEHf")

  return Path(abs), nil
}


func logOnSet(key string, v any, isDefault bool) {
  fmt.Printf("env %s=%v (isDefault=%v)\n", key, v, isDefault)
}

// func debugEnvironment() map[string]string {
//   return map[string]string{
//     "HARNESS_PROJECT":        "TEST_PROJECT_NAME",
//     "HARNESS_ROOT_DIR":       "../..",

//     "HARNESS_READY_TIMEOUT":  "170s",
//     "HARNESS_POLL_EVERY":     "600ms",

//     // "HARNESS_DOWN_VOLUMES":    "false",
//   }
// }

func configOptions() env.Options {
  return env.Options{
    Prefix:                 "HARNESS_",
    UseFieldNameByDefault:  true,
    OnSet:                  logOnSet,
    // Environment:            debugEnvironment(),

    FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[Path](): parsePath,
		}, 
  }
}

func LoadConfig() (Config, error) {
  return env.ParseAsWithOptions[Config](configOptions())
}

func (c Config) Show() {
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println("Loaded config:")
	fmt.Println(string(b))
}
