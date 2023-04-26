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
