package xerr

import (
	"errors"
	"fmt"
	"testing"

	gferrors "github.com/go-faster/errors"
	"github.com/stretchr/testify/require"
)

func TestXErr(t *testing.T) {
	errtype := errors.New("ErrorType")
	type DetailA int
	type DetailB int
	type DetailC int

	base := errors.New("postgress error")
	err := Wrap(base, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))
	err = Wrap(err, WithDetails(DetailB(23)))
	err = Wrap(err, WithType(errtype))
	err = Wrap(err, WithDetails(DetailC(34)))
	err = Wrap(err, WithDetails(DetailA(45)))

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(45))
	detailB := ExtractDetails[DetailB](err)
	require.Equal(t, *detailB, DetailB(23))
	detailC := ExtractDetails[DetailC](err)
	require.Equal(t, *detailC, DetailC(34))

	require.ErrorIs(t, err, base)
}

func TestXErr_MultiStack(t *testing.T) {
	errtype := errors.New("ErrorType")
	type DetailA int

	err := errors.New("postgress error")
	err = Wrap(err, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))
	err = Wrap(err, WithStack())
	err = Wrap(err, WithType(errtype))
	err = Wrap(err, WithStack())

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(12))
}

func TestXErr_Join(t *testing.T) {
	errtype := errors.New("ErrorType")
	type DetailA int

	errA := errors.New("postgress error")
	errA = Wrap(errA, WithStack())

	errB := errors.New("schmostgress error")
	errB = Wrap(errB, WithStack())

	err := Join(errA, errB)

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	err = Wrap(err, WithType(errtype))

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
}

func TestXErr_New(t *testing.T) {
	type DetailA int

	err := New(fmt.Sprintf("postgress error: %d-%d", 1, 2))
	err = Wrap(err, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(12))
}

func TestXErr_NewNilCall(t *testing.T) {
	type DetailA int

	err := NilCall()
	err = Wrap(err, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(12))
}

func TestXErr_NewNilArg(t *testing.T) {
	type DetailA int

	err := NilArg("nil arg name")
	err = Wrap(err, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(12))
}

func TestXErr_GoFaster(t *testing.T) {
	type DetailA int
	err := NilArg("nil arg name")
	err = Wrap(err, WithStack())
	err = Wrap(err, WithDetails(DetailA(12)))
	err = gferrors.Wrap(err, "with go-faster info")

	t.Log(err)
	t.Log(err.Error())
	t.Logf("%v", err)
	t.Logf("%+v", err)
	t.Log(Render(err))

	detailA := ExtractDetails[DetailA](err)
	require.Equal(t, *detailA, DetailA(12))
}
