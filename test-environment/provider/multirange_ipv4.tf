resource "netaddr_range_ipv4" "multirange_ipv4" {
    key_prefix = "/test/multirange-ipv4/"
    first_address = "192.170.0.1"
    last_address = "192.170.0.2"
}

resource "netaddr_range_ipv4" "multirange_ipv4_range2" {
    key_prefix = "/test/multirange-ipv4-range2/"
    first_address = "192.170.0.10"
    last_address = "192.170.0.20"
}

resource "netaddr_address_ipv4_v2" "multirange_ipv4_addr1" {
    range_ids = [netaddr_range_ipv4.multirange_ipv4.id, netaddr_range_ipv4.multirange_ipv4_range2.id]
    name = "addr1"
}

output "multirange_ipv4_addr1" {
  value = netaddr_address_ipv4_v2.multirange_ipv4_addr1
}

resource "netaddr_address_ipv4_v2" "multirange_ipv4_addr2" {
    range_ids = [netaddr_range_ipv4.multirange_ipv4.id, netaddr_range_ipv4.multirange_ipv4_range2.id]
    name = "addr2"
}

output "multirange_ipv4_addr2" {
  value = netaddr_address_ipv4_v2.multirange_ipv4_addr2
}

resource "netaddr_address_ipv4_v2" "multirange_ipv4_addr3" {
    range_ids = [netaddr_range_ipv4.multirange_ipv4.id, netaddr_range_ipv4.multirange_ipv4_range2.id]
    name = "addr3"
}

output "multirange_ipv4_addr3" {
  value = netaddr_address_ipv4_v2.multirange_ipv4_addr3
}

resource "netaddr_address_ipv4_v2" "multirange_ipv4_addr4" {
    range_ids = [netaddr_range_ipv4.multirange_ipv4.id, netaddr_range_ipv4.multirange_ipv4_range2.id]
    name = "addr4"
}

output "multirange_ipv4_addr4" {
  value = netaddr_address_ipv4_v2.multirange_ipv4_addr4
}

resource "netaddr_address_ipv4_v2" "multirange_ipv4_addr5" {
    range_ids = [netaddr_range_ipv4.multirange_ipv4.id, netaddr_range_ipv4.multirange_ipv4_range2.id]
    name = "addr5"
    hardcoded_address = "192.170.0.20"
}

output "multirange_ipv4_addr5" {
  value = netaddr_address_ipv4_v2.multirange_ipv4_addr5
}