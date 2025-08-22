package errdefs

import (
	"fmt"
	"testing"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testError int

const (
	ErrTestConfig testError = iota + 1
	ErrTestFileAccess
)

func (t testError) Error() string {
	return fmt.Sprintf("test error %d", t)
}

type anonInner struct {
	Name string
}

func (i *anonInner) privateMethod() error {
	anonFunc := func() error {
		return func() error {
			return i.privatePrivateMethod()
		}()
	}
	return func() error {
		return WrapWith(anonFunc(), "detail A")
	}()
}

func (i *anonInner) privatePrivateMethod() error {
	return WrapWithStack(ErrTestConfig)
}
func (i *anonInner) PublicMethod() error {
	return i.privateMethod()
}

type Outer struct{}

func (o Outer) DoSomething() error {
	fn := func() error {
		inner := &anonInner{Name: "anonInner"}
		return WrapWith(inner.PublicMethod(), "detail B")
	}
	return fn()
}

func Fn() error {
	return func() error {
		o := &Outer{}
		err := o.DoSomething()
		return err
	}()
}

func getFn() func() error {
	return func() error {
		return Fn()
	}
}

func TestErrorPrinting(t *testing.T) {
	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	f := getFn()
	err = f()
	if err != nil {
		log.Error("test", zap.Error(err))
	}

	topAnon := func() error {
		return New("err config", With("nothing"))
	}

	err = topAnon()
	if err != nil {
		log.Error("test", zap.Error(err))
	}
}

func TestErrorAs(t *testing.T) {
	f := func() error {
		return func() error {
			return WrapWithStack(ErrTestFileAccess)
		}()
	}()
	var te testError
	require.True(t, errors.As(f, &te))
}

func TestErrorIs(t *testing.T) {
	f := func() error {
		return func() error {
			return WrapWithStack(ErrTestFileAccess)
		}()
	}()
	require.True(t, errors.Is(f, ErrTestFileAccess))
}
