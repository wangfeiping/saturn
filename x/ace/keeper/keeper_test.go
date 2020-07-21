package keeper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"

	acekeeper "github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var _ = Describe("x/ace/keeper", func() {

	var (
		play types.Play = types.Play{
			AceID:   "ace_test",
			Height:  1000,
			Address: "xxx",
			Card:    0}

		ctx    sdk.Context
		keeper acekeeper.AceKeeper
	)

	ctx, keeper = CreateMockAceKeeper()

	BeforeEach(func() {

	})

	Describe("Create an ace keeper", func() {
		Context("with mem db", func() {
			It("should be success", func() {
				Expect(keeper).NotTo(BeZero())
				log := keeper.Logger(ctx)
				log.Debug("test")
				// Expect(true).To(BeTrue())
				Expect(log).NotTo(BeZero())
			})

			It(`should be success for calling method "Set"`, func() {
				keeper.Set(ctx, "test", play)
				Expect(keeper.Has(ctx, "test")).To(BeTrue())
			})

			It(`should be success for calling method "Get"`, func() {
				// keeper.Set(ctx, "test", play)
				var p types.Play
				err := keeper.Get(ctx, "test", &p)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(p.AceID).To(Equal(play.AceID))
				Expect(p.Address).To(Equal(play.Address))
			})

			It(`should be success for calling method "Delete"`, func() {
				keeper.Delete(ctx, "test")
				var p types.Play
				err := keeper.Get(ctx, "test", &p)
				Expect(err).Should(HaveOccurred())
				Expect(p).To(BeZero())
			})
		})

	})
})
