# Service Bindind for Golang Apps

> Kubernetes Service Binding Library for Golang Apps

## Specifications

The recommanded service binding spec from [servicebinding.io](https://servicebinding.io/):

* `host` : A DNS hostname or IP (should resolve)
* `port` : A valid port number
* `uri` : A valid URI as defined by [RFC3986](https://tools.ietf.org/html/rfc3986)
* `username` : A string-based username credentials
* `password` : A string-based password credentials
* `database` : Extended Spec for database requirements
* `certificates` : A collection of PEM-encoded X.509 public certificates, representing a certificate chain used to trust TLS connections
* `private-key` : A PEM-encoded private key used in mTLS client authentication
* `ssl` : Extended spec with SSL enabled

And the `workload projection` :

```text
$SERVICE_BINDING_ROOT
├── account-database
│   ├── type
│   ├── provider
│   ├── uri
│   ├── username
│   └── password
└── transaction-event-stream
    ├── type
    ├── connection-count
    ├── uri
    ├── certificates
    └── private-key
```

## Get started

Download the `servicebinding` package with `go get` command:

```shell
go get github.com/bzhtux/servicebinding/bindings
```

To use the `servicebinding` package import it as below:

```go
import (
    "github.com/bzhtux/servicebinding/bindings"
)
```

Create a new `bindings` object like this:

```go
bindings.NewBinding("<binding type>")
```

Example:

```go
sb, err := bindings.NewBinding("redis")
if err != nil {
  log.Printf("NewBinding Error: %s\n", err.Error())
}
fmt.Printf("Redis Host = %s\n", sb.Host)
```

## Troubleshooting

Need more tests !!!

## Contribute

You can contribute in many ways:

* opening issues for bugs
* opening issues to request new useful features
* contributing to this repo