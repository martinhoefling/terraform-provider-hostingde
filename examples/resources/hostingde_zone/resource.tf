# Manage example DNS zone.
resource "hostingde_zone" "sample" {
  name = "example.test"
  type = "NATIVE"
}
