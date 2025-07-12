package wrap

import (
	"context"
	"errors"
)

//go:generate go run github.com/AugustineAurelius/eos/ generator wrapper  --name Test --logging --tracing
type Test struct {
}

type Test222 struct {
	Name string
}

func (t *Test) Test1(a int, b *Test222) (int, error) {
	return 1, errors.New("123")
}

func (t *Test) Test2(a int, b float64) error {
	return nil
}

func (t *Test) Test3(ctx context.Context, c int, b float64) error {
	return nil
}
func (t *Test) Test5(ctx context.Context, a int, b float64) (int, error) {
	return 0, nil
}

func (t *Test) testPriv(ctx context.Context, a int, b float64) error {
	return nil
}
