resource "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.0.254"
}

resource "netaddr_address_ipv4" "test" {
    range_id = netaddr_range_ipv4.test.id
    name = "test"
    hardcoded_address = "192.168.0.5"
}

resource "netaddr_address_ipv4" "test2" {
    range_id = netaddr_range_ipv4.test.id
    name = "test2"
}

output "test_addr" {
  value = netaddr_address_ipv4.test.address
}

output "test2_addr" {
  value = netaddr_address_ipv4.test2.address
}