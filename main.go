package main

import (
	"fmt"

	Number "github.com/go-composites/number/src"
	Range "github.com/go-composites/range/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	// An inclusive range 1..5.
	inclusive := Range.New(1, 5).Payload().(Range.Interface)
	fmt.Println("1..5 length :", inclusive.Len())
	fmt.Println("1..5 has 5  :", inclusive.Includes(5))

	// An exclusive range 1...5 (Ruby's "...").
	exclusive := Range.New(1, 5, Range.Exclusive()).Payload().(Range.Interface)
	fmt.Println("1...5 length:", exclusive.Len())
	fmt.Println("1...5 has 5 :", exclusive.Includes(5))

	// A stepped range 0,2,4,6,8,10.
	stepped := Range.New(0, 10, Range.WithStep(2)).Payload().(Range.Interface)
	fmt.Print("0..10 step 2:")
	stepped.Each(func(n Number.Interface) Result.Interface {
		fmt.Print(" ", n.ToGoString())
		return Result.New()
	})
	fmt.Println()

	// A zero step is a value, not a panic.
	bad := Range.New(0, 10, Range.WithStep(0))
	fmt.Println("step 0 error:", bad.HasError())
	fmt.Println(bad.Error().Message())
}
