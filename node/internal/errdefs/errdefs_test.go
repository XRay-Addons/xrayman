package errdefs

import (
	"fmt"
	"testing"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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
		return With(anonFunc(), "detail A")
	}()
}

func (i *anonInner) privatePrivateMethod() error {
	return WithStack(ErrCmdExec)
}
func (i *anonInner) PublicMethod() error {
	return i.privateMethod()
}

type Outer struct{}

func (o Outer) DoSomething() error {
	fn := func() error {
		inner := &anonInner{Name: "anonInner"}
		return With(inner.PublicMethod(), "detail B")
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
		return With(ErrConfig, "")
	}

	err = topAnon()
	if err != nil {
		log.Error("test", zap.Error(err))
	}
}

func TestErrorAs(t *testing.T) {
	f := func() error {
		return func() error {
			return WithStack(ErrFileAccess)
		}()
	}()
	var oe OriginError
	require.True(t, errors.As(f, &oe))
	fmt.Printf("%+v\n", oe)
}

func TestErrorIs(t *testing.T) {
	f := func() error {
		return func() error {
			return WithStack(ErrFileAccess)
		}()
	}()
	require.True(t, errors.Is(f, ErrFileAccess))
}
