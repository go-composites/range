package Range

import (
	Array "github.com/go-composites/array/src"
	Error "github.com/go-composites/error/src"
	Number "github.com/go-composites/number/src"
	Result "github.com/go-composites/result/src"
)

/*
Range is a numeric interval composite over Go int64 bounds.

Modelled on Ruby's Range, it supports both an inclusive end (Ruby's "1..5")
and an exclusive end (Ruby's "1...5"), and a non-zero integer step. Its
fallible construction returns a Result.Interface so an invalid step (zero) is a
value rather than a panic. The elements it exposes are go-composites Numbers.

Note: Range is deliberately integer-only — its bounds and step are Go int64
values — which keeps element enumeration and coverage tractable.
*/
type Interface interface {
	Begin() Number.Interface
	End() Number.Interface
	Step() Number.Interface
	ExcludesEnd() bool
	Includes(n int64) bool
	Len() int
	IsEmpty() bool
	Each(fn func(Number.Interface) Result.Interface) Result.Interface
	ToArray() Array.Interface
	IsNull() bool
}

type data struct {
	begin      int64
	end        int64
	step       int64
	excludeEnd bool
}

type Option func(*data)

/*
WithStep is a functional parameter setting the Range step.

The step must be non-zero; New returns an error Result when it is zero.
*/
func WithStep(step int64) Option {
	return func(d *data) {
		d.step = step
	}
}

/*
Exclusive is a functional parameter making the Range exclude its end bound
(Ruby's "1...5").
*/
func Exclusive() Option {
	return func(d *data) {
		d.excludeEnd = true
	}
}

/*
New is the Range constructor.

It yields a Result whose payload is a Range over [begin, end] with a step of 1
and an inclusive end. The end can be excluded with Exclusive() and the step set
with WithStep(). A zero step is invalid and produces a Result carrying an Error
("step cannot be zero") instead of a payload — construction never panics and
never returns nil.

	r1 := Range.New(1, 5)                              // 1..5  step 1
	r2 := Range.New(1, 5, Range.Exclusive())           // 1...5
	r3 := Range.New(0, 10, Range.WithStep(2))          // 0,2,4,6,8,10
*/
func New(begin, end int64, options ...Option) Result.Interface {
	d := &data{
		begin:      begin,
		end:        end,
		step:       1,
		excludeEnd: false,
	}
	for _, opt := range options {
		opt(d)
	}
	if d.step == 0 {
		return Result.New(
			Result.WithError(
				Error.New("step cannot be zero"),
			),
		)
	}
	return Result.New(
		Result.WithPayload(Interface(d)),
	)
}

/*
Begin returns the lower bound of the Range as a Number.
*/
func (d data) Begin() Number.Interface {
	return Number.New(Number.WithInt(d.begin))
}

/*
End returns the upper bound of the Range as a Number.
*/
func (d data) End() Number.Interface {
	return Number.New(Number.WithInt(d.end))
}

/*
Step returns the Range step as a Number.
*/
func (d data) Step() Number.Interface {
	return Number.New(Number.WithInt(d.step))
}

/*
ExcludesEnd reports whether the Range excludes its end bound (Ruby's "1...5").
*/
func (d data) ExcludesEnd() bool {
	return d.excludeEnd
}

// withinBound reports whether v is on the correct side of the end bound,
// honouring the step direction and the inclusive / exclusive flag.
func (d data) withinBound(v int64) bool {
	if d.step > 0 {
		if d.excludeEnd {
			return v < d.end
		}
		return v <= d.end
	}
	if d.excludeEnd {
		return v > d.end
	}
	return v >= d.end
}

/*
Includes reports whether n is an element of the Range — that is, n lies between
the bounds (honouring the inclusive / exclusive end) and is reachable from begin
by a whole number of steps.
*/
func (d data) Includes(n int64) bool {
	if d.step > 0 && n < d.begin {
		return false
	}
	if d.step < 0 && n > d.begin {
		return false
	}
	if !d.withinBound(n) {
		return false
	}
	return (n-d.begin)%d.step == 0
}

/*
Len returns the number of elements in the Range.
*/
func (d data) Len() int {
	count := 0
	for v := d.begin; d.withinBound(v); v += d.step {
		count++
	}
	return count
}

/*
IsEmpty reports whether the Range has no elements.
*/
func (d data) IsEmpty() bool {
	return d.Len() == 0
}

/*
Each iterates over the Range elements as Numbers, calling fn for each. It
short-circuits and returns the first error Result fn produces; otherwise it
returns a success Result once every element has been visited.
*/
func (d data) Each(fn func(Number.Interface) Result.Interface) Result.Interface {
	for v := d.begin; d.withinBound(v); v += d.step {
		if result := fn(Number.New(Number.WithInt(v))); result.HasError() {
			return result
		}
	}
	return Result.New()
}

/*
ToArray materialises the Range elements into an Array of Numbers.
*/
func (d data) ToArray() Array.Interface {
	out := Array.New()
	for v := d.begin; d.withinBound(v); v += d.step {
		out.Push(Number.New(Number.WithInt(v)))
	}
	return out
}

/*
IsNull reports whether the Range is the Null-Object variant.

A concrete Range is never null.
*/
func (d data) IsNull() bool {
	return false
}
