package provider

import (
  "errors"
  "path/filepath"
)

type Path string

func (p Path) String() string {
  return string(p)
}

func (p Path) Join(elem ...string) Path {
  parts := append([]string{p.String()}, elem...)
  return Path(filepath.Join(parts...))
}

func ParsePath(v string) (Path, error) {
  var p Path

  if v == "" { 
    return p, errors.New("path is empty")
  }

  abs, err := filepath.Abs(filepath.Clean(v))
  if err != nil { 
    return p, err
  }

  return Path(abs), nil
}

func (p *Path) UnmarshalText(text []byte) error {
	pp, err := ParsePath(string(text))
	if err != nil {
		return err
	}
	*p = pp
	return nil
}
