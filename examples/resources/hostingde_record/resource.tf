# Manage example DNS record.
resource "hostingde_record" "example" {
  zone_id = hostingde_zone.sample.id
  name = "test.example.test"
  type = "CNAME"
  content = "www.example.com"
  ttl = 300
}

# Manage example DNS MX record.
resource "hostingde_record" "example" {
  zone_id = hostingde_zone.sample.id
  name = "test.example.test"
  type = "MX"
  content = "mail.example.com"
  ttl = 300
  priority = 10
}
