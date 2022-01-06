resource "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.0.254"
}