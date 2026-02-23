package harness

type Config struct {
  ReadyTimeout NonNegativeDuration `envDefault:"120s"`
}
