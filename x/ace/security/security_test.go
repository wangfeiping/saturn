package security_test

import (
	"fmt"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/wangfeiping/saturn/x/ace/security"
	"github.com/wangfeiping/saturn/x/ace/security/paillier"
)

var _ = Describe("x/ace/security", func() {

	var (
		privkey security.PrivateKey
		pubkey  security.PublicKey
		plainA  []byte
		plainB  []byte
		plainC  []byte
		result  []byte
	)

	BeforeEach(func() {
	})

	Describe("Encryption operations with paillier keys", func() {
		Context("create a pair of keys", func() {
			It("should be success", func() {
				privkey = paillier.Create()
				pubkey = privkey.PublicKey()
				Expect(privkey).ShouldNot(BeNil())
			})
		})

		Context("Encrypt the number 7011, 1058, 39099", func() {
			It("should be success", func() {
				var err error
				plainA, err = pubkey.Encrypt(big.NewInt(7011).Bytes())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(plainA)).NotTo(BeZero())
				plainB, err = pubkey.Encrypt(big.NewInt(1058).Bytes())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(plainB)).NotTo(BeZero())
				plainC, err = pubkey.Encrypt(big.NewInt(39099).Bytes())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(plainC)).NotTo(BeZero())
			})
		})

		Context("Add up the three numbers.", func() {
			It("should be success", func() {
				tmp, err := pubkey.Add(plainA, plainB)
				Expect(err).ShouldNot(HaveOccurred())
				result, err = pubkey.Add(tmp, plainC)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("Decrypt the result 39099", func() {
			It("should be success", func() {
				r, err := privkey.Decrypt(result)
				Expect(err).ShouldNot(HaveOccurred())
				i := new(big.Int).SetBytes(r)
				Expect(i).To(Equal(big.NewInt(47168)))
			})
		})
	})
})

var _ = Describe("Paillier", func() {

	var (
		cdc     *codec.Codec
		privkey security.PrivateKey
		pubkey  security.PublicKey
	)

	cdc = codec.New()
	cdc.RegisterInterface((*security.PrivateKey)(nil), nil)
	cdc.RegisterInterface((*security.PublicKey)(nil), nil)
	cdc.RegisterConcrete(paillier.PaillierPrivKey{}, "Saturn/PrivateKey/Paillier", nil)
	cdc.RegisterConcrete(paillier.PaillierPubKey{}, "Saturn/PublicKey/Paillier", nil)
	cdc.Seal()

	BeforeEach(func() {

	})

	Describe("Codec encode and decode the security keys", func() {
		Context("encode", func() {
			privkey = paillier.Create()
			pubkey = privkey.PublicKey()

			if cdc == nil {
				fmt.Println("nil!!!")
			}
			bytes := cdc.MustMarshalJSON(privkey)
			// if err != nil {
			// 	fmt.Println("err: ", err)
			// }
			// Expect(err).ShouldNot(HaveOccurred())
			fmt.Println("privkey: ", string(bytes))

			bytes = cdc.MustMarshalJSON(pubkey)
			fmt.Println("pubkey: ", string(bytes))

			var pub security.PublicKey
			cdc.MustUnmarshalJSON(bytes, &pub)
			_, _ = pub.Encrypt(big.NewInt(7011).Bytes())

		})
	})
})
