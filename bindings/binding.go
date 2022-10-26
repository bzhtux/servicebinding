package bindings

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	. "github.com/mitchellh/mapstructure"
)

const (
	DEFAULT_SERVICE_BINDINGS_ROOT = "/bindings"
	serviceBindingRoot            = "SERVICE_BINDING_ROOT"
)

type ServiceBinding struct {
	Root string
}

type BindingsSpec struct {
	Host         string   `mapstructure:"host"`
	Port         uint16   `mapstructure:"port"`
	Uri          string   `mapstructure:"uri"`
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
	Database     string   `mapstructure:"database"`
	SSL          bool     `mapstructure:"ssl"`
	Certificates []string `mapstructure:"certificates"`
}

type Binding interface {
	NewServiceBinding()
	NewBinding()
}

type SecretSpec struct {
	*BindingsSpec
	Type     string
	Provider string
}

func NewServiceBinding() (*ServiceBinding, error) {
	root, exists := os.LookupEnv(serviceBindingRoot)
	if !exists {
		return nil, errors.New("environment variable not set: " + serviceBindingRoot)
	}
	sb := &ServiceBinding{Root: root}
	return sb, nil
}

// Get only one binding type (e.g redis) even if Type is a slice
func NewBinding(Type ...string) (*BindingsSpec, error) {
	var t string
	result := make(map[string]interface{})
	if len(Type) == 0 {
		log.Fatal("No binding provided")
	} else {
		t = Type[0]
		fmt.Printf("Binding type: %s\n", t)
	}
	sb, err := NewServiceBinding()
	if err != nil {
		log.Printf("NewServiceBindings error: %s\n", err.Error())
		return nil, err
	}

	log.Printf("Service Binding Root: %s\n", sb.Root)
	err = filepath.Walk(sb.Root, func(bpath string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("filepath.Walk error: %s\n", err.Error())
			return err
		}

		if info.Name() == "type" && info.Mode()&os.ModeSymlink != os.ModeSymlink {

			fct, err := os.ReadFile(bpath)
			if err == nil && string(fct) == t {
				files, err := os.ReadDir(path.Dir(bpath))
				if err != nil {
					log.Printf("Error reading dir %s", bpath)
					return err
				}

				for _, f := range files {
					fc, err := os.ReadFile(filepath.Join(path.Dir(bpath), f.Name()))
					if err != nil {
						log.Printf("Error getting file content: %s/%s", path.Dir(bpath), f.Name())
						return err
					}
					if f.Name() == "port" {
						dataPort, _ := strconv.Atoi(string(fc))
						result[f.Name()] = dataPort
					} else {
						if f.Name() == "ssl" {
							boolSSL, _ := strconv.ParseBool(string(fc))
							result[f.Name()] = boolSSL
						} else {
							result[f.Name()] = string(fc)
						}
					}

				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("filepath.Walk error: %q\n", err.Error())
		return nil, err
	}

	bs := new(BindingsSpec)
	if err := Decode(result, &bs); err != nil {
		log.Printf("Error when decoding result into struct: %s\n", err.Error())
		return nil, err
	} else {
		return bs, err
	}
}
