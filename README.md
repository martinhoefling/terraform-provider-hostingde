# Terraform Provider for [hosting.de](https://hosting.de)

Currently only supports maintaining DNS Zones and Records.

### Environment variables for provider configuration
- Required: `HOSTINGDE_AUTH_TOKEN`, go to your [hosting.de profile](https://secure.hosting.de/profile) and create an API Key
- Optional: `HOSTINGDE_ACCOUNT_ID`

# Development and testing

Compile and install the provider into your `$GOPATH/bin`

```shell
make install
```

Then, navigate to the `example` directory. 

```shell
cd example
```

Run the following command to initialize the workspace and apply the example configuration.

```shell
terraform init && terraform apply
```

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
