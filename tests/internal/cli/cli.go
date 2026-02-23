package cli

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type cmdError struct {
	name	string
	args	[]string
	out 	[]byte
	cause	error
}

func (e *cmdError) Error() string {
		return fmt.Sprintf(
			"command failed: %s %s: %v\n%s",
			e.name,
			strings.Join(e.args, " "),
			e.cause,
			strings.TrimSpace(string(e.out)),
		)
}

func (e *cmdError) Unwrap() error { return e.cause }

func Run(ctx context.Context, name string, args ...string) error {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()

	if err != nil {
		return &cmdError{
			name:		name,
			args:		args,
			out:		out,
			cause:	err,
		}
	}
	return nil
}