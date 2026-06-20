<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/range" width="720"></p>

# range

[![ci](https://github.com/go-composites/range/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/range/actions/workflows/ci.yml)

A numeric **Range** composite for **Composition-Oriented Programming**, modelled
on Ruby's `Range`. A `Range` is an integer interval with a non-zero step that
supports both an **inclusive** end (Ruby's `1..5`) and an **exclusive** end
(Ruby's `1...5`). Construction is **fallible and returns a `Result`** — an
invalid step (zero) is a *value*, never a panic and never `nil`.

The elements of a `Range` are exposed as go-composites `Number`s:

```golang
inclusive := Range.New(1, 5).Payload().(Range.Interface)
inclusive.Each(func(n Number.Interface) Result.Interface {
    fmt.Print(" ", n.ToGoString()) // 1 2 3 4 5
    return Result.New()
})
```

> Range is deliberately **integer-only** (its bounds and step are Go `int64`),
> which keeps element enumeration and 100% coverage tractable.

`Range` follows the org's Null-Object / never-nil invariant (enforced by the
`nonnil` CI analyzer): the `NullRange` variant in `src/null` satisfies the same
`Interface`, is an empty range, and reports `IsNull() == true`.

## Install

```bash
export GOPRIVATE=github.com/go-composites GOPROXY=direct GOSUMDB=off
go get github.com/go-composites/range@main
```

## Usage

> [!NOTE] main.go

```golang
package main

import (
    "fmt"

    Number "github.com/go-composites/number/src"
    Range "github.com/go-composites/range/src"
    Result "github.com/go-composites/result/src"
)

func main() {
    inclusive := Range.New(1, 5).Payload().(Range.Interface)
    fmt.Println("1..5 length :", inclusive.Len())   // 5
    fmt.Println("1..5 has 5  :", inclusive.Includes(5)) // true

    exclusive := Range.New(1, 5, Range.Exclusive()).Payload().(Range.Interface)
    fmt.Println("1...5 has 5 :", exclusive.Includes(5)) // false

    bad := Range.New(0, 10, Range.WithStep(0))
    fmt.Println("step 0 error:", bad.HasError())   // true
    fmt.Println(bad.Error().Message())             // step cannot be zero
}
```

```bash
$ task build
```

## API

Constructor

- `New(begin, end int64, options ...Option) Result.Interface` — payload is a
  `Range`, or an `Error` (`"step cannot be zero"`) when the step is zero.
- `WithStep(int64) Option` — set the (non-zero) step (default `1`).
- `Exclusive() Option` — exclude the end bound (Ruby's `1...5`).
- `null.New() Interface` — the `NullRange` Null-Object (`IsNull() == true`).

Bounds (each returns `Number.Interface`)

- `Begin()`, `End()`, `Step()`.
- `ExcludesEnd() bool` — whether the end bound is excluded.

Membership and size

- `Includes(n int64) bool` — honours the step and the inclusive / exclusive end.
- `Len() int`, `IsEmpty() bool`.

Iteration and materialisation

- `Each(fn func(Number.Interface) Result.Interface) Result.Interface` —
  short-circuits on the first error `Result`.
- `ToArray() Array.Interface` — materialise the elements as `Number`s.

Inspection

- `IsNull() bool`.

## License

BSD-3-Clause — see [LICENSE](./LICENSE).
