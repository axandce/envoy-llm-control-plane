package harness

import (
  "errors"
  "path/filepath"
)

type Path string

func (p Path) String() string {
  return string(p)
}

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
