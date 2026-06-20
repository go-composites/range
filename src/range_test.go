package Range_test

import (
	Error "github.com/go-composites/error/src"
	Number "github.com/go-composites/number/src"
	Range "github.com/go-composites/range/src"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// payload extracts the Range from a successful constructor Result.
func payload(r Result.Interface) Range.Interface {
	gomega.ExpectWithOffset(1, r.HasError()).To(gomega.BeFalse())
	return r.Payload().(Range.Interface)
}

var _ = ginkgo.Describe("Range", func() {

	ginkgo.Describe("construction", func() {
		ginkgo.It("builds an inclusive range with a default step of 1", func() {
			r := payload(Range.New(1, 5))
			gomega.Expect(r.Begin().ToGoInt()).To(gomega.BeEquivalentTo(1))
			gomega.Expect(r.End().ToGoInt()).To(gomega.BeEquivalentTo(5))
			gomega.Expect(r.Step().ToGoInt()).To(gomega.BeEquivalentTo(1))
			gomega.Expect(r.ExcludesEnd()).To(gomega.BeFalse())
			gomega.Expect(r.IsNull()).To(gomega.BeFalse())
		})
		ginkgo.It("builds an exclusive range with Exclusive()", func() {
			r := payload(Range.New(1, 5, Range.Exclusive()))
			gomega.Expect(r.ExcludesEnd()).To(gomega.BeTrue())
		})
		ginkgo.It("builds a stepped range with WithStep()", func() {
			r := payload(Range.New(0, 10, Range.WithStep(2)))
			gomega.Expect(r.Step().ToGoInt()).To(gomega.BeEquivalentTo(2))
		})
		ginkgo.It("returns an error Result for a zero step", func() {
			r := Range.New(0, 10, Range.WithStep(0))
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.Equal("step cannot be zero"))
			var _ Error.Interface = r.Error()
		})
	})

	ginkgo.Describe("length and emptiness", func() {
		ginkgo.It("counts an inclusive range", func() {
			gomega.Expect(payload(Range.New(1, 5)).Len()).To(gomega.Equal(5))
		})
		ginkgo.It("counts an exclusive range", func() {
			gomega.Expect(payload(Range.New(1, 5, Range.Exclusive())).Len()).To(gomega.Equal(4))
		})
		ginkgo.It("counts a stepped range", func() {
			gomega.Expect(payload(Range.New(0, 10, Range.WithStep(2))).Len()).To(gomega.Equal(6))
		})
		ginkgo.It("counts a descending range", func() {
			gomega.Expect(payload(Range.New(5, 1, Range.WithStep(-1))).Len()).To(gomega.Equal(5))
		})
		ginkgo.It("counts a descending exclusive range", func() {
			gomega.Expect(payload(Range.New(5, 1, Range.WithStep(-2), Range.Exclusive())).Len()).To(gomega.Equal(2))
		})
		ginkgo.It("reports a non-empty range", func() {
			gomega.Expect(payload(Range.New(1, 5)).IsEmpty()).To(gomega.BeFalse())
		})
		ginkgo.It("reports an empty range (begin past end)", func() {
			r := payload(Range.New(5, 1))
			gomega.Expect(r.Len()).To(gomega.Equal(0))
			gomega.Expect(r.IsEmpty()).To(gomega.BeTrue())
		})
		ginkgo.It("reports an empty exclusive single-point range", func() {
			r := payload(Range.New(3, 3, Range.Exclusive()))
			gomega.Expect(r.IsEmpty()).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Includes", func() {
		ginkgo.It("includes a member of an inclusive range", func() {
			gomega.Expect(payload(Range.New(1, 5)).Includes(5)).To(gomega.BeTrue())
		})
		ginkgo.It("excludes the end of an exclusive range", func() {
			gomega.Expect(payload(Range.New(1, 5, Range.Exclusive())).Includes(5)).To(gomega.BeFalse())
		})
		ginkgo.It("excludes a value below begin", func() {
			gomega.Expect(payload(Range.New(1, 5)).Includes(0)).To(gomega.BeFalse())
		})
		ginkgo.It("excludes a value above the inclusive end", func() {
			gomega.Expect(payload(Range.New(1, 5)).Includes(6)).To(gomega.BeFalse())
		})
		ginkgo.It("respects the step for a hit", func() {
			gomega.Expect(payload(Range.New(0, 10, Range.WithStep(2))).Includes(4)).To(gomega.BeTrue())
		})
		ginkgo.It("respects the step for a miss", func() {
			gomega.Expect(payload(Range.New(0, 10, Range.WithStep(2))).Includes(5)).To(gomega.BeFalse())
		})
		ginkgo.It("includes within a descending range", func() {
			gomega.Expect(payload(Range.New(5, 1, Range.WithStep(-1))).Includes(3)).To(gomega.BeTrue())
		})
		ginkgo.It("excludes a value above begin in a descending range", func() {
			gomega.Expect(payload(Range.New(5, 1, Range.WithStep(-1))).Includes(6)).To(gomega.BeFalse())
		})
		ginkgo.It("excludes a value below the descending end", func() {
			gomega.Expect(payload(Range.New(5, 1, Range.WithStep(-1))).Includes(0)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("Each", func() {
		ginkgo.It("visits every element in order", func() {
			var seen []int64
			r := payload(Range.New(1, 3)).Each(func(n Number.Interface) Result.Interface {
				seen = append(seen, n.ToGoInt())
				return Result.New()
			})
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(seen).To(gomega.Equal([]int64{1, 2, 3}))
		})
		ginkgo.It("short-circuits on the first error Result", func() {
			var seen []int64
			r := payload(Range.New(1, 5)).Each(func(n Number.Interface) Result.Interface {
				seen = append(seen, n.ToGoInt())
				if n.ToGoInt() == 3 {
					return Result.New(Result.WithError(Error.New("stop")))
				}
				return Result.New()
			})
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).To(gomega.Equal("stop"))
			gomega.Expect(seen).To(gomega.Equal([]int64{1, 2, 3}))
		})
	})

	ginkgo.Describe("ToArray", func() {
		ginkgo.It("materialises the elements as Numbers", func() {
			arr := payload(Range.New(1, 3)).ToArray()
			gomega.Expect(arr.Len()).To(gomega.Equal(3))
			gomega.Expect(arr.First().Payload().(Number.Interface).ToGoInt()).To(gomega.BeEquivalentTo(1))
			gomega.Expect(arr.Last().Payload().(Number.Interface).ToGoInt()).To(gomega.BeEquivalentTo(3))
		})
		ginkgo.It("materialises an empty range to an empty Array", func() {
			arr := payload(Range.New(5, 1)).ToArray()
			gomega.Expect(arr.IsEmpty()).To(gomega.BeTrue())
		})
	})
})
