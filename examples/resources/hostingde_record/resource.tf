# Manage example DNS record.
resource "hostingde_record" "example" {
  zone_id = hostingde_zone.sample.id
  name = "test.example.test"
  type = "CNAME"
  content = "www.example.com"
}
