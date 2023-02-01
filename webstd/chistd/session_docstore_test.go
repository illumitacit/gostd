package chistd_test

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi/v5"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Implicit imports to register drivers for docstore
	_ "github.com/fensak-io/gostd/webstd/chistd"
	_ "gocloud.dev/docstore/memdocstore"
)

var _ = Describe("SessionDocstore", func() {
	var c chi.Router

	BeforeEach(func() {
		randBytes := make([]byte, 16)
		_, err := rand.Read(randBytes)
		Ω(err).ShouldNot(HaveOccurred())

		suffix := hex.EncodeToString(randBytes)
		opt := session.Options{
			Provider:       "docstore",
			ProviderConfig: "mem://session" + suffix,
		}
		c = chi.NewRouter()
		c.Use(session.Sessioner(opt))
		addTestSessionRoutes(c)
	})

	When("getting registered session", func() {
		It("works", func() {
			cookie := makeTestRequest(c, "", "/")
			makeTestRequest(c, cookie, "/get")
		})
	})

	When("regenerating existing session", func() {
		It("works", func() {
			cookie := makeTestRequest(c, "", "/")
			regCookie := makeTestRequest(c, cookie, "/regen")
			makeTestRequest(c, regCookie, "/get")
		})
	})

	When("regenerating empty session", func() {
		It("works", func() {
			cookie := "MacaronSession=ad2c7e3cbecfcf48; Path=/;"
			makeTestRequest(c, cookie, "/regen")
		})
	})
})

func addTestSessionRoutes(c chi.Router) {
	c.Get("/", func(resp http.ResponseWriter, req *http.Request) {
		sess := session.GetSession(req)
		err := sess.Set("uname", "unknwon")
		Ω(err).ShouldNot(HaveOccurred())
	})
	c.Get("/get", func(resp http.ResponseWriter, req *http.Request) {
		sess := session.GetSession(req)
		sid := sess.ID()
		Ω(sid).ShouldNot(BeEmpty())

		raw, err := sess.Read(sid)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(raw).ShouldNot(BeNil())

		uname := sess.Get("uname")
		Ω(uname).Should(Equal("unknwon"))

		delErr := sess.Delete("uname")
		Ω(delErr).ShouldNot(HaveOccurred())

		unameRef := sess.Get("uname")
		Ω(unameRef).Should(BeNil())

		delErrAgain := sess.Delete("uname")
		Ω(delErrAgain).ShouldNot(HaveOccurred())
	})
	c.Get("/regen", func(resp http.ResponseWriter, req *http.Request) {
		sess := session.GetSession(req)
		raw, err := sess.RegenerateID(resp, req)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(raw).ShouldNot(BeNil())
	})
}

func makeTestRequest(c chi.Router, cookie, endpoint string) string {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", endpoint, nil)
	Ω(err).ShouldNot(HaveOccurred())
	req.Header.Set("Cookie", cookie)
	c.ServeHTTP(resp, req)
	return resp.Header().Get("Set-Cookie")
}
