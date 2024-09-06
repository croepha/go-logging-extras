package errordump

import (
	"log/slog"
	"reflect"

	"gitlab.com/croepha/common-utils/lostandfound"
)

/*

Some tools to print details about errors beyond just a simple string representation.

*/

// This is a kind of function that can provide details about a given error
// Detailers should return nil, map[string]any or a struct
// nil means that the detailer is not suitable for the given error
type Detailer func(err error) (details Details)
type Details any

// Stops at first one returning non-empty details
func ChainDetailers(detailers ...Detailer) Detailer {
	return func(err error) Details {
		for _, d := range detailers {
			if r := d(err); r != nil {
				return r
			}
		}
		return nil
	}
}

type wrappingDetails struct {
	Error        Details
	WrappedError Details
}
type wrappingManyDetails struct {
	Error         Details
	WrappedErrors []Details
}

func unwrapError(err error, next Detailer) Details {
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		return wrappingDetails{
			Error:        next(err),
			WrappedError: unwrapError(x.Unwrap(), next),
		}
	case interface{ Unwrap() []error }:
		return wrappingManyDetails{
			Error: next(err),
			WrappedErrors: lostandfound.MapApply(
				x.Unwrap(),
				func(err error) Details {
					return unwrapError(err, next)
				},
			),
		}
	default:
		return next(err)
	}
}

func NewUnwrappingDetailer(next Detailer) Detailer {
	return func(err error) Details {
		return unwrapError(err, next)
	}
}

// This should get any JSON Marshallable info from the error
func RawDetailer(err error) Details {
	return err
}

// TODO: Also kinda would be nice to reflect all the struct members, some
// errors have details/codes are hidden away as private members

func ReflectionDetailer(next Detailer) Detailer {
	return func(err error) Details {
		r := reflect.TypeOf(err)
		for r.Kind() == reflect.Pointer {
			r = r.Elem()
		}
		return reflectionDetails{
			String:               err.Error(),
			ReflectedName:        r.Name(),
			ReflectedPackagePath: r.PkgPath(),
			NextDetails:          next(err),
		}
	}
}

type reflectionDetails struct {
	String               string
	ReflectedName        string
	ReflectedPackagePath string
	NextDetails          Details
}

var GlobalDetailer = NewUnwrappingDetailer(ReflectionDetailer(RawDetailer))

type slogValue struct {
	err error
}

func (v *slogValue) LogValue() slog.Value {
	return slog.AnyValue(GlobalDetailer(v.err))
}

func NewSlog(name string, err error) slog.Attr {
	return slog.Any(name, &slogValue{err: err})
}
