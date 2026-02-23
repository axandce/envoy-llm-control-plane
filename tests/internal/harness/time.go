package harness

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type NonNegativeDuration time.Duration

func (d NonNegativeDuration) Duration() time.Duration {
	return time.Duration(d)
}

func (d NonNegativeDuration) SecondsCeil() int {
	return int(math.Ceil(d.Duration().Seconds()))
}

func (d NonNegativeDuration) SecondsCeilStr() string {
	return strconv.Itoa(d.SecondsCeil())
}

func (d *NonNegativeDuration) UnmarshalText(text []byte) error {
	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	if parsed < 0 {
		return fmt.Errorf("duration must be non-negative")
	}
	*d = NonNegativeDuration(parsed)
	return nil
}
