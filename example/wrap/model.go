package wrap

import (
	"context"
	"errors"
)

//go:generate go run github.com/AugustineAurelius/eos/ generator wrapper  --name Test
type Test struct {
}

func (t *Test) Test1(a int, b float64) (int, error) {
	return a + int(b), errors.New("123")
}

func (t *Test) Test2(a int, b float64) error {
	return nil
}

func (t *Test) Test3(ctx context.Context, a int, b float64) error {
	return nil
}

func (t *Test) testPriv(ctx context.Context, a int, b float64) error {
	return nil
}
