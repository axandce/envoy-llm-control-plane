package integration

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("[integration] TestMain starting (v0: no harness yet)")

	// TODO: v1+ will:
	// - load HARNESS_* env vars
	// - harness.Up()
	// - harness.WaitReady()
	// - run tests
	// - harness.Down() (unless KEEP_STACK_UP)

	code := m.Run()

	fmt.Println("[integration] TestMain finished")
	os.Exit(code)
}
