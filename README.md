# Terraform Provider for [hosting.de](https://hosting.de)

Currently only supports maintaining DNS Zones and Records.

- [Provider documentation](https://registry.terraform.io/providers/pub-solar/hostingde/latest/docs)

### Environment variables for provider configuration
- Required: `HOSTINGDE_AUTH_TOKEN`, go to your [hosting.de profile](https://secure.hosting.de/profile) and create an API Key
- Optional: `HOSTINGDE_ACCOUNT_ID`

### Quick start
Copy the example `main.tf`, then run the following commands:

```shell
export HOSTINGDE_AUTH_TOKEN=your-token
```
```shell
terraform init
```
```shell
terraform plan -out "hostingde.plan"
```
```shell
terraform apply "hostingde.plan"
```

### Example `main.tf`
```
terraform {
  required_providers {
    hostingde = {
      source = "pub-solar/hostingde"
      version = ">= 0.0.2"
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

### Importing existing zones and records
#### Zones
- Create `*.tf` files containing the existing zone and records you'd like to import
- Go to https://secure.hosting.de/dns/ then click "Show details" on the zone you'd like to import
- Copy the `ZONE_CONFIG_ID` from the URL: https://secure.hosting.de/dns/zones/id/$ZONE_CONFIG_ID/edit
- Import the `hostingde_zone` resouce using:
```shell
terraform import hostingde_zone.your_zone_name $ZONE_CONFIG_ID
```

#### Records
- Importing records is a little more involved, let's go:
- Write a shell function to prepare `curl` JSON data (this assumes you have your
  API token set in the environment and that you replace `$ZONE_CONFIG_ID` with
  the ID from above)
```shell
generate_post_data()
{
  cat <<EOF
{
    "authToken": "$HOSTINGDE_AUTH_TOKEN",
    "filter": {
        "field": "zoneConfigId",
        "value": "$ZONE_CONFIG_ID"
    },
    "limit": 10,
    "page": 1,
    "sort": {
        "field": "recordName",
        "order": "asc"
    }
}
EOF
}
```
- Do the curl POST request to get all DNS record IDs of that zone
```shell
curl \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -X POST \
  -d "$(generate_post_data)" \
  https://secure.hosting.de/api/dns/v1/json/recordsFind
```
- Example response:
```
{
    "errors": [
    ],
    "metadata": {
        "clientTransactionId": "",
        "serverTransactionId": "20230411151239132-dnsrobot-robots1-26486-0"
    },
    "response": {
        "data": [
            {
                "accountId": "ACCOUNT_ID",
                "addDate": "2023-02-03T13:33:26Z",
                "comments": "",
                "content": "\"v=DMARC1; p=reject;\"",
                "id": "RECORD_ID",
                "lastChangeDate": "2023-02-03T13:33:26Z",
                "name": "_dmarc.your.domain",
                "priority": null,
                "recordTemplateId": null,
                "ttl": 3600,
                "type": "TXT",
                "zoneConfigId": "ZONE_CONFIG_ID"
            },
            ...
```
- One by one, import the records:
```shell
terraform import hostingde_record.your_record_name $RECORD_ID
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
      "registry.terraform.io/pub-solar/hostingde" = "<PATH>"
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

Generate documentation to `docs/`

```shell
make docs
```

Run linters

```shell
make lint
```

Run resource tests
```shell
export HOSTINGDE_AUTH_TOKEN=YOUR-API-TOKEN

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
