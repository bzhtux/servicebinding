package bindings

import (
	"errors"
	"fmt"
	"io/fs"
	_ "log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
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
	Provider     string   `mapstructure:"provider"`
	Certificates []string `mapstructure:"certificates"`
	PrivateKey   string   `mapstructure:"privatekey"`
	Type         string   `mapstructure:"type"`
}

func GetServiceBindingRoot() (*ServiceBinding, error) {
	root, exists := os.LookupEnv(serviceBindingRoot)
	if !exists {
		err_msg := fmt.Sprintf("Environment variable not set: %s", serviceBindingRoot)
		return nil, errors.New(err_msg)
	}
	sb := &ServiceBinding{
		Root: root,
	}
	return sb, nil
}

func NewBinding(bindingtype string) (*BindingsSpec, error) {
	bs := new(BindingsSpec)
	result := make(map[string]interface{})

	sb, err := GetServiceBindingRoot()
	if err != nil {
		// log.Printf("GetServiceBindingRoot error: %s\n", err.Error())
		return nil, err
	}

	// log.Printf("Found ServiceBinding.Root %s, processing ...\n", sb.Root)
	// Walk through SERVICE_BINDING_ROOT directory
	err = filepath.Walk(sb.Root, func(bpath string, info fs.FileInfo, err error) error {
		if err != nil {
			// log.Printf("Error with ServiceBinding.Root = %s : %s\n", sb.Root, err.Error())
			return err
		}

		// If file is not a link continue walking, else stop (do nothing)
		if info.Name() == "type" && info.Mode()&os.ModeSymlink != os.ModeSymlink {

			fct, err := os.ReadFile(bpath)
			// If file content match the binding type (e.g postgresql)
			// loop into the parent directory to list files and get content
			if err == nil && string(fct) == bindingtype {
				files, err := os.ReadDir(path.Dir(bpath))
				if err != nil {
					return err
				}

				for _, f := range files {
					fc, err := os.ReadFile(filepath.Join(path.Dir(bpath), f.Name()))
					if err != nil {
						return err
					}
					if f.Name() == "port" {
						dataPort, err := strconv.Atoi(string(fc))
						if err != nil {
							return err
						}
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
		if err != nil {
			return err
		} else {
			return nil
		}

	})

	if err != nil {
		return nil, err
	}

	if err := mapstructure.Decode(result, &bs); err != nil {
		return nil, err
	} else {
		return bs, err
	}
}
