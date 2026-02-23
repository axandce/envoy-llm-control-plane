package configenv

import "github.com/caarlos0/env/v11"

func DefaultOptions() env.Options {
	return env.Options{
		Prefix:                "HARNESS_",
		UseFieldNameByDefault: true,
	}
}

func Load[T any]() (T, error) {
	return env.ParseAsWithOptions[T](DefaultOptions())
}

func MustLoad[T any]() T {
	return env.Must(Load[T]())
}

func LoadWith[T any](opts env.Options) (T, error) {
	return env.ParseAsWithOptions[T](opts)
}

func MustLoadWith[T any](opts env.Options) T {
	return env.Must(LoadWith[T](opts))
}
