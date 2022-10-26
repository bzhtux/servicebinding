package bindings_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bzhtux/servicebinding/bindings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BindingsSpec", func() {

	var (
		sb *bindings.ServiceBinding
		// SoftBinding *BindingsSpec
	)

	BeforeEach(func() {
		folderPath := "redis/..20420101"
		sbd, err := ioutil.TempDir("/tmp", "bindings")

		sb = &bindings.ServiceBinding{
			Root: sbd,
		}

		if err != nil {
			log.Printf("Error creating tmp dir: %s\n", err)
		}
		if err := os.MkdirAll(filepath.Join(sbd, folderPath), 0755); err != nil {
			log.Printf("Error creating folderPath: %s\n", err.Error())
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "type"), []byte("redis"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/type\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "host"), []byte("redis_host"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/host\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "username"), []byte("redis_username"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/username\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "password"), []byte("redis_password"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/password\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "port"), []byte("6379"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/port\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "database"), []byte("0"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/database\n", sbd, folderPath)
		}
		if err := ioutil.WriteFile(filepath.Join(sbd, folderPath, "ssl"), []byte("false"), 0644); err != nil {
			log.Printf("Error writing file: %s/%s/ssl\n", sbd, folderPath)
		}
		// defer os.RemoveAll(sbd)
	})

	Describe("Run NewServiceBinding", func() {
		Context("Without SERVICE_BINDING_ROOT env var defined", func() {
			It("Should return an Error", func() {
				os.Unsetenv("SERVICE_BINDING_ROOT")
				_, err := bindings.NewServiceBinding()
				Expect(err.Error()).To(Equal("environment variable not set: SERVICE_BINDING_ROOT"))
			})
		})
		Context("With SERVICE_BINDING_ROOT env var set to sbd", func() {
			It("Should return the SERVICE_BINDING_ROOT dir", func() {
				os.Setenv("SERVICE_BINDING_ROOT", sb.Root)
				s, _ := bindings.NewServiceBinding()
				Expect(s.Root).To(Equal(sb.Root))
			})
			It("Should return Redis host equal to redis_host", func() {
				os.Setenv("SERVICE_BINDING_ROOT", sb.Root)
				// s, _ := bindings.NewServiceBinding()
				// fmt.Printf("*** Service binding root: %s\n", s.Root)
				sb, _ := bindings.NewBinding("redis")
				Expect(sb.Host).To(Equal("redis_host"))
			})
		})
	})
	AfterEach(func() {
		os.RemoveAll(sb.Root)
	})
})
