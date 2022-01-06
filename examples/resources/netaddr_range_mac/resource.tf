resource "netaddr_range_mac" "test" {
    key_prefix = "/test/mac/"
    first_address = "52:54:00:00:00:00"
    last_address = "52:54:00:ff:ff:ff"
}