package dht_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-identity"

	"time"

	. "github.com/republicprotocol/go-dht"
)

var _ = Describe("Distributed Hash Table", func() {
	var dht *DHT
	var randomAddress identity.Address
	var randomMulti identity.MultiAddress

	Context("updates", func() {
		BeforeEach(func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht = NewDHT(address, 20)

			randomAddress, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())

			randomMulti, err = randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())

		})

		It("should be able to find address after it is updated", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			multi, err := dht.FindMultiAddress(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*multi).Should(Equal(randomMulti))
		})

		It("should error when the bucket is full", func() {
			var err error
			// The maximum number of entries that can be held by a DHT with
			// Bucket size 20.
			for i := 0; i < 160*20+1; i++ {
				address, _, e := identity.NewAddress()
				Ω(e).ShouldNot(HaveOccurred())
				multi, e := address.MultiAddress()
				Ω(e).ShouldNot(HaveOccurred())
				e = dht.Update(multi)
				if err == nil && e != nil {
					err = e
					break
				}
			}
			Ω(err).Should(HaveOccurred())
		})

		It("should not update the time stamp for existing addresses", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			t := bucket.Get(0).Time

			time.Sleep(time.Millisecond)

			err = dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(t).Should(Equal(bucket.Get(0).Time))
			t = bucket.Get(0).Time
		})
	})

	Context("removing addresses", func() {

		BeforeEach(func() {
			address, _, err := identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())
			dht = NewDHT(address, 20)

			randomAddress, _, err = identity.NewAddress()
			Ω(err).ShouldNot(HaveOccurred())

			randomMulti, err = randomAddress.MultiAddress()
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should remove the address", func() {
			err := dht.Update(randomMulti)
			Ω(err).ShouldNot(HaveOccurred())

			bucket, err := dht.FindBucket(randomAddress)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(bucket.Length()).Should(Equal(1))

			dht.Remove(randomMulti)
			Ω(bucket.Length()).Should(Equal(0))
		})
	})
})
