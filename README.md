# Terraform Provider for SpiceDB


- A resource and a data source and [provider](./internal/provider),
- Examples [examples](./examples) and generated  [documentation](./docs/index.md),
- Miscellaneous meta files.

## Example

#### Configure provider
```terraform
terraform {
  required_providers {
    spicedb = {
      source = "educationperfect/spicedb"
    }
  }
}

provider "spicedb" {
  endpoint = "localhost:50051"
  token = "happylittlekey"
  insecure = true
}
```

#### Define schema with `spicedb_schema` resource
```terraform

resource "spicedb_schema" "test" {
   schema = <<EOF
     definition user {}

     definition organization {
         permission is_member = member
         relation member : user
     }
EOF
}
```

#### Use `spicedb_schema` data resource
```terraform

data "spicedb_schema" "test" { }

output "schema" {
  value = data.spicedb_schema.test.schema
}
```

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

or you can do it manually

```bash
docker compose up -d
TEST_ACC=1 go test -v -cover ./internal/provider/
```

## Local dev hack

Writing this snippet in `~/.terraformrc` allows Terraform to resolve and use
locally build provider:

```
provider_installation {
  dev_overrides {
    "educationperfect/spicedb" = "$(go env GOPATH)/bin"
  }
  direct {}
}
```