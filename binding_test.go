package bindings_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	bindings "github.com/bzhtux/servicebinding"
	"github.com/spf13/afero"
)

var _ = Describe("Binding", func() {
	// var bs *bindings.BindingsSpec
	BeforeEach(func() {
		var TestFs = afero.NewOsFs()
		TestFs.MkdirAll("/tmp/testsb/..19700101.12345/pg", 0755)
		afero.WriteFile(TestFs, "/tmp/testsb/..19700101.12345/pg/type", []byte("postgresql"), 0644)
		os.Unsetenv("SERVICE_BINDING_ROOT")
	})

	Describe("Test SERVICE_BINDING_ROOT exists", func() {
		Context("Without env var set", func() {
			It("Should return an error", func() {

				_, err := bindings.GetServiceBindingRoot()
				Expect(err).NotTo(BeNil())
			})
		})
		Context("With env var set", func() {
			It("Should not return an error", func() {
				os.Setenv("SERVICE_BINDING_ROOT", "/bindings")
				sbr, err := bindings.GetServiceBindingRoot()
				Expect(err).To(BeNil())
				Expect(sbr.Root).To(Equal("/bindings"))
			})
		})
	})

	Describe("Test NewBinding", func() {
		Context("Without fs layout set", func() {
			It("Should return an error", func() {
				var TestFs = afero.NewOsFs()
				TestFs.RemoveAll("/tmp/testsb")
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				_, err := bindings.NewBinding("postgresql")
				Expect(err).NotTo(BeNil())
			})
		})
		Context("With fs layout set", func() {
			It("Should return value for type", func() {
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				nb, _ := bindings.NewBinding("postgresql")
				Expect(nb.Type).To(Equal("postgresql"))
			})
		})
		Context("With fs layout set and ssl not a boolean", func() {
			It("Should return false", func() {
				var TestFs = afero.NewOsFs()
				afero.WriteFile(TestFs, "/tmp/testsb/..19700101.12345/pg/ssl", []byte("NothingElseMatter"), 0644)
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				nb, _ := bindings.NewBinding("postgresql")
				Expect(nb.SSL).To(Equal(false))
			})
		})
		Context("With fs layout set and ssl a boolean", func() {
			It("Should return true", func() {
				var TestFs = afero.NewOsFs()
				afero.WriteFile(TestFs, "/tmp/testsb/..19700101.12345/pg/ssl", []byte("1"), 0644)
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				nb, _ := bindings.NewBinding("postgresql")
				Expect(nb.SSL).To(Equal(true))
			})
		})
		Context("With fs layout set and port not an integer", func() {
			It("Should return an error", func() {
				var TestFs = afero.NewOsFs()
				afero.WriteFile(TestFs, "/tmp/testsb/..19700101.12345/pg/port", []byte("fivefourthreetwo"), 0644)
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				_, err := bindings.NewBinding("postgresql")
				Expect(err).NotTo(BeNil())
			})
		})
		Context("With fs layout set and port an integer", func() {
			It("Should not return an error and should return the same port number: uint16(5432)", func() {
				var TestFs = afero.NewOsFs()
				afero.WriteFile(TestFs, "/tmp/testsb/..19700101.12345/pg/port", []byte("5432"), 0644)
				os.Setenv("SERVICE_BINDING_ROOT", "/tmp/testsb")
				nb, err := bindings.NewBinding("postgresql")
				Expect(err).To(BeNil())
				Expect(nb.Port).To(Equal(uint16(5432)))
			})
		})
	})
})
