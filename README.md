# Terraform Provider for [hosting.de](https://hosting.de)

Currently only supports maintaining DNS Zones and Records.

### Environment variables for provider configuration
- Required: `HOSTINGDE_AUTH_TOKEN`, go to your [hosting.de profile](https://secure.hosting.de/profile) and create an API Key
- Optional: `HOSTINGDE_ACCOUNT_ID`


### Example `main.tf`
```
terraform {
  required_providers {
    hostingde = {
      source = "hostingde/hostingde"
      version = ">= 0.0.1"
    }
  }
}

# Not recommended, use environment variables to configure the provider
#provider "hostingde" {
#  auth_token = "YOUR_API_TOKEN"
#  account_id = "YOUR_ACCOUNT_ID"
#}

resource "hostingde_zone" "sample" {
  name = "example.test"
  type = "NATIVE"
}

resource "hostingde_record" "example" {
  zone_id = hostingde_zone.sample.id
  name = "test.example.test"
  type = "CNAME"
  content = "www.example.com"
}

output "hostingde_zone" {
  value = hostingde_zone.sample
}

output "hostingde_record" {
  value = hostingde_record.example
}
```

# Development and testing
Prepare Terraform for local provider install
```shell
go env GOBIN
```

Add your GOBIN PATH to `~/.terraformrc`
```
provider_installation {

  dev_overrides {
      "registry.terraform.io/hostingde/hostingde" = "<PATH>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Compile and install the provider into your `$GOPATH/bin`

```shell
make install
```

Run resource tests
```shell
make testacc
```

Then, navigate to the `example` directory. 

```shell
cd example
```

Run the following command to initialize the workspace and apply the example configuration.

```shell
terraform init && terraform apply
```

Useful resources:
- [Tutorial for Terraform Plugin Framework](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider)
- [Terraform Plugin Framework documentation](https://developer.hashicorp.com/terraform/plugin/framework)
