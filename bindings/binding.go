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
)

const (
	// Type                          = "type"
	DEFAULT_SERVICE_BINDINGS_ROOT = "/bindings"
	serviceBindingRoot            = "SERVICE_BINDING_ROOT"
)

type ServiceBinding struct {
	root string
}

type BindingsSpec struct {
	Host         string
	Port         uint16
	Uri          string
	Username     string
	Password     string
	SSL          bool
	Certificates []string
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
	sb := &ServiceBinding{root: root}
	return sb, nil
}

func NewBinding(Type ...string) (*BindingsSpec, error) {
	var t string
	result := make(map[string]string)
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

	log.Printf("Service Binding Root: %s\n", sb.root)
	err = filepath.Walk(sb.root, func(bpath string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("filepath.Walk error: %s\n", err.Error())
			return err
		}

		if info.Name() == "type" && info.Mode()&os.ModeSymlink != os.ModeSymlink {
			// log.Printf("File: %s", bpath)
			// fmt.Println(path.Dir(bpath))

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
				// log.Printf("file: %s - content: %s", f.Name(), fc)
				result[f.Name()] = string(fc)
			}
		}

		return nil
	})
	// log.Printf("Map: %s\n", result)
	if err != nil {
		log.Printf("filepath.Walk error: %q\n", err.Error())
		return nil, err
	}

	bs := new(BindingsSpec)

	// bs := new(BindingsSpec)
	port, _ := strconv.Atoi(result["port"])
	bs.Host = result["host"]
	bs.Port = uint16(port)
	bs.Username = result["username"]
	bs.Password = result["password"]
	bs.SSL = result["ssl"] == "true"
	bs.Certificates = []string{""}

	return bs, err
}
