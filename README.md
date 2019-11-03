# Terraform Provider for Hosting.de

Currently only supports maintaining DNS Zones and Records.

# Build and Install

    go build -o terraform-provider-hostingde
    cp terraform-provider-hostingde ~/.terraform/plugins

# Example TF File

    provider "hostingde" {}

    resource "hostingde_zone" "sample" {
      name = "sample.example.com"
    }
    
    resource "hostingde_record" "example" {
      zone_id = hostingde_zone.sample.id
      name = "test.sample.example.com"
      type = "CNAME"
      content = "www.example.com"
    }

