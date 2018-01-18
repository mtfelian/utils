package cookies

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with Ginkgo", func() {
	cookiesString := "token=tokenString;murka= murkaString; zhmurka=aga"
	type nameValue map[string]string

	It("checks New", func() {
		cookies := New(cookiesString)
		Expect(cookies).To(HaveLen(3))

		expectedNewResult := nameValue{
			"token":   "tokenString",
			"murka":   " murkaString",
			"zhmurka": "aga",
		}

		for i, cookie := range cookies {
			By(fmt.Sprintf("testing case %d", i))
			Expect(cookie.Value).To(Equal(expectedNewResult[cookie.Name]))
		}
	})

	It("checks Get", func() {
		cookies := New(cookiesString)
		expectedGetResult := nameValue{
			"token":   "tokenString",
			"murka":   " murkaString",
			"zhmurka": "aga",
			"q":       "",
			"":        "",
		}

		for key, expectedValue := range expectedGetResult {
			By(fmt.Sprintf("testing case %s: %s", key, expectedValue))
			receivedCookie := cookies.Get(key)
			if expectedValue == "" {
				Expect(receivedCookie).To(BeNil())
				continue
			}
			Expect(receivedCookie.Value).To(Equal(expectedValue))
		}
	})

	It("checks GetValue", func() {
		cookies := New(cookiesString)
		expectedGetResult := nameValue{
			"token":   "tokenString",
			"murka":   " murkaString",
			"zhmurka": "aga",
			"q":       "",
			"":        "",
		}

		for key, expectedValue := range expectedGetResult {
			By(fmt.Sprintf("testing case %s: %s", key, expectedValue))
			Expect(cookies.GetValue(key)).To(Equal(expectedValue))
		}
	})
})
