package NullRange

import (
	Array "github.com/go-composites/array/src"
	Number "github.com/go-composites/number/src"
	Range "github.com/go-composites/range/src"
	Result "github.com/go-composites/result/src"
)

/*
NullRange is the Null-Object variant of Range.

It satisfies Range.Interface so callers never have to test for a bare nil: it is
an empty range — its bounds and step are the zero Number, it contains nothing,
Each is a no-op success, ToArray yields an empty Array, Includes is always
false, and IsNull() returns true.
*/
type Interface interface {
	Range.Interface
}

type data struct{}

/*
New returns a NullRange.
*/
func New() Interface {
	return &data{}
}

func (d data) Begin() Number.Interface {
	return Number.New()
}

func (d data) End() Number.Interface {
	return Number.New()
}

func (d data) Step() Number.Interface {
	return Number.New()
}

func (d data) ExcludesEnd() bool {
	return false
}

func (d data) Includes(int64) bool {
	return false
}

func (d data) Len() int {
	return 0
}

func (d data) IsEmpty() bool {
	return true
}

func (d data) Each(func(Number.Interface) Result.Interface) Result.Interface {
	return Result.New()
}

func (d data) ToArray() Array.Interface {
	return Array.New()
}

func (d data) IsNull() bool {
	return true
}
