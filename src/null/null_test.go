package NullRange_test

import (
	Number "github.com/go-composites/number/src"
	Range "github.com/go-composites/range/src"
	NullRange "github.com/go-composites/range/src/null"
	Result "github.com/go-composites/result/src"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("NullRange", func() {
	var n NullRange.Interface
	ginkgo.BeforeEach(func() {
		n = NullRange.New()
	})

	ginkgo.It("satisfies the Range interface", func() {
		var _ Range.Interface = n
	})
	ginkgo.It("reports IsNull() true", func() {
		gomega.Expect(n.IsNull()).To(gomega.BeTrue())
	})
	ginkgo.It("is an empty range", func() {
		gomega.Expect(n.Len()).To(gomega.Equal(0))
		gomega.Expect(n.IsEmpty()).To(gomega.BeTrue())
	})
	ginkgo.It("exposes zero Numbers for its bounds and step", func() {
		gomega.Expect(n.Begin().ToGoInt()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.End().ToGoInt()).To(gomega.BeEquivalentTo(0))
		gomega.Expect(n.Step().ToGoInt()).To(gomega.BeEquivalentTo(0))
	})
	ginkgo.It("does not exclude an end", func() {
		gomega.Expect(n.ExcludesEnd()).To(gomega.BeFalse())
	})
	ginkgo.It("includes nothing", func() {
		gomega.Expect(n.Includes(0)).To(gomega.BeFalse())
	})
	ginkgo.It("Each is a no-op success and never calls fn", func() {
		called := false
		r := n.Each(func(Number.Interface) Result.Interface {
			called = true
			return Result.New()
		})
		gomega.Expect(r.HasError()).To(gomega.BeFalse())
		gomega.Expect(called).To(gomega.BeFalse())
	})
	ginkgo.It("ToArray yields an empty Array", func() {
		gomega.Expect(n.ToArray().IsEmpty()).To(gomega.BeTrue())
	})
})
