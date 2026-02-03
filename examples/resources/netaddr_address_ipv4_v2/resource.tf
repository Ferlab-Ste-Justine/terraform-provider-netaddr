resource "netaddr_range_ipv4" "range1" {
    key_prefix = "/test/ipv4/"
    first_address = "192.168.0.1"
    last_address = "192.168.0.64"
}

resource "netaddr_range_ipv4" "range2" {
    key_prefix = "/test/ipv4-extras/"
    first_address = "192.168.0.128"
    last_address = "192.168.0.254"
}

resource "netaddr_address_ipv4_v2" "test" {
    range_ids = [netaddr_range_ipv4.range1.id, netaddr_range_ipv4.range2.id]
    name = "test"
    hardcoded_address = "192.168.0.5"
}

resource "netaddr_address_ipv4_v2" "test2" {
    range_ids = [netaddr_range_ipv4.range1.id, netaddr_range_ipv4.range2.id]
    name = "test2"
}

output "test_addr" {
  value = netaddr_address_ipv4.test.address
}

output "test2_addr" {
  value = netaddr_address_ipv4.test2.address
}